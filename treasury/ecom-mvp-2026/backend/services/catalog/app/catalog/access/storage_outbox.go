package access

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OutboxStorage handles outbox_events persistence and relay scanning.
type OutboxStorage struct {
	db *pgxpool.Pool
}

// NewOutboxStorage constructs an OutboxStorage.
func NewOutboxStorage(db *pgxpool.Pool) *OutboxStorage {
	return &OutboxStorage{db: db}
}

// InsertOutboxParams carries the values for a new outbox_events row.
type InsertOutboxParams struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	PayloadJSON []byte
}

// Insert writes one outbox event row within the caller-managed transaction q.
// published_at is NULL, marking the row as pending.
func (s *OutboxStorage) Insert(ctx context.Context, q Querier, p InsertOutboxParams) error {
	const query = `
		INSERT INTO catalog.outbox_events
			(id, aggregate_id, event_type, payload_json, created_at, published_at)
		VALUES
			($1, $2, $3, $4, NOW(), NULL)`

	if _, err := q.Exec(ctx, query, p.ID, p.AggregateID, p.EventType, p.PayloadJSON); err != nil {
		return fmt.Errorf("storage: insert outbox event: %w", err)
	}
	return nil
}

// OutboxRow is a pending event record fetched by the relay goroutine.
type OutboxRow struct {
	ID          uuid.UUID
	PayloadJSON []byte
	EventType   string
}

// FetchPending selects up to limit unprocessed outbox rows using
// SELECT ... FOR UPDATE SKIP LOCKED to prevent concurrent relay instances
// from processing the same row.
func (s *OutboxStorage) FetchPending(ctx context.Context, limit int) ([]OutboxRow, error) {
	const query = `
		SELECT id, event_type, payload_json
		FROM catalog.outbox_events
		WHERE published_at IS NULL
		ORDER BY id ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("storage: fetch pending outbox: %w", err)
	}
	defer rows.Close()

	var result []OutboxRow
	for rows.Next() {
		var r OutboxRow
		if err := rows.Scan(&r.ID, &r.EventType, &r.PayloadJSON); err != nil {
			return nil, fmt.Errorf("storage: scan outbox row: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

// MarkPublished sets published_at = NOW() for the given IDs.
func (s *OutboxStorage) MarkPublished(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	const query = `
		UPDATE catalog.outbox_events
		SET published_at = $1
		WHERE id = ANY($2::uuid[])`

	_, err := s.db.Exec(ctx, query, time.Now().UTC(), ids)
	if err != nil {
		return fmt.Errorf("storage: mark outbox published: %w", err)
	}
	return nil
}
