package access

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory"
)

// OutboxStorage provides data access for the outbox_events table.
type OutboxStorage struct {
	pool *pgxpool.Pool
}

var _ inventory.OutboxStorage = (*OutboxStorage)(nil)

// NewOutboxStorage constructs an OutboxStorage backed by the given pool.
func NewOutboxStorage(pool *pgxpool.Pool) *OutboxStorage {
	return &OutboxStorage{pool: pool}
}

// Insert writes an outbox row inside the given transaction.
// Called in the same tx as the business mutation (sweeper EXPIRED transition).
// At-least-once: the relayer goroutine will publish and set published_at.
func (s *OutboxStorage) Insert(ctx context.Context, tx pgx.Tx, evt *inventory.OutboxEvent) error {
	const q = `
		INSERT INTO inventory.outbox_events
		    (id, aggregate_id, event_type, payload_json)
		VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(ctx, q, evt.ID, evt.AggregateID, evt.EventType, evt.PayloadJSON)
	if err != nil {
		return fmt.Errorf("OutboxStorage.Insert exec: %w", err)
	}
	return nil
}

// FetchUnpublished fetches up to limit rows where published_at IS NULL, ordered by id,
// using FOR UPDATE SKIP LOCKED to allow concurrent relay goroutines on multiple replicas.
func (s *OutboxStorage) FetchUnpublished(ctx context.Context, tx pgx.Tx, limit int) ([]*inventory.OutboxEvent, error) {
	const q = `
		SELECT id, aggregate_id, event_type, payload_json, created_at, published_at
		FROM inventory.outbox_events
		WHERE published_at IS NULL
		ORDER BY id
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := tx.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("FetchUnpublished query: %w", err)
	}
	defer rows.Close()

	var events []*inventory.OutboxEvent
	for rows.Next() {
		var e inventory.OutboxEvent
		if err := rows.Scan(&e.ID, &e.AggregateID, &e.EventType, &e.PayloadJSON, &e.CreatedAt, &e.PublishedAt); err != nil {
			return nil, fmt.Errorf("FetchUnpublished scan: %w", err)
		}
		events = append(events, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FetchUnpublished rows: %w", err)
	}
	return events, nil
}

// MarkPublished sets published_at = NOW() for a single outbox row.
// Called by the relay goroutine after a successful Kafka publish.
func (s *OutboxStorage) MarkPublished(ctx context.Context, tx pgx.Tx, id interface{}) error {
	now := time.Now().UTC()
	const q = `
		UPDATE inventory.outbox_events
		SET published_at = $2
		WHERE id = $1`

	cmd, err := tx.Exec(ctx, q, id, now)
	if err != nil {
		return fmt.Errorf("MarkPublished exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("MarkPublished: no row for id=%v", id)
	}
	return nil
}
