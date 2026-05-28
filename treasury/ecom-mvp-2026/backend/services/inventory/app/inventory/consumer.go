package inventory

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
)

// ─── Consumer event topic / name constants ────────────────────────────────────

const (
	// Topics (match contracts.json).
	TopicPaymentEvents  = "ecom.payment.events"
	TopicOrderEvents    = "ecom.order.events"
	TopicCatalogEvents  = "ecom.catalog.events"

	// Event type strings used as routing keys.
	EventPaymentCompleted = "payment.completed"
	EventPaymentFailed    = "payment.failed"
	EventPaymentExpired   = "payment.expired"
	EventOrderCancelled   = "order.cancelled"
	EventProductCreated   = "product.created"

	// Consumer names for the consumed_events dedup table.
	ConsumerPaymentCompleted = "inventory.payment-completed"
	ConsumerPaymentFailed    = "inventory.payment-failed"
	ConsumerPaymentExpired   = "inventory.payment-expired"
	ConsumerOrderCancelled   = "inventory.order-cancelled"
	ConsumerProductCreated   = "inventory.product-created"
)

// ─── Shared payload shapes ────────────────────────────────────────────────────

// paymentEventPayload is the common payload shape for payment.* events.
type paymentEventPayload struct {
	EventID string `json:"eventId"`
	OrderID string `json:"orderId"`
}

// orderCancelledPayload is the payload shape for order.cancelled.
// NOTE: consumers MUST NOT branch on fromStatus — read current DB row state.
// See td.json event_consumers[order-cancelled].notes.
type orderCancelledPayload struct {
	EventID    string `json:"eventId"`
	OrderID    string `json:"orderId"`
	FromStatus string `json:"fromStatus"` // informational only — DO NOT use for branching
}

// productCreatedPayload is the payload for product.created.
type productCreatedPayload struct {
	EventID string `json:"eventId"`
	SKU     string `json:"sku"`
}

// ─── Consumer event router ────────────────────────────────────────────────────

// EventHandlers returns the map of event_type → KafkaHandler for the inventory service.
// Register this map with the common/kafka event router.
func (svc *Service) EventHandlers() map[string]kafka.KafkaHandler {
	return map[string]kafka.KafkaHandler{
		EventPaymentCompleted: svc.handlePaymentCompleted,
		EventPaymentFailed:    svc.handlePaymentFailed,
		EventPaymentExpired:   svc.handlePaymentExpired,
		EventOrderCancelled:   svc.handleOrderCancelled,
		EventProductCreated:   svc.handleProductCreated,
	}
}

// ─── FULL: events.payment.completed ──────────────────────────────────────────

// handlePaymentCompleted processes events.payment.completed.
// State-driven per TD spec and PR-003 race resolution:
//   - Reads CURRENT reservation row state; ignores any fromStatus on the event.
//   - RESERVED  → stock: reserved_qty -= qty, sold_qty += qty; COMMITTED.
//   - COMMITTED → no-op (idempotent re-delivery).
//   - RELEASED/EXPIRED → log conflict metric; no-op.
//
// Dedup: INSERT INTO consumed_events ON CONFLICT DO NOTHING; 0 rows → already processed.
//
// Lock-order (Option A two-phase read — matches TD spec concurrency_model.lock_ordering):
//
//	Phase 1 — Snapshot-read reservations WITHOUT FOR UPDATE (FindByOrderID) to learn
//	           the ordered SKU set. No row-locks held here.
//	Phase 2 — For each row in lex SKU order: lock stock_levels(sku) FOR UPDATE first,
//	           then re-lock the reservation row by PK (GetByIDForUpdate).
//	This matches the sweeper's lock order and eliminates the cross-path deadlock
//	between concurrent sweep + payment.completed on the same orderId.
func (svc *Service) handlePaymentCompleted(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
	var payload paymentEventPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return serror.Wrap(err)
	}

	eventID, err := uuid.Parse(payload.EventID)
	if err != nil {
		return serror.Wrap(err).With(slog.String("raw_event_id", payload.EventID))
	}
	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return serror.Wrap(err).With(slog.String("raw_order_id", payload.OrderID))
	}

	tx, err := svc.Pool.Begin(ctx)
	if err != nil {
		return serror.Wrap(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Dedup: INSERT ON CONFLICT DO NOTHING.
	inserted, err := svc.ConsumedStorage.TryInsert(ctx, tx, eventID, ConsumerPaymentCompleted)
	if err != nil {
		return serror.Wrap(err)
	}
	if !inserted {
		// Already processed — commit and ack.
		return tx.Commit(ctx)
	}

	// ── Phase 1: snapshot read (no lock) ────────────────────────────────────────
	// LOCK-ORDER PIN (Option A): read reservations WITHOUT FOR UPDATE to get the
	// ordered SKU set. Acquiring reservations locks here (before stock_levels) would
	// invert the pin and risk deadlock against the sweeper or release path.
	snaps, err := svc.ResvStorage.FindByOrderID(ctx, tx, orderID)
	if err != nil {
		return serror.Wrap(err)
	}

	// ── Phase 2: lock stock_levels first (lex SKU order), then re-lock each row ─
	// snaps are already ORDER BY sku ASC from the snapshot query.
	for _, snap := range snaps {
		switch snap.Status {
		case StatusReserved:
			// LOCK-ORDER PIN: acquire stock_levels(sku) BEFORE reservations(id).
			if _, err := svc.StockStorage.GetBySKUForUpdate(ctx, tx, snap.SKU); err != nil {
				return serror.Wrap(err).With(slog.String("sku", snap.SKU))
			}
			// Re-lock the reservation row by PK now that stock_levels is held.
			r, err := svc.ResvStorage.GetByIDForUpdate(ctx, tx, snap.ID)
			if err != nil {
				return serror.Wrap(err).With(slog.String("reservation_id", snap.ID.String()))
			}
			if r.Status != StatusReserved {
				// Raced to RELEASED/EXPIRED between snapshot and re-lock — no-op.
				slog.WarnContext(ctx, "payment.completed: reservation status changed between snapshot and lock",
					slog.String("order_id", orderID.String()),
					slog.String("reservation_id", r.ID.String()),
					slog.String("current_status", string(r.Status)),
				)
				continue
			}
			if err := svc.StockStorage.DecrementReservedIncrementSold(ctx, tx, r.SKU, r.Qty); err != nil {
				return serror.Wrap(err).With(slog.String("sku", r.SKU))
			}
			if err := svc.ResvStorage.MarkCommitted(ctx, tx, r.ID); err != nil {
				return serror.Wrap(err).With(slog.String("reservation_id", r.ID.String()))
			}

		case StatusCommitted:
			// Idempotent re-delivery — no-op.

		case StatusReleased, StatusExpired:
			// PR-003 race: payment.completed arrived after sweeper or release consumer
			// already moved the reservation. Log and no-op.
			slog.WarnContext(ctx, "payment.completed arrived after reservation released/expired",
				slog.String("order_id", orderID.String()),
				slog.String("reservation_id", snap.ID.String()),
				slog.String("current_status", string(snap.Status)),
			)
			// TODO: increment metric inventory_reservation_conflict_total{event=payment.completed}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return serror.Wrap(err)
	}
	return nil
}

// ─── STUB: events.payment.failed ─────────────────────────────────────────────

// handlePaymentFailed processes events.payment.failed.
// TODO: wire state-driven handler per TD spec (event_consumers[payment-failed]):
//   - Dedup via consumed_events.
//   - FindByOrderIDForUpdate ORDER BY sku ASC.
//   - RESERVED → release (available += qty, reserved -= qty); release_reason=PAYMENT_FAILED.
//   - COMMITTED → defensive log; no-op (should be unreachable: PAYMENT_INTENT_TERMINAL guard).
//   - RELEASED/EXPIRED → no-op.
func (svc *Service) handlePaymentFailed(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
	slog.DebugContext(ctx, "payment.failed received (STUB)",
		slog.String("event_id", msg.EventID),
		slog.String("aggregate_id", msg.AggregateID),
	)
	// TODO: implement per td.json consumer inventory.consumer.payment-failed
	return nil
}

// ─── STUB: events.payment.expired ────────────────────────────────────────────

// handlePaymentExpired processes events.payment.expired.
// TODO: wire state-driven handler per TD spec (event_consumers[payment-expired]):
//   - Identical to payment.failed for inventory's purposes per contract.
//   - release_reason=PAYMENT_EXPIRED.
func (svc *Service) handlePaymentExpired(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
	slog.DebugContext(ctx, "payment.expired received (STUB)",
		slog.String("event_id", msg.EventID),
		slog.String("aggregate_id", msg.AggregateID),
	)
	// TODO: implement per td.json consumer inventory.consumer.payment-expired
	return nil
}

// ─── STUB: events.order.cancelled ────────────────────────────────────────────

// handleOrderCancelled processes events.order.cancelled.
// TODO: wire state-driven handler per TD spec (event_consumers[order-cancelled]):
//   - CRITICAL: NEVER branch on event.fromStatus — read CURRENT DB row state (PR-003).
//   - RESERVED  → release reserved→available; release_reason=ORDER_CANCELLED.
//   - COMMITTED → release sold→available; release_reason=PAID_SOFT_CANCEL (PR-006 BA MVP).
//   - RELEASED/EXPIRED → no-op.
func (svc *Service) handleOrderCancelled(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
	slog.DebugContext(ctx, "order.cancelled received (STUB)",
		slog.String("event_id", msg.EventID),
		slog.String("aggregate_id", msg.AggregateID),
	)
	// TODO: implement per td.json consumer inventory.consumer.order-cancelled.
	// IMPORTANT: DO NOT USE event.fromStatus for branching — read current reservation row.
	return nil
}

// ─── STUB: events.product.created ────────────────────────────────────────────

// handleProductCreated processes events.product.created.
// TODO: wire state-driven handler per TD spec (event_consumers[product-created]):
//   - Dedup via consumed_events.
//   - INSERT INTO stock_levels (sku, 0, 0, 0, 0) ON CONFLICT (sku) DO NOTHING.
//   - Initial qty = 0; admin uses stock.adjust to seed inventory.
func (svc *Service) handleProductCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
	slog.DebugContext(ctx, "product.created received (STUB)",
		slog.String("event_id", msg.EventID),
		slog.String("aggregate_id", msg.AggregateID),
	)
	// TODO: implement per td.json consumer inventory.consumer.product-created
	return nil
}
