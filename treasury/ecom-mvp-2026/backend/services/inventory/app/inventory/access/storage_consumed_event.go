package access

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory"
)

// ConsumedEventStorage provides dedup insertion for Kafka event consumers.
// All consumer handlers insert into consumed_events inside the same tx that
// applies the event's side-effect. A PK conflict means the event was already
// processed — the handler should commit and ack without further action.
type ConsumedEventStorage struct {
	pool *pgxpool.Pool
}

var _ inventory.ConsumedStorage = (*ConsumedEventStorage)(nil)

// NewConsumedEventStorage constructs a ConsumedEventStorage.
func NewConsumedEventStorage(pool *pgxpool.Pool) *ConsumedEventStorage {
	return &ConsumedEventStorage{pool: pool}
}

// TryInsert attempts to insert a dedup row for (eventID, consumerName).
// Returns (true, nil) if the row was newly inserted (first delivery).
// Returns (false, nil) if the row already existed (duplicate delivery — caller should no-op and ack).
// Returns (false, err) on unexpected errors.
func (s *ConsumedEventStorage) TryInsert(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, consumerName string) (inserted bool, err error) {
	const q = `
		INSERT INTO inventory.consumed_events (event_id, consumer_name)
		VALUES ($1, $2)
		ON CONFLICT (event_id, consumer_name) DO NOTHING`

	cmd, err := tx.Exec(ctx, q, eventID, consumerName)
	if err != nil {
		return false, fmt.Errorf("ConsumedEventStorage.TryInsert exec: %w", err)
	}
	return cmd.RowsAffected() > 0, nil
}
