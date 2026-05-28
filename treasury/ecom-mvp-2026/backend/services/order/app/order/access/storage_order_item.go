package access

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order"
)

// OrderItemStorage implements order.StorageOrderItem against PostgreSQL.
// ORD-007: NO Update method is exposed. The schema also enforces this via a
// defensive trigger on UPDATE (deferred to post-MVP; repo discipline covers MVP).
type OrderItemStorage struct {
	pool *pgxpool.Pool
}

func NewOrderItemStorage(pool *pgxpool.Pool) *OrderItemStorage {
	return &OrderItemStorage{pool: pool}
}

// BulkInsert persists all order_items for a new order in a single batch.
// Called INSIDE the create-from-checkout transaction.
// Snapshot fields (name, image, price) are frozen here — never updated.
func (s *OrderItemStorage) BulkInsert(ctx context.Context, tx pgx.Tx, items []order.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	// Build a multi-row VALUES clause for a single INSERT statement.
	const colCount = 10
	placeholders := make([]string, 0, len(items))
	args := make([]any, 0, len(items)*colCount)

	for i, it := range items {
		base := i * colCount
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9, base+10,
		))
		args = append(args,
			it.ID,
			it.OrderID,
			it.ProductID,
			it.SKU,
			it.Qty,
			it.NameSnapshot,
			it.ImageURLSnapshot,
			it.PriceSnapshot,
			it.LineSubtotal,
			it.CreatedAt,
		)
	}

	q := fmt.Sprintf(`
		INSERT INTO "order".order_items
			(id, order_id, product_id, sku, qty, name_snapshot, image_url_snapshot, price_snapshot, line_subtotal, created_at)
		VALUES %s`,
		strings.Join(placeholders, ", "),
	)

	if _, err := tx.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("bulk insert order items: %w", err)
	}
	return nil
}

// ListByOrderID fetches all line items for an order (used by order.detail).
func (s *OrderItemStorage) ListByOrderID(ctx context.Context, orderID string) ([]order.OrderItem, error) {
	const q = `
		SELECT id, order_id, product_id, sku, qty,
		       name_snapshot, image_url_snapshot, price_snapshot, line_subtotal, created_at
		FROM "order".order_items
		WHERE order_id = $1
		ORDER BY created_at ASC`

	rows, err := s.pool.Query(ctx, q, orderID)
	if err != nil {
		return nil, fmt.Errorf("list order items: %w", err)
	}
	defer rows.Close()

	items := make([]order.OrderItem, 0)
	for rows.Next() {
		var it order.OrderItem
		if err := rows.Scan(
			&it.ID, &it.OrderID, &it.ProductID, &it.SKU, &it.Qty,
			&it.NameSnapshot, &it.ImageURLSnapshot, &it.PriceSnapshot, &it.LineSubtotal, &it.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("order items rows: %w", err)
	}
	return items, nil
}
