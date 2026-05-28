package access

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/cart/app/cart"
)

// CartStorage handles all cart + cart_items DB operations.
// All methods are scoped by user_id — no cross-customer leakage is possible.
type CartStorage struct {
	pool *pgxpool.Pool
}

// Compile-time check: CartStorage must satisfy cart.CartStoragePort.
var _ cart.CartStoragePort = (*CartStorage)(nil)

func NewCartStorage(pool *pgxpool.Pool) *CartStorage {
	return &CartStorage{pool: pool}
}

// GetOrCreateByUserID returns the existing cart for userID, creating one if absent.
// The cart row is inserted via ON CONFLICT DO NOTHING so concurrent calls are safe.
func (s *CartStorage) GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (cart.Cart, error) {
	cartID := uuid.New()
	now := time.Now().UTC()

	// Try to insert; ignore conflict if the cart already exists.
	const upsertSQL = `
		INSERT INTO cart.carts (cart_id, user_id, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING`
	_, err := s.pool.Exec(ctx, upsertSQL, cartID, userID, now)
	if err != nil {
		return cart.Cart{}, err
	}

	// Always read back the authoritative row (whether newly inserted or pre-existing).
	const selectSQL = `SELECT cart_id, user_id, updated_at FROM cart.carts WHERE user_id = $1`
	row := s.pool.QueryRow(ctx, selectSQL, userID)

	var c cart.Cart
	if err := row.Scan(&c.CartID, &c.UserID, &c.UpdatedAt); err != nil {
		return cart.Cart{}, err
	}
	return c, nil
}

// UpsertItem merges a product into the cart. If the (cart_id, product_id) pair already
// exists, qty is added to the existing quantity (CART-006 ON CONFLICT merge invariant).
// The caller must supply a transaction so this and TouchUpdatedAt are atomic.
func (s *CartStorage) UpsertItem(ctx context.Context, tx pgx.Tx, cartID, productID uuid.UUID, qty int) error {
	const sql = `
		INSERT INTO cart.cart_items (cart_item_id, cart_id, product_id, qty, added_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET
			qty      = cart.cart_items.qty + EXCLUDED.qty,
			added_at = NOW()`
	_, err := tx.Exec(ctx, sql, uuid.New(), cartID, productID, qty)
	return err
}

// UpdateItemQty replaces the qty on an existing cart item (last-writer-wins).
// Returns pgx.ErrNoRows if the item does not belong to this cart.
func (s *CartStorage) UpdateItemQty(ctx context.Context, tx pgx.Tx, cartItemID, cartID uuid.UUID, qty int) (int64, error) {
	const sql = `
		UPDATE cart.cart_items
		SET qty = $1, added_at = NOW()
		WHERE cart_item_id = $2 AND cart_id = $3`
	tag, err := tx.Exec(ctx, sql, qty, cartItemID, cartID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// DeleteItem removes a single cart item. Returns rowsAffected (0 if already gone — idempotent).
func (s *CartStorage) DeleteItem(ctx context.Context, tx pgx.Tx, cartItemID, cartID uuid.UUID) (int64, error) {
	const sql = `DELETE FROM cart.cart_items WHERE cart_item_id = $1 AND cart_id = $2`
	tag, err := tx.Exec(ctx, sql, cartItemID, cartID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// DeleteItems removes a batch of cart items by ID (for cart.clear-on-checkout).
// Uses pgx array parameter — no string concatenation.
func (s *CartStorage) DeleteItems(ctx context.Context, tx pgx.Tx, cartID uuid.UUID, itemIDs []uuid.UUID) (int64, error) {
	const sql = `DELETE FROM cart.cart_items WHERE cart_id = $1 AND cart_item_id = ANY($2::uuid[])`
	tag, err := tx.Exec(ctx, sql, cartID, itemIDs)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// GetCartWithItems fetches the cart row plus all its items for the given user.
// Returns (zero Cart, nil items, nil) when no cart exists for the user.
func (s *CartStorage) GetCartWithItems(ctx context.Context, userID uuid.UUID) (cart.Cart, []cart.CartItem, error) {
	const cartSQL = `SELECT cart_id, user_id, updated_at FROM cart.carts WHERE user_id = $1`
	row := s.pool.QueryRow(ctx, cartSQL, userID)

	var c cart.Cart
	if err := row.Scan(&c.CartID, &c.UserID, &c.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return cart.Cart{}, nil, nil
		}
		return cart.Cart{}, nil, err
	}

	const itemsSQL = `
		SELECT cart_item_id, cart_id, product_id, qty, added_at
		FROM cart.cart_items
		WHERE cart_id = $1
		ORDER BY added_at`
	rows, err := s.pool.Query(ctx, itemsSQL, c.CartID)
	if err != nil {
		return cart.Cart{}, nil, err
	}
	defer rows.Close()

	var items []cart.CartItem
	for rows.Next() {
		var item cart.CartItem
		if err := rows.Scan(&item.CartItemID, &item.CartID, &item.ProductID, &item.Qty, &item.AddedAt); err != nil {
			return cart.Cart{}, nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return cart.Cart{}, nil, err
	}

	return c, items, nil
}

// GetCartByUserID returns the cart row only (no items), used internally for mutation paths.
func (s *CartStorage) GetCartByUserID(ctx context.Context, userID uuid.UUID) (cart.Cart, error) {
	const sql = `SELECT cart_id, user_id, updated_at FROM cart.carts WHERE user_id = $1`
	row := s.pool.QueryRow(ctx, sql, userID)

	var c cart.Cart
	if err := row.Scan(&c.CartID, &c.UserID, &c.UpdatedAt); err != nil {
		return cart.Cart{}, err
	}
	return c, nil
}

// GetItemByID fetches a single cart_item row, scoped to cartID.
func (s *CartStorage) GetItemByID(ctx context.Context, cartItemID, cartID uuid.UUID) (cart.CartItem, error) {
	const sql = `
		SELECT cart_item_id, cart_id, product_id, qty, added_at
		FROM cart.cart_items
		WHERE cart_item_id = $1 AND cart_id = $2`
	row := s.pool.QueryRow(ctx, sql, cartItemID, cartID)

	var item cart.CartItem
	if err := row.Scan(&item.CartItemID, &item.CartID, &item.ProductID, &item.Qty, &item.AddedAt); err != nil {
		return cart.CartItem{}, err
	}
	return item, nil
}

// TouchUpdatedAt bumps the cart's updated_at timestamp. Must be called inside a transaction.
func (s *CartStorage) TouchUpdatedAt(ctx context.Context, tx pgx.Tx, cartID uuid.UUID) error {
	const sql = `UPDATE cart.carts SET updated_at = NOW() WHERE cart_id = $1`
	_, err := tx.Exec(ctx, sql, cartID)
	return err
}

// BeginTx starts a new transaction on the pool.
func (s *CartStorage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.pool.Begin(ctx)
}
