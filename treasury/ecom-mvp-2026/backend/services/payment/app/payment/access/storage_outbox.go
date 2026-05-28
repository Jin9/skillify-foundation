package access

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/app/payment"
)

type outboxStorage struct {
	db *pgxpool.Pool
}

var _ payment.OutboxStorage = (*outboxStorage)(nil)

// NewOutboxStorage constructs the concrete OutboxStorage.
func NewOutboxStorage(db *pgxpool.Pool) payment.OutboxStorage {
	return &outboxStorage{db: db}
}

func (s *outboxStorage) Insert(ctx context.Context, tx pgx.Tx, evt payment.OutboxEvent) error {
	const q = `
		INSERT INTO payment.outbox
			(id, aggregate_type, aggregate_id, event_id, event_type, topic,
			 partition_key, payload, status, created_at, attempts)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'PENDING',now(),0)`

	_, err := tx.Exec(ctx, q,
		evt.ID,
		evt.AggregateType,
		evt.AggregateID,
		evt.EventID,
		evt.EventType,
		evt.Topic,
		evt.PartitionKey,
		evt.Payload,
	)
	return err
}

func (s *outboxStorage) ListPending(ctx context.Context, limit int) ([]payment.OutboxEvent, error) {
	const q = `
		SELECT id, aggregate_type, aggregate_id, event_id, event_type, topic,
		       partition_key, payload, status, created_at, published_at, attempts, last_error
		FROM payment.outbox
		WHERE status = 'PENDING'
		ORDER BY created_at
		LIMIT $1`

	rows, err := s.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evts []payment.OutboxEvent
	for rows.Next() {
		var e payment.OutboxEvent
		if err := rows.Scan(
			&e.ID, &e.AggregateType, &e.AggregateID, &e.EventID,
			&e.EventType, &e.Topic, &e.PartitionKey, &e.Payload,
			&e.Status, &e.CreatedAt, &e.PublishedAt, &e.Attempts, &e.LastError,
		); err != nil {
			return nil, err
		}
		evts = append(evts, e)
	}
	return evts, rows.Err()
}

func (s *outboxStorage) MarkPublished(ctx context.Context, id string) error {
	const q = `
		UPDATE payment.outbox
		SET status = 'PUBLISHED', published_at = now()
		WHERE id = $1`
	_, err := s.db.Exec(ctx, q, id)
	return err
}

func (s *outboxStorage) MarkFailed(ctx context.Context, id, lastError string, maxAttempts int) error {
	const q = `
		UPDATE payment.outbox
		SET attempts   = attempts + 1,
		    last_error = $2,
		    status     = CASE WHEN attempts + 1 >= $3 THEN 'FAILED' ELSE status END
		WHERE id = $1`
	_, err := s.db.Exec(ctx, q, id, lastError, maxAttempts)
	return err
}
