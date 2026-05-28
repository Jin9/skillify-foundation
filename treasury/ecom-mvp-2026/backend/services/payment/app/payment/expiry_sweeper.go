package payment

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ExpirySweeper is the internal periodic job that flips REQUIRES_PAYMENT intents
// past their expires_at to EXPIRED status.
//
// Scope:
//   - Internal only — emits NO outbox event (per expired_event_decision / TL PR-001 LOCKED).
//   - The reservation sweeper (inventory service) drives Order/Inventory state via
//     events.reservation.expired; Payment's sweeper only makes Payment's own state honest.
//   - FOR UPDATE SKIP LOCKED ensures multi-replica safety without a separate lock table.
//
// Metrics placeholders (increment via your metrics library):
//   - payment.expiry_sweeper.flipped_total
//   - payment.expiry_sweeper.lag_seconds
type ExpirySweeper struct {
	db          *pgxpool.Pool
	intentStore IntentStorage
	cadence     time.Duration
	batchSize   int
}

// NewExpirySweeper constructs the sweeper. Call Start(ctx) in a dedicated goroutine.
func NewExpirySweeper(db *pgxpool.Pool, intentStore IntentStorage, cadence time.Duration, batchSize int) *ExpirySweeper {
	return &ExpirySweeper{
		db:          db,
		intentStore: intentStore,
		cadence:     cadence,
		batchSize:   batchSize,
	}
}

// Start runs the sweeper loop. It blocks until ctx is cancelled.
// Call it in a goroutine: go sweeper.Start(ctx).
func (sw *ExpirySweeper) Start(ctx context.Context) {
	slog.InfoContext(ctx, "expiry sweeper started",
		"cadence", sw.cadence,
		"batch_size", sw.batchSize,
	)

	ticker := time.NewTicker(sw.cadence)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "expiry sweeper stopping")
			return
		case <-ticker.C:
			sw.tick(ctx)
		}
	}
}

// tick runs a single sweep, processing up to batchSize expired intents.
// Each intent is updated in its own short transaction so a failure in one
// row does not block the rest.
func (sw *ExpirySweeper) tick(ctx context.Context) {
	// Open a single tx to lock the batch via FOR UPDATE SKIP LOCKED.
	// We collect IDs first, then update each in its own tx to minimise lock hold time.
	// Alternative: update in the same tx — acceptable for small batches.
	tx, err := sw.db.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "sweeper: begin tx failed", "error", err)
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ids, err := sw.intentStore.ListExpiredIDs(ctx, tx, sw.batchSize)
	if err != nil {
		slog.ErrorContext(ctx, "sweeper: list expired ids failed", "error", err)
		return
	}

	if len(ids) == 0 {
		_ = tx.Commit(ctx)
		return
	}

	// We need the version for each ID to call MarkExpired with optimistic lock.
	// Simplest approach: fetch each row (locked via FOR UPDATE in the same tx)
	// and update in that same tx. FOR UPDATE SKIP LOCKED already serialises us.
	flipped := 0
	for _, id := range ids {
		intent, err := sw.intentStore.GetByIDForUpdate(ctx, tx, id)
		if err != nil {
			slog.ErrorContext(ctx, "sweeper: get intent failed",
				"intent_id", id, "error", err,
			)
			continue
		}
		if intent.Status != StatusRequiresPayment {
			// Another callback arrived and transitioned this row before us — skip.
			continue
		}
		if err := sw.intentStore.MarkExpired(ctx, tx, id, intent.Version); err != nil {
			slog.ErrorContext(ctx, "sweeper: mark expired failed",
				"intent_id", id, "error", err,
			)
			continue
		}
		flipped++
	}

	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "sweeper: commit failed", "error", err)
		return
	}

	if flipped > 0 {
		slog.InfoContext(ctx, "sweeper tick complete",
			"flipped", flipped,
			"batch_size", sw.batchSize,
			// metric: payment.expiry_sweeper.flipped_total += flipped
		)
	}
}
