package access

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/app/payment"
)

// intentStorage implements payment.IntentStorage.

type intentStorage struct {
	db *pgxpool.Pool
}

var _ payment.IntentStorage = (*intentStorage)(nil)

// NewIntentStorage constructs the concrete IntentStorage backed by Postgres.
func NewIntentStorage(db *pgxpool.Pool) payment.IntentStorage {
	return &intentStorage{db: db}
}

func (s *intentStorage) Insert(ctx context.Context, pi payment.PaymentIntent) error {
	const q = `
		INSERT INTO payment.payment_intents
			(intent_id, order_id, owner_user_id, amount_minor, currency,
			 status, expires_at, created_at, updated_at, version)
		VALUES ($1,$2,$3,$4,$5,$6,$7,now(),now(),0)
		ON CONFLICT (order_id) DO NOTHING`

	tag, err := s.db.Exec(ctx, q,
		pi.IntentID, pi.OrderID, pi.OwnerUserID,
		pi.AmountMinor, pi.Currency, pi.Status, pi.ExpiresAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return payment.ErrDuplicateOrderID
	}
	return nil
}

func (s *intentStorage) GetByOrderID(ctx context.Context, orderID uuid.UUID) (payment.PaymentIntent, error) {
	const q = `
		SELECT intent_id, order_id, owner_user_id, amount_minor, currency,
		       status, mock_provider_ref, provider_status, paid_at,
		       expires_at, created_at, updated_at, version
		FROM payment.payment_intents
		WHERE order_id = $1`

	return scanIntent(s.db.QueryRow(ctx, q, orderID))
}

func (s *intentStorage) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, intentID uuid.UUID) (payment.PaymentIntent, error) {
	const q = `
		SELECT intent_id, order_id, owner_user_id, amount_minor, currency,
		       status, mock_provider_ref, provider_status, paid_at,
		       expires_at, created_at, updated_at, version
		FROM payment.payment_intents
		WHERE intent_id = $1
		FOR UPDATE`

	return scanIntent(tx.QueryRow(ctx, q, intentID))
}

func (s *intentStorage) MarkSucceeded(ctx context.Context, tx pgx.Tx, intentID uuid.UUID, version int, providerRef string, paidAt time.Time) error {
	const q = `
		UPDATE payment.payment_intents
		SET status = 'SUCCEEDED',
		    mock_provider_ref = $3,
		    provider_status = 'SUCCEEDED',
		    paid_at = $4,
		    updated_at = now(),
		    version = version + 1
		WHERE intent_id = $1 AND version = $2`

	_, err := tx.Exec(ctx, q, intentID, version, providerRef, paidAt)
	return err
}

func (s *intentStorage) MarkFailed(ctx context.Context, tx pgx.Tx, intentID uuid.UUID, version int, providerRef string) error {
	const q = `
		UPDATE payment.payment_intents
		SET status = 'FAILED',
		    mock_provider_ref = $3,
		    provider_status = 'FAILED',
		    updated_at = now(),
		    version = version + 1
		WHERE intent_id = $1 AND version = $2`

	_, err := tx.Exec(ctx, q, intentID, version, providerRef)
	return err
}

func (s *intentStorage) MarkExpired(ctx context.Context, tx pgx.Tx, intentID uuid.UUID, version int) error {
	const q = `
		UPDATE payment.payment_intents
		SET status = 'EXPIRED',
		    updated_at = now(),
		    version = version + 1
		WHERE intent_id = $1 AND version = $2 AND status = 'REQUIRES_PAYMENT'`

	_, err := tx.Exec(ctx, q, intentID, version)
	return err
}

func (s *intentStorage) ListExpiredIDs(ctx context.Context, tx pgx.Tx, batchSize int) ([]uuid.UUID, error) {
	const q = `
		SELECT intent_id
		FROM payment.payment_intents
		WHERE status = 'REQUIRES_PAYMENT'
		  AND expires_at < now()
		ORDER BY expires_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := tx.Query(ctx, q, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// scanIntent reads a PaymentIntent from a pgx.Row.
func scanIntent(row pgx.Row) (payment.PaymentIntent, error) {
	var pi payment.PaymentIntent
	err := row.Scan(
		&pi.IntentID,
		&pi.OrderID,
		&pi.OwnerUserID,
		&pi.AmountMinor,
		&pi.Currency,
		&pi.Status,
		&pi.MockProviderRef,
		&pi.ProviderStatus,
		&pi.PaidAt,
		&pi.ExpiresAt,
		&pi.CreatedAt,
		&pi.UpdatedAt,
		&pi.Version,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return payment.PaymentIntent{}, pgx.ErrNoRows
	}
	return pi, err
}
