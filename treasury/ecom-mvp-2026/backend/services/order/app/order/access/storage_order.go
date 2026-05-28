package access

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order"
)

// OrderStorage implements order.StorageOrder against PostgreSQL (pgx).
// Schema: "order". All status writes go through UpdateStatus — the ONLY place
// that touches orders.status — enforcing the sole-writer invariant.
type OrderStorage struct {
	pool *pgxpool.Pool
}

func NewOrderStorage(pool *pgxpool.Pool) *OrderStorage {
	return &OrderStorage{pool: pool}
}

// NextOrderNumber allocates the next order number from the global sequence and
// formats it as ORD-YYYYMMDD-NNNNNN (td.json §order_number_strategy).
// Called INSIDE an open transaction so the sequence allocation is rolled back
// if the surrounding tx fails (Postgres sequences are NOT rollback-safe by default,
// but gaps are acceptable and preferable to duplicate keys).
func (s *OrderStorage) NextOrderNumber(ctx context.Context, tx pgx.Tx, datePart string) (string, error) {
	var seq int64
	err := tx.QueryRow(ctx, "SELECT nextval('order.order_number_seq')").Scan(&seq)
	if err != nil {
		return "", fmt.Errorf("next order number seq: %w", err)
	}
	return fmt.Sprintf("ORD-%s-%06d", datePart, seq), nil
}

// Insert persists a new order row. Called once per order lifecycle — no UPDATE path here.
func (s *OrderStorage) Insert(ctx context.Context, tx pgx.Tx, o order.Order) error {
	const q = `
		INSERT INTO "order".orders (
			id, order_number, user_id, status,
			subtotal, shipping_fee, coupon_discount, grand_total, currency,
			address_snapshot, buyer_email_snapshot,
			tracking_number, idempotency_key, version, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11,
			$12, $13, $14, $15, $16
		)`
	_, err := tx.Exec(ctx, q,
		o.ID, o.OrderNumber, o.UserID, string(o.Status),
		o.Subtotal, o.ShippingFee, o.CouponDiscount, o.GrandTotal, o.Currency,
		o.AddressSnapshot, o.BuyerEmailSnapshot,
		o.TrackingNumber, o.IdempotencyKey, o.Version, o.CreatedAt, o.UpdatedAt,
	)
	return err
}

// GetByIDForUpdate fetches an order row with SELECT ... FOR UPDATE.
// Used by every state-changing path (consumers, cancel endpoints).
// Returns nil, nil when the order does not exist.
func (s *OrderStorage) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id string) (*order.Order, error) {
	const q = `
		SELECT id, order_number, user_id, status,
		       subtotal, shipping_fee, coupon_discount, grand_total, currency,
		       address_snapshot, buyer_email_snapshot,
		       tracking_number, idempotency_key, version, created_at, updated_at
		FROM "order".orders
		WHERE id = $1
		FOR UPDATE`
	o, err := scanOrder(tx.QueryRow(ctx, q, id))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return o, err
}

// GetByIDScopedForUpdate fetches an order row belonging to userID with SELECT ... FOR UPDATE.
// Used by order.cancel-mine (combined ownership + lock, per td.json cancel-mine tx step 2).
// Returns nil, nil when the order does not exist OR belongs to a different user
// (enumeration resistance: callers must return 404 ORDER_NOT_OWNED in both cases).
func (s *OrderStorage) GetByIDScopedForUpdate(ctx context.Context, tx pgx.Tx, id, userID string) (*order.Order, error) {
	const q = `
		SELECT id, order_number, user_id, status,
		       subtotal, shipping_fee, coupon_discount, grand_total, currency,
		       address_snapshot, buyer_email_snapshot,
		       tracking_number, idempotency_key, version, created_at, updated_at
		FROM "order".orders
		WHERE id = $1 AND user_id = $2
		FOR UPDATE`
	o, err := scanOrder(tx.QueryRow(ctx, q, id, userID))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return o, err
}

// GetByID fetches an order row without a lock (read-only path: order.detail for ADMIN).
// Returns nil, nil when the order does not exist.
func (s *OrderStorage) GetByID(ctx context.Context, id string) (*order.Order, error) {
	const q = `
		SELECT id, order_number, user_id, status,
		       subtotal, shipping_fee, coupon_discount, grand_total, currency,
		       address_snapshot, buyer_email_snapshot,
		       tracking_number, idempotency_key, version, created_at, updated_at
		FROM "order".orders
		WHERE id = $1`
	o, err := scanOrder(s.pool.QueryRow(ctx, q, id))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return o, err
}

// UpdateStatus is the SOLE METHOD that writes orders.status.
// Called within an active transaction after SELECT ... FOR UPDATE.
// Bumps version (optimistic lock belt-and-suspenders) and updated_at.
// trackingNumber is set only on PACKING→SHIPPED transitions (ORD-006).
func (s *OrderStorage) UpdateStatus(ctx context.Context, tx pgx.Tx, id string, newStatus order.OrderStatus, version int, trackingNumber *string) error {
	const q = `
		UPDATE "order".orders
		SET status = $1,
		    version = version + 1,
		    updated_at = now(),
		    tracking_number = COALESCE($2, tracking_number)
		WHERE id = $3 AND version = $4`
	tag, err := tx.Exec(ctx, q, string(newStatus), trackingNumber, id, version)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update order status: version conflict or order not found (id=%s, version=%d)", id, version)
	}
	return nil
}

// ListByUserID returns a paginated (DESC createdAt) list of orders belonging to userID.
// Returns the rows, the total count (for pagination), and any error.
// The optional status filter is applied when non-nil.
func (s *OrderStorage) ListByUserID(ctx context.Context, userID string, status *order.OrderStatus, limit, offset int) ([]order.Order, int, error) {
	args := []any{userID}
	filterClause := ""
	if status != nil {
		args = append(args, string(*status))
		filterClause = fmt.Sprintf("AND status = $%d", len(args))
	}

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM "order".orders WHERE user_id = $1 %s`, filterClause)
	var total int
	if err := s.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("list orders count: %w", err)
	}

	args = append(args, limit, offset)
	listQ := fmt.Sprintf(`
		SELECT id, order_number, user_id, status,
		       subtotal, shipping_fee, coupon_discount, grand_total, currency,
		       address_snapshot, buyer_email_snapshot,
		       tracking_number, idempotency_key, version, created_at, updated_at
		FROM "order".orders
		WHERE user_id = $1 %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		filterClause, len(args)-1, len(args),
	)

	rows, err := s.pool.Query(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders query: %w", err)
	}
	defer rows.Close()

	orders := make([]order.Order, 0)
	for rows.Next() {
		o, err := scanOrderRow(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("list orders scan: %w", err)
		}
		orders = append(orders, *o)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list orders rows: %w", err)
	}

	return orders, total, nil
}

// scanOrder reads a single order row from a pgx.Row.
func scanOrder(row pgx.Row) (*order.Order, error) {
	var o order.Order
	var statusStr string
	err := row.Scan(
		&o.ID, &o.OrderNumber, &o.UserID, &statusStr,
		&o.Subtotal, &o.ShippingFee, &o.CouponDiscount, &o.GrandTotal, &o.Currency,
		&o.AddressSnapshot, &o.BuyerEmailSnapshot,
		&o.TrackingNumber, &o.IdempotencyKey, &o.Version, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.Status = order.OrderStatus(statusStr)
	return &o, nil
}

// scanOrderRow reads a single order row from pgx.Rows (multi-row query path).
func scanOrderRow(rows pgx.Rows) (*order.Order, error) {
	var o order.Order
	var statusStr string
	err := rows.Scan(
		&o.ID, &o.OrderNumber, &o.UserID, &statusStr,
		&o.Subtotal, &o.ShippingFee, &o.CouponDiscount, &o.GrandTotal, &o.Currency,
		&o.AddressSnapshot, &o.BuyerEmailSnapshot,
		&o.TrackingNumber, &o.IdempotencyKey, &o.Version, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	o.Status = order.OrderStatus(statusStr)
	return &o, nil
}
