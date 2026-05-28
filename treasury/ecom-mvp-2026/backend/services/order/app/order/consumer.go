package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	commonkafka "gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
)

// --- Inbound event payload shapes ---

// paymentEventPayload is shared across payment.completed, payment.failed, payment.expired.
type paymentEventPayload struct {
	EventID   string `json:"eventId"`
	OrderID   string `json:"orderId"`
	PaymentID string `json:"paymentId,omitempty"`
}

// reservationExpiredPayload is the payload for events.reservation.expired.
type reservationExpiredPayload struct {
	EventID       string `json:"eventId"`
	OrderID       string `json:"orderId"`
	ReservationID string `json:"reservationId,omitempty"`
}

// --- Event handler registration ---

// RegisterConsumerRoutes returns the kafka.KafkaHandler map for all topics consumed
// by this service. Routes are keyed by EventName (matching the common/kafka Message envelope).
func (s *Service) RegisterConsumerRoutes() map[string]commonkafka.KafkaHandler {
	return map[string]commonkafka.KafkaHandler{
		"payment.completed":    s.onPaymentCompleted,
		"payment.failed":       s.onPaymentFailed,
		"payment.expired":      s.onPaymentExpired,
		"reservation.expired":  s.onReservationExpired,
	}
}

// --- FULL handler: events.payment.completed ---
//
// State-driven PENDING_PAYMENT → PAID transition.
// Two-layer idempotency:
//  1. consumed_events PK on event_id (absorbs Kafka redelivery)
//  2. State-driven check inside FOR UPDATE (absorbs cross-topic interleavings)
//
// Pseudocode matches td.json §event_consumers[0].pseudocode exactly.
func (s *Service) onPaymentCompleted(ctx context.Context, msg commonkafka.Message[json.RawMessage]) error {
	var payload paymentEventPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("payment.completed: unmarshal payload: %w", err)
	}

	eventID := payload.EventID
	if eventID == "" {
		eventID = msg.EventID // fall back to envelope event_id
	}
	orderID := payload.OrderID
	if orderID == "" {
		slog.WarnContext(ctx, "payment.completed: missing orderId, skipping")
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("payment.completed: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Layer 1: consumed_events dedup (PK conflict = already processed → ack).
	inserted, err := s.consumed.InsertOrConflict(ctx, tx, eventID, "order.payment.completed", "payment.completed", orderID)
	if err != nil {
		return fmt.Errorf("payment.completed: insert consumed_event: %w", err)
	}
	if !inserted {
		slog.InfoContext(ctx, "payment.completed: duplicate event, skipping", "eventId", eventID, "orderId", orderID)
		return tx.Commit(ctx)
	}

	// SELECT ... FOR UPDATE (sole-writer invariant — no other service touches orders.status).
	order, err := s.orders.GetByIDForUpdate(ctx, tx, orderID)
	if err != nil {
		return fmt.Errorf("payment.completed: select for update: %w", err)
	}
	if order == nil {
		slog.WarnContext(ctx, "payment.completed: unknown orderId, skipping", "orderId", orderID)
		return tx.Commit(ctx)
	}

	// Layer 2: state-driven no-op for any status other than PENDING_PAYMENT.
	if order.Status != StatusPendingPayment {
		slog.InfoContext(ctx, "payment.completed: ignored — status not PENDING_PAYMENT",
			"orderId", orderID,
			"status", order.Status,
			"eventId", eventID,
		)
		return tx.Commit(ctx)
	}

	// Perform transition: PENDING_PAYMENT → PAID.
	now := time.Now().UTC()
	if err := s.orders.UpdateStatus(ctx, tx, orderID, StatusPaid, order.Version, nil); err != nil {
		return fmt.Errorf("payment.completed: update status: %w", err)
	}

	reason := fmt.Sprintf("payment.completed eventId=%s", eventID)
	histRow := OrderStatusHistory{
		ID:          uuid.New().String(),
		OrderID:     orderID,
		FromStatus:  ptrStatus(StatusPendingPayment),
		ToStatus:    StatusPaid,
		ActorUserID: nil,
		ActorRole:   RoleSystem,
		Reason:      &reason,
		OccurredAt:  now,
	}
	if err := s.statusHistory.Insert(ctx, tx, histRow); err != nil {
		return fmt.Errorf("payment.completed: insert history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("payment.completed: commit: %w", err)
	}

	slog.InfoContext(ctx, "payment.completed: order transitioned to PAID",
		"orderId", orderID,
		"eventId", eventID,
	)
	return nil
}

// --- STUB handlers ---

// onPaymentFailed is a stub for events.payment.failed (PENDING_PAYMENT → PAYMENT_FAILED).
// TODO: implement state-driven framework matching onPaymentCompleted pattern.
func (s *Service) onPaymentFailed(ctx context.Context, msg commonkafka.Message[json.RawMessage]) error {
	slog.InfoContext(ctx, "TODO: state-driven handler", "event", "payment.failed", "eventId", msg.EventID)
	return nil
}

// onPaymentExpired is a stub for events.payment.expired (PENDING_PAYMENT → PAYMENT_EXPIRED).
// TODO: implement state-driven framework. Note: race with events.reservation.expired resolved by FOR UPDATE + state-driven guard (td.json §deduplication_with_reservation_expired).
func (s *Service) onPaymentExpired(ctx context.Context, msg commonkafka.Message[json.RawMessage]) error {
	slog.InfoContext(ctx, "TODO: state-driven handler", "event", "payment.expired", "eventId", msg.EventID)
	return nil
}

// onReservationExpired is a stub for events.reservation.expired (PENDING_PAYMENT → PAYMENT_EXPIRED).
// TODO: implement state-driven framework. Same final state as payment.expired; first-writer-wins via FOR UPDATE (td.json §deduplication_with_reservation_expired PR-008).
func (s *Service) onReservationExpired(ctx context.Context, msg commonkafka.Message[json.RawMessage]) error {
	slog.InfoContext(ctx, "TODO: state-driven handler", "event", "reservation.expired", "eventId", msg.EventID)
	return nil
}

// ptrStatus is a helper that returns a pointer to an OrderStatus value.
func ptrStatus(s OrderStatus) *OrderStatus {
	return &s
}
