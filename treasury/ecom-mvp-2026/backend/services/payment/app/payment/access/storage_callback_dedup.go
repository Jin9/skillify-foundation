package access

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/app/payment"
)

type callbackDedupStorage struct {
	db *pgxpool.Pool
}

var _ payment.CallbackDedupStorage = (*callbackDedupStorage)(nil)

// NewCallbackDedupStorage constructs the concrete CallbackDedupStorage.
func NewCallbackDedupStorage(db *pgxpool.Pool) payment.CallbackDedupStorage {
	return &callbackDedupStorage{db: db}
}

func (s *callbackDedupStorage) LookupTx(ctx context.Context, tx pgx.Tx, dedupKey string) (payment.CallbackDedup, error) {
	const q = `
		SELECT dedup_key, intent_id, provider_status, envelope, http_status, created_at, expires_at
		FROM payment.payment_callback_dedup
		WHERE dedup_key = $1`

	var d payment.CallbackDedup
	err := tx.QueryRow(ctx, q, dedupKey).Scan(
		&d.DedupKey,
		&d.IntentID,
		&d.ProviderStatus,
		&d.Envelope,
		&d.HTTPStatus,
		&d.CreatedAt,
		&d.ExpiresAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return payment.CallbackDedup{}, pgx.ErrNoRows
	}
	return d, err
}

func (s *callbackDedupStorage) Insert(ctx context.Context, tx pgx.Tx, row payment.CallbackDedup) error {
	const q = `
		INSERT INTO payment.payment_callback_dedup
			(dedup_key, intent_id, provider_status, envelope, http_status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, now(), $6)
		ON CONFLICT (dedup_key) DO NOTHING`

	expiresAt := row.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = time.Now().UTC().AddDate(0, 0, 30)
	}

	_, err := tx.Exec(ctx, q,
		row.DedupKey,
		row.IntentID,
		row.ProviderStatus,
		row.Envelope,
		row.HTTPStatus,
		expiresAt,
	)
	return err
}
