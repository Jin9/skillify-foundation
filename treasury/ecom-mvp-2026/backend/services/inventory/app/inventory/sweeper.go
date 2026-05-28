package inventory

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/config"
)

// Sweeper runs the reservation-expiry sweep goroutine.
// It fires every cfg.Sweeper.CadenceSeconds seconds and processes up to
// cfg.Sweeper.BatchSize expired reservations per tick.
//
// Each tick:
//  1. Candidate scan: SELECT id,order_id,sku,qty FROM reservations
//     WHERE status='RESERVED' AND expires_at < NOW()
//     ORDER BY expires_at ASC LIMIT batch FOR UPDATE SKIP LOCKED
//     (runs in a short tx; SKIP LOCKED enables horizontal replica safety).
//  2. Per-row tx: re-lock stock_levels(sku) FIRST, then re-lock reservations(id)
//     — satisfying the canonical lock-order pin. Flip to EXPIRED, adjust stock,
//     write outbox row for events.reservation.expired.
type Sweeper struct {
	svc *Service
	cfg config.Config
}

// NewSweeper constructs a Sweeper.
func NewSweeper(svc *Service, cfg config.Config) *Sweeper {
	return &Sweeper{svc: svc, cfg: cfg}
}

// Run starts the sweeper loop. Blocks until ctx is cancelled.
// Start as a goroutine from main.
//
// Named constants used (NOT magic literals):
//   - s.cfg.Sweeper.CadenceSeconds  (SWEEPER_CADENCE_SECONDS, default 30)
//   - s.cfg.Sweeper.BatchSize       (SWEEPER_BATCH_SIZE, default 200)
func (s *Sweeper) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.cfg.Sweeper.CadenceSeconds) * time.Second)
	defer ticker.Stop()

	slog.InfoContext(ctx, "sweeper started",
		slog.Int("cadence_seconds", s.cfg.Sweeper.CadenceSeconds),
		slog.Int("batch_size", s.cfg.Sweeper.BatchSize),
	)

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "sweeper stopped")
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

// tick runs a single sweep iteration.
func (s *Sweeper) tick(ctx context.Context) {
	expiredCount := 0
	skippedCount := 0
	errCount := 0

	// ── Candidate scan ───────────────────────────────────────────────────────
	// Short tx just to scan candidates with SKIP LOCKED.
	// Committing this tx releases the locks; each candidate is then
	// re-locked individually in its own per-row tx.
	scanTx, err := s.svc.Pool.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "sweeper candidate scan begin tx", slog.String("error", err.Error()))
		return
	}

	candidates, err := s.svc.ResvStorage.FindExpiredCandidates(ctx, scanTx, s.cfg.Sweeper.BatchSize)
	if err != nil {
		_ = scanTx.Rollback(ctx)
		slog.ErrorContext(ctx, "sweeper candidate scan query", slog.String("error", err.Error()))
		return
	}
	// Commit the scan tx — releases SKIP LOCKED row locks so other replicas can see them.
	// Each row is re-locked individually below.
	if err := scanTx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "sweeper candidate scan commit", slog.String("error", err.Error()))
		return
	}

	if len(candidates) == 0 {
		return // nothing expired this tick
	}

	// ── Per-row tx ──────────────────────────────────────────────────────────
	for _, c := range candidates {
		released, err := s.processCandidate(ctx, c)
		if err != nil {
			errCount++
			slog.ErrorContext(ctx, "sweeper process candidate error",
				slog.String("reservation_id", c.ID.String()),
				slog.String("order_id", c.OrderID.String()),
				slog.String("error", err.Error()),
			)
			continue
		}
		if released {
			expiredCount++
		} else {
			skippedCount++
		}
	}

	// ── Batch metrics log ────────────────────────────────────────────────────
	slog.InfoContext(ctx, "sweeper tick complete",
		slog.Int("candidates", len(candidates)),
		slog.Int("released", expiredCount),
		slog.Int("skipped", skippedCount),
		slog.Int("errored", errCount),
	)
	// TODO: emit Prometheus counters reservations_swept_total{outcome=released/skipped/errored}
	// and gauge reservations_pending_expiry via SELECT COUNT(*) from the partial index.
}

// processCandidate processes a single expired reservation candidate in its own tx.
// Returns (true, nil) if the row was transitioned to EXPIRED.
// Returns (false, nil) if the row was already moved by another path (COMMITTED/RELEASED/EXPIRED).
// Returns (false, err) on infrastructure error.
func (s *Sweeper) processCandidate(ctx context.Context, c *Reservation) (released bool, err error) {
	tx, err := s.svc.Pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// ── LOCK-ORDER PIN: stock_levels(sku) BEFORE reservations(id) ────────────
	// "ALWAYS lock stock_levels(sku) BEFORE reservations(id) inside any
	//  state-mutating tx." — td.json concurrency_model.lock_ordering
	stockRow, err := s.svc.StockStorage.GetBySKUForUpdate(ctx, tx, c.SKU)
	if err != nil {
		return false, err
	}

	resvRow, err := s.svc.ResvStorage.GetByIDForUpdate(ctx, tx, c.ID)
	if err != nil {
		return false, err
	}

	// Re-check status: someone else may have committed or released this row
	// between the candidate scan and this per-row lock.
	if resvRow.Status != StatusReserved {
		// Already moved — skip; no stock change needed.
		_ = tx.Rollback(ctx)
		return false, nil
	}

	// Defensive: log if stock would go negative (should never happen; CHECK constraint also guards).
	if stockRow.ReservedQty < resvRow.Qty {
		slog.ErrorContext(ctx, "sweeper: reserved_qty would go negative (defensive panic-class)",
			slog.String("sku", resvRow.SKU),
			slog.Int("reserved_qty", stockRow.ReservedQty),
			slog.Int("r_qty", resvRow.Qty),
		)
	}

	// Transition stock: reserved → available.
	if err = s.svc.StockStorage.IncrementAvailableDecrementReserved(ctx, tx, resvRow.SKU, resvRow.Qty); err != nil {
		return false, err
	}

	// Transition reservation: RESERVED → EXPIRED.
	if err = s.svc.ResvStorage.MarkExpired(ctx, tx, resvRow.ID); err != nil {
		return false, err
	}

	// Insert outbox event for events.reservation.expired.
	// aggregate_id = orderId (Kafka partition key per contracts.json).
	now := time.Now().UTC()
	payload := ReservationExpiredPayload{
		EventID:       uuid.New().String(),
		OccurredAt:    now,
		OrderID:       resvRow.OrderID.String(),
		ReservationID: resvRow.ID.String(),
		SweptAt:       now,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	outboxEvt := &OutboxEvent{
		ID:          uuid.New(),
		AggregateID: resvRow.OrderID.String(), // partition key = orderId per contracts.json
		EventType:   "reservation.expired",
		PayloadJSON: payloadJSON,
	}
	if err = s.svc.OutboxStorage.Insert(ctx, tx, outboxEvt); err != nil {
		return false, err
	}

	if err = tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}
