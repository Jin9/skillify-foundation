package access

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order"
)

// OutboxStorage implements order.StorageOutbox against PostgreSQL.
// Outbox rows are written in the SAME transaction as the status mutation and
// published asynchronously by the outbox publisher goroutine.
type OutboxStorage struct {
	pool *pgxpool.Pool
}

func NewOutboxStorage(pool *pgxpool.Pool) *OutboxStorage {
	return &OutboxStorage{pool: pool}
}

// Insert writes a new outbox row inside the given transaction.
// aggregate_id becomes the Kafka partition key (events.order.cancelled.ordering_guarantees).
func (s *OutboxStorage) Insert(ctx context.Context, tx pgx.Tx, event order.OutboxEvent) error {
	const q = `
		INSERT INTO "order".outbox_events (id, aggregate_id, event_type, payload_json, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	if _, err := tx.Exec(ctx, q,
		event.ID,
		event.AggregateID,
		event.EventType,
		event.PayloadJSON,
		event.CreatedAt,
	); err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}
	return nil
}

// MarkPublished sets published_at to NOW() for the given outbox row.
// Called by the outbox publisher goroutine after the Kafka send succeeds.
func (s *OutboxStorage) MarkPublished(ctx context.Context, id string) error {
	const q = `UPDATE "order".outbox_events SET published_at = now() WHERE id = $1`
	if _, err := s.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("mark outbox published: %w", err)
	}
	return nil
}

// PollUnpublished returns up to limit outbox rows where published_at IS NULL.
// Uses the partial index idx_outbox_unpublished for efficiency.
// Called by the outbox publisher goroutine on a 1-second ticker.
func (s *OutboxStorage) PollUnpublished(ctx context.Context, limit int) ([]order.OutboxEvent, error) {
	const q = `
		SELECT id, aggregate_id, event_type, payload_json, created_at, published_at
		FROM "order".outbox_events
		WHERE published_at IS NULL
		ORDER BY id ASC
		LIMIT $1`

	rows, err := s.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("poll unpublished outbox: %w", err)
	}
	defer rows.Close()

	events := make([]order.OutboxEvent, 0)
	for rows.Next() {
		var e order.OutboxEvent
		if err := rows.Scan(&e.ID, &e.AggregateID, &e.EventType, &e.PayloadJSON, &e.CreatedAt, &e.PublishedAt); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("outbox rows: %w", err)
	}
	return events, nil
}
