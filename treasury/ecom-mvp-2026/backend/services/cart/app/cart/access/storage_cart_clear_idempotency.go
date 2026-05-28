package access

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/cart/app/cart"
)

// GetClearIdempotency checks whether orderId was already processed.
// Returns (record, nil) on hit, (zero, pgx.ErrNoRows) on miss.
func (s *CartStorage) GetClearIdempotency(ctx context.Context, orderID uuid.UUID) (cart.ClearIdempotencyRecord, error) {
	const sql = `SELECT order_id, removed FROM cart.cart_clear_idempotency WHERE order_id = $1`
	row := s.pool.QueryRow(ctx, sql, orderID)

	var rec cart.ClearIdempotencyRecord
	if err := row.Scan(&rec.OrderID, &rec.Removed); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return cart.ClearIdempotencyRecord{}, pgx.ErrNoRows
		}
		return cart.ClearIdempotencyRecord{}, err
	}
	return rec, nil
}

// SaveClearIdempotency persists the cleared count keyed by orderId.
// ON CONFLICT DO NOTHING guarantees safety on concurrent replay.
// Must be called inside a transaction.
func (s *CartStorage) SaveClearIdempotency(ctx context.Context, tx pgx.Tx, orderID uuid.UUID, removed int) error {
	const sql = `
		INSERT INTO cart.cart_clear_idempotency (order_id, removed, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (order_id) DO NOTHING`
	_, err := tx.Exec(ctx, sql, orderID, removed)
	return err
}
