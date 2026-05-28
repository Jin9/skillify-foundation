package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory"
)

// StockLevelStorage provides data access for the stock_levels table.
type StockLevelStorage struct {
	pool *pgxpool.Pool
}

var _ inventory.StockStorage = (*StockLevelStorage)(nil)

// NewStockLevelStorage constructs a StockLevelStorage backed by the given pool.
func NewStockLevelStorage(pool *pgxpool.Pool) *StockLevelStorage {
	return &StockLevelStorage{pool: pool}
}

// GetBySKU fetches a single stock_levels row by sku. Returns pgx.ErrNoRows on miss.
func (s *StockLevelStorage) GetBySKU(ctx context.Context, sku string) (*inventory.StockLevel, error) {
	const q = `
		SELECT sku, available_qty, reserved_qty, sold_qty, version, updated_at
		FROM inventory.stock_levels
		WHERE sku = $1`

	row := s.pool.QueryRow(ctx, q, sku)
	var sl inventory.StockLevel
	err := row.Scan(&sl.SKU, &sl.AvailableQty, &sl.ReservedQty, &sl.SoldQty, &sl.Version, &sl.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("GetBySKU scan: %w", err)
	}
	return &sl, nil
}

// GetBySKUForUpdate fetches and row-locks a single stock_levels row.
// Must be called inside an open transaction (tx).
//
// LOCK-ORDER PIN: callers MUST acquire stock_levels locks (lex order by SKU)
// BEFORE acquiring any reservations lock. See td.json concurrency_model.lock_ordering.
func (s *StockLevelStorage) GetBySKUForUpdate(ctx context.Context, tx pgx.Tx, sku string) (*inventory.StockLevel, error) {
	const q = `
		SELECT sku, available_qty, reserved_qty, sold_qty, version, updated_at
		FROM inventory.stock_levels
		WHERE sku = $1
		FOR UPDATE`

	row := tx.QueryRow(ctx, q, sku)
	var sl inventory.StockLevel
	err := row.Scan(&sl.SKU, &sl.AvailableQty, &sl.ReservedQty, &sl.SoldQty, &sl.Version, &sl.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("GetBySKUForUpdate scan: %w", err)
	}
	return &sl, nil
}

// GetManyBySKUsForUpdate fetches and row-locks multiple stock_levels rows in a single
// SELECT WHERE sku = ANY($1) FOR UPDATE. Caller must sort skus lexicographically
// before calling to satisfy the lock-order pin.
//
// LOCK-ORDER PIN: SKUs MUST be sorted in lexicographic order by the caller before
// this call. See td.json concurrency_model.lock_ordering.
func (s *StockLevelStorage) GetManyBySKUsForUpdate(ctx context.Context, tx pgx.Tx, skus []string) (map[string]*inventory.StockLevel, error) {
	const q = `
		SELECT sku, available_qty, reserved_qty, sold_qty, version, updated_at
		FROM inventory.stock_levels
		WHERE sku = ANY($1)
		ORDER BY sku ASC
		FOR UPDATE`

	rows, err := tx.Query(ctx, q, skus)
	if err != nil {
		return nil, fmt.Errorf("GetManyBySKUsForUpdate query: %w", err)
	}
	defer rows.Close()

	result := make(map[string]*inventory.StockLevel, len(skus))
	for rows.Next() {
		var sl inventory.StockLevel
		if err := rows.Scan(&sl.SKU, &sl.AvailableQty, &sl.ReservedQty, &sl.SoldQty, &sl.Version, &sl.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetManyBySKUsForUpdate scan: %w", err)
		}
		result[sl.SKU] = &sl
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetManyBySKUsForUpdate rows: %w", err)
	}
	return result, nil
}

// DecrementAvailableIncrementReserved updates available_qty and reserved_qty atomically.
// Must be called inside an open transaction after the row is already locked via
// GetManyBySKUsForUpdate or GetBySKUForUpdate.
// No-negative invariant: caller must have verified available_qty >= qty before calling.
func (s *StockLevelStorage) DecrementAvailableIncrementReserved(ctx context.Context, tx pgx.Tx, sku string, qty int) error {
	const q = `
		UPDATE inventory.stock_levels
		SET available_qty = available_qty - $2,
		    reserved_qty  = reserved_qty  + $2,
		    version       = version + 1,
		    updated_at    = NOW()
		WHERE sku = $1`

	cmd, err := tx.Exec(ctx, q, sku, qty)
	if err != nil {
		return fmt.Errorf("DecrementAvailableIncrementReserved exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("DecrementAvailableIncrementReserved: no row for sku=%s", sku)
	}
	return nil
}

// IncrementAvailableDecrementReserved reverses a reservation (release path).
// Must be called inside an open transaction after the stock row is locked.
//
// LOCK-ORDER PIN: stock_levels BEFORE reservations.
func (s *StockLevelStorage) IncrementAvailableDecrementReserved(ctx context.Context, tx pgx.Tx, sku string, qty int) error {
	const q = `
		UPDATE inventory.stock_levels
		SET available_qty = available_qty + $2,
		    reserved_qty  = reserved_qty  - $2,
		    version       = version + 1,
		    updated_at    = NOW()
		WHERE sku = $1`

	cmd, err := tx.Exec(ctx, q, sku, qty)
	if err != nil {
		return fmt.Errorf("IncrementAvailableDecrementReserved exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("IncrementAvailableDecrementReserved: no row for sku=%s", sku)
	}
	return nil
}

// DecrementReservedIncrementSold transitions reserved→sold (payment.completed path).
// Must be called inside an open transaction after the stock row is locked.
//
// LOCK-ORDER PIN: stock_levels BEFORE reservations.
func (s *StockLevelStorage) DecrementReservedIncrementSold(ctx context.Context, tx pgx.Tx, sku string, qty int) error {
	const q = `
		UPDATE inventory.stock_levels
		SET reserved_qty = reserved_qty - $2,
		    sold_qty     = sold_qty     + $2,
		    version      = version + 1,
		    updated_at   = NOW()
		WHERE sku = $1`

	cmd, err := tx.Exec(ctx, q, sku, qty)
	if err != nil {
		return fmt.Errorf("DecrementReservedIncrementSold exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("DecrementReservedIncrementSold: no row for sku=%s", sku)
	}
	return nil
}

// DecrementSoldIncrementAvailable transitions sold→available (order.cancelled on COMMITTED path / PR-006).
// Must be called inside an open transaction after the stock row is locked.
//
// LOCK-ORDER PIN: stock_levels BEFORE reservations.
func (s *StockLevelStorage) DecrementSoldIncrementAvailable(ctx context.Context, tx pgx.Tx, sku string, qty int) error {
	const q = `
		UPDATE inventory.stock_levels
		SET sold_qty      = sold_qty      - $2,
		    available_qty = available_qty + $2,
		    version       = version + 1,
		    updated_at    = NOW()
		WHERE sku = $1`

	cmd, err := tx.Exec(ctx, q, sku, qty)
	if err != nil {
		return fmt.Errorf("DecrementSoldIncrementAvailable exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("DecrementSoldIncrementAvailable: no row for sku=%s", sku)
	}
	return nil
}

// AdjustQuantities applies a signed delta to available_qty (admin stock.adjust path).
// The caller must have already checked available_qty + delta >= 0 under FOR UPDATE.
func (s *StockLevelStorage) AdjustQuantities(ctx context.Context, tx pgx.Tx, sku string, delta int) error {
	const q = `
		UPDATE inventory.stock_levels
		SET available_qty = available_qty + $2,
		    version       = version + 1,
		    updated_at    = NOW()
		WHERE sku = $1`

	cmd, err := tx.Exec(ctx, q, sku, delta)
	if err != nil {
		return fmt.Errorf("AdjustQuantities exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("AdjustQuantities: no row for sku=%s", sku)
	}
	return nil
}

// Insert inserts a new stock_levels row (product.created consumer path).
// Uses ON CONFLICT (sku) DO NOTHING for natural idempotency.
func (s *StockLevelStorage) Insert(ctx context.Context, tx pgx.Tx, sku string) error {
	const q = `
		INSERT INTO inventory.stock_levels (sku, available_qty, reserved_qty, sold_qty, version)
		VALUES ($1, 0, 0, 0, 0)
		ON CONFLICT (sku) DO NOTHING`

	_, err := tx.Exec(ctx, q, sku)
	if err != nil {
		return fmt.Errorf("StockLevelStorage.Insert exec: %w", err)
	}
	return nil
}
