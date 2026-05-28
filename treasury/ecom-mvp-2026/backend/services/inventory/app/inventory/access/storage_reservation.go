package access

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory"
)

// ReservationStorage provides data access for the reservations table.
type ReservationStorage struct {
	pool *pgxpool.Pool
}

var _ inventory.ReservationStorage = (*ReservationStorage)(nil)

// NewReservationStorage constructs a ReservationStorage backed by the given pool.
func NewReservationStorage(pool *pgxpool.Pool) *ReservationStorage {
	return &ReservationStorage{pool: pool}
}

// Insert creates a new RESERVED reservation row inside the given transaction.
// Must be called AFTER the corresponding stock_levels row is already locked
// (lock-order pin: stock_levels first, then reservations).
func (s *ReservationStorage) Insert(ctx context.Context, tx pgx.Tx, r *inventory.Reservation) error {
	const q = `
		INSERT INTO inventory.reservations
		    (id, order_id, sku, qty, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.Exec(ctx, q,
		r.ID, r.OrderID, r.SKU, r.Qty, string(r.Status), r.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("ReservationStorage.Insert exec: %w", err)
	}
	return nil
}

// FindByOrderID fetches all reservations for an order as a snapshot (no lock).
// Use this to obtain the SKU set before acquiring stock_levels locks, so that
// stock_levels can be locked first (lock-order pin: stock_levels BEFORE reservations).
// After calling this, lock stock_levels in lex SKU order, then re-lock each
// reservation by PK using GetByIDForUpdate.
func (s *ReservationStorage) FindByOrderID(ctx context.Context, tx pgx.Tx, orderID uuid.UUID) ([]*inventory.Reservation, error) {
	const q = `
		SELECT id, order_id, sku, qty, status, expires_at, created_at,
		       released_at, release_reason
		FROM inventory.reservations
		WHERE order_id = $1
		ORDER BY sku ASC`

	rows, err := tx.Query(ctx, q, orderID)
	if err != nil {
		return nil, fmt.Errorf("FindByOrderID query: %w", err)
	}
	defer rows.Close()
	return scanReservations(rows)
}

// FindByOrderIDForUpdate fetches all reservations for an order and row-locks them
// in ascending sku order to honour the canonical lock-order pin.
//
// LOCK-ORDER PIN: stock_levels(sku) must ALREADY be locked before calling this.
// See td.json concurrency_model.lock_ordering.
func (s *ReservationStorage) FindByOrderIDForUpdate(ctx context.Context, tx pgx.Tx, orderID uuid.UUID) ([]*inventory.Reservation, error) {
	const q = `
		SELECT id, order_id, sku, qty, status, expires_at, created_at,
		       released_at, release_reason
		FROM inventory.reservations
		WHERE order_id = $1
		ORDER BY sku ASC
		FOR UPDATE`

	rows, err := tx.Query(ctx, q, orderID)
	if err != nil {
		return nil, fmt.Errorf("FindByOrderIDForUpdate query: %w", err)
	}
	defer rows.Close()
	return scanReservations(rows)
}

// FindExpiredCandidates scans for RESERVED rows whose expires_at is in the past.
// Uses FOR UPDATE SKIP LOCKED so two sweeper replicas never contend on the same row.
// The candidate scan itself runs in its own short tx; each row is then re-locked
// individually in the per-row tick tx.
func (s *ReservationStorage) FindExpiredCandidates(ctx context.Context, tx pgx.Tx, batchSize int) ([]*inventory.Reservation, error) {
	const q = `
		SELECT id, order_id, sku, qty, status, expires_at, created_at,
		       released_at, release_reason
		FROM inventory.reservations
		WHERE status = 'RESERVED'
		  AND expires_at < NOW()
		ORDER BY expires_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := tx.Query(ctx, q, batchSize)
	if err != nil {
		return nil, fmt.Errorf("FindExpiredCandidates query: %w", err)
	}
	defer rows.Close()
	return scanReservations(rows)
}

// GetByIDForUpdate fetches and row-locks a single reservation by PK.
// Must be called AFTER the corresponding stock_levels row is already locked.
//
// LOCK-ORDER PIN: stock_levels BEFORE reservations.
func (s *ReservationStorage) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*inventory.Reservation, error) {
	const q = `
		SELECT id, order_id, sku, qty, status, expires_at, created_at,
		       released_at, release_reason
		FROM inventory.reservations
		WHERE id = $1
		FOR UPDATE`

	row := tx.QueryRow(ctx, q, id)
	r, err := scanReservation(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("GetByIDForUpdate scan: %w", err)
	}
	return r, nil
}

// MarkExpired transitions a reservation from RESERVED to EXPIRED.
// Must be called inside a tx where the row is already locked.
func (s *ReservationStorage) MarkExpired(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	reason := string(inventory.ReasonExpired)
	const q = `
		UPDATE inventory.reservations
		SET status         = 'EXPIRED',
		    released_at    = NOW(),
		    release_reason = $2
		WHERE id = $1`

	cmd, err := tx.Exec(ctx, q, id, reason)
	if err != nil {
		return fmt.Errorf("MarkExpired exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("MarkExpired: no row for id=%s", id)
	}
	return nil
}

// MarkCommitted transitions a reservation from RESERVED to COMMITTED.
// Must be called inside a tx where the row is already locked.
func (s *ReservationStorage) MarkCommitted(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	const q = `
		UPDATE inventory.reservations
		SET status = 'COMMITTED'
		WHERE id = $1`

	cmd, err := tx.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("MarkCommitted exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("MarkCommitted: no row for id=%s", id)
	}
	return nil
}

// MarkReleased transitions a reservation to RELEASED with the supplied reason.
// Must be called inside a tx where the row is already locked.
func (s *ReservationStorage) MarkReleased(ctx context.Context, tx pgx.Tx, id uuid.UUID, reason inventory.ReleaseReason) error {
	now := time.Now().UTC()
	const q = `
		UPDATE inventory.reservations
		SET status         = 'RELEASED',
		    released_at    = $2,
		    release_reason = $3
		WHERE id = $1`

	cmd, err := tx.Exec(ctx, q, id, now, string(reason))
	if err != nil {
		return fmt.Errorf("MarkReleased exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("MarkReleased: no row for id=%s", id)
	}
	return nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func scanReservations(rows pgx.Rows) ([]*inventory.Reservation, error) {
	var result []*inventory.Reservation
	for rows.Next() {
		r, err := scanReservationFromRows(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

// scanReservation scans a pgx.Row into a Reservation.
func scanReservation(row pgx.Row) (*inventory.Reservation, error) {
	var r inventory.Reservation
	var statusStr string
	var reasonStr *string

	err := row.Scan(
		&r.ID, &r.OrderID, &r.SKU, &r.Qty,
		&statusStr, &r.ExpiresAt, &r.CreatedAt,
		&r.ReleasedAt, &reasonStr,
	)
	if err != nil {
		return nil, err
	}
	r.Status = inventory.ReservationStatus(statusStr)
	if reasonStr != nil {
		rr := inventory.ReleaseReason(*reasonStr)
		r.ReleaseReason = &rr
	}
	return &r, nil
}

func scanReservationFromRows(rows pgx.Rows) (*inventory.Reservation, error) {
	var r inventory.Reservation
	var statusStr string
	var reasonStr *string

	err := rows.Scan(
		&r.ID, &r.OrderID, &r.SKU, &r.Qty,
		&statusStr, &r.ExpiresAt, &r.CreatedAt,
		&r.ReleasedAt, &reasonStr,
	)
	if err != nil {
		return nil, fmt.Errorf("scanReservationFromRows: %w", err)
	}
	r.Status = inventory.ReservationStatus(statusStr)
	if reasonStr != nil {
		rr := inventory.ReleaseReason(*reasonStr)
		r.ReleaseReason = &rr
	}
	return &r, nil
}
