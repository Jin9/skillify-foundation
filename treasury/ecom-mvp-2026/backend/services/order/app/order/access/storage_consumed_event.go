package access

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConsumedEventStorage implements order.StorageConsumedEvent against PostgreSQL.
// Provides per-service idempotency for inbound Kafka events (layer 1 of two-layer dedup).
type ConsumedEventStorage struct {
	pool *pgxpool.Pool
}

func NewConsumedEventStorage(pool *pgxpool.Pool) *ConsumedEventStorage {
	return &ConsumedEventStorage{pool: pool}
}

// InsertOrConflict attempts to insert a consumed_events row.
// Returns (true, nil) when the row was newly inserted (first processing).
// Returns (false, nil) when the event_id PK conflicts (duplicate/redelivery — ack silently).
// Any other error is returned as-is.
//
// Must be called INSIDE an open transaction so the dedup row and the status update
// commit or rollback atomically.
func (s *ConsumedEventStorage) InsertOrConflict(
	ctx context.Context,
	tx pgx.Tx,
	eventID, consumerName, eventType, orderID string,
) (inserted bool, err error) {
	const q = `
		INSERT INTO "order".consumed_events (event_id, consumer_name, event_type, order_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (event_id) DO NOTHING
		RETURNING event_id`

	var returnedID string
	err = tx.QueryRow(ctx, q, eventID, consumerName, eventType, orderID).Scan(&returnedID)
	if err == pgx.ErrNoRows {
		// ON CONFLICT DO NOTHING — row already existed, no RETURNING row emitted.
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("insert consumed event: %w", err)
	}
	return true, nil
}
