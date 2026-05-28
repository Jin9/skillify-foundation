package access

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("idempotency key not found")

// IdempotencyStatus mirrors checkout.IdempotencyStatus without a circular import.
type IdempotencyStatus string

const (
	StatusInflight  IdempotencyStatus = "INFLIGHT"
	StatusCompleted IdempotencyStatus = "COMPLETED"
	StatusAbandoned IdempotencyStatus = "ABANDONED"
)

// IdempotencyRow is the raw DB row.
type IdempotencyRow struct {
	Key              string
	CustomerUserID   string
	RequestHash      string
	Status           IdempotencyStatus
	ResponseEnvelope json.RawMessage
	HTTPStatus       *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExpiresAt        time.Time
}

// IdempotencyStorage handles all DB operations on the idempotency_keys table
// within an existing pgx transaction. All methods accept a pgx.Tx so the caller
// keeps the transaction boundary.
type IdempotencyStorage struct{}

func NewIdempotencyStorage() *IdempotencyStorage {
	return &IdempotencyStorage{}
}

// TryClaimOrLookup atomically attempts to INSERT an INFLIGHT row.
//
// Race-safety: the INSERT ... ON CONFLICT DO NOTHING is a single atomic
// statement — Postgres acquires the unique-index lock before any other
// connection's concurrent INSERT can proceed, making it impossible for two
// callers to both observe "row not found" and both proceed to insert.
//
// Returns:
//   - winner=true, existing=nil  → this caller claimed the key; proceed with orchestration.
//   - winner=false, existing!=nil → conflict; caller must inspect existing.Status:
//       COMPLETED + matching hash   → return cached envelope.
//       COMPLETED + different hash  → 409 IDEMPOTENCY_KEY_REUSED.
//       INFLIGHT                    → 409 IDEMPOTENCY_KEY_INFLIGHT.
//       ABANDONED                   → treat as fresh (janitor cleared it); re-claim via DELETE+INSERT.
func (s *IdempotencyStorage) TryClaimOrLookup(ctx context.Context, tx pgx.Tx, key, customerUserID, requestHash string) (winner bool, existing *IdempotencyRow, err error) {
	// Attempt atomic INSERT; RETURNING key tells us whether the row was created.
	const insertQ = `
		INSERT INTO checkout.idempotency_keys
		    (key, customer_user_id, request_hash, status, response_envelope, http_status, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, 'INFLIGHT', NULL, NULL, NOW(), NOW(), NOW() + INTERVAL '24 hours')
		ON CONFLICT (key, customer_user_id) DO NOTHING
		RETURNING key`

	var returnedKey string
	err = tx.QueryRow(ctx, insertQ, key, customerUserID, requestHash).Scan(&returnedKey)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, nil, fmt.Errorf("idempotency try-claim: %w", err)
	}

	if err == nil {
		// INSERT returned a row — we are the winner.
		return true, nil, nil
	}

	// INSERT returned 0 rows (ON CONFLICT DO NOTHING) — a concurrent or prior
	// row already exists. SELECT FOR UPDATE to read it (and block on any open
	// transaction holding the row-lock, ensuring we see the final status once
	// the leader commits or rolls back).
	const selectQ = `
		SELECT key, customer_user_id, request_hash, status, response_envelope, http_status,
		       created_at, updated_at, expires_at
		FROM checkout.idempotency_keys
		WHERE key = $1 AND customer_user_id = $2
		FOR UPDATE`

	row := tx.QueryRow(ctx, selectQ, key, customerUserID)
	var r IdempotencyRow
	var envelope []byte
	scanErr := row.Scan(
		&r.Key, &r.CustomerUserID, &r.RequestHash, &r.Status,
		&envelope, &r.HTTPStatus,
		&r.CreatedAt, &r.UpdatedAt, &r.ExpiresAt,
	)
	if errors.Is(scanErr, pgx.ErrNoRows) {
		// Row vanished between the failed INSERT and the SELECT (extremely rare:
		// the janitor deleted an ABANDONED row in that tiny window). Treat as
		// ErrNotFound so the caller can retry the claim.
		return false, nil, ErrNotFound
	}
	if scanErr != nil {
		return false, nil, fmt.Errorf("idempotency post-conflict lookup: %w", scanErr)
	}
	r.ResponseEnvelope = json.RawMessage(envelope)
	return false, &r, nil
}

// Lookup performs SELECT ... FOR UPDATE inside the caller's transaction.
// Returns ErrNotFound if no row exists.
//
// Deprecated: prefer TryClaimOrLookup for first-time claim paths to avoid the
// INFLIGHT race window (H02). Lookup is retained for read-only paths (e.g.
// janitor, admin tooling) that only ever read an existing row.
func (s *IdempotencyStorage) Lookup(ctx context.Context, tx pgx.Tx, key, customerUserID string) (IdempotencyRow, error) {
	const q = `
		SELECT key, customer_user_id, request_hash, status, response_envelope, http_status,
		       created_at, updated_at, expires_at
		FROM checkout.idempotency_keys
		WHERE key = $1 AND customer_user_id = $2
		FOR UPDATE`

	row := tx.QueryRow(ctx, q, key, customerUserID)

	var r IdempotencyRow
	var envelope []byte
	err := row.Scan(
		&r.Key, &r.CustomerUserID, &r.RequestHash, &r.Status,
		&envelope, &r.HTTPStatus,
		&r.CreatedAt, &r.UpdatedAt, &r.ExpiresAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return IdempotencyRow{}, ErrNotFound
	}
	if err != nil {
		return IdempotencyRow{}, fmt.Errorf("idempotency lookup: %w", err)
	}
	r.ResponseEnvelope = json.RawMessage(envelope)
	return r, nil
}

// InsertInflight writes a new INFLIGHT row. Must run inside the caller's transaction.
//
// Deprecated: prefer TryClaimOrLookup which closes the Lookup→InsertInflight race
// window (H02). InsertInflight is retained for backward-compatibility only.
func (s *IdempotencyStorage) InsertInflight(ctx context.Context, tx pgx.Tx, key, customerUserID, requestHash string, expiresAt time.Time) error {
	const q = `
		INSERT INTO checkout.idempotency_keys
		    (key, customer_user_id, request_hash, status, response_envelope, http_status, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, 'INFLIGHT', NULL, NULL, NOW(), NOW(), $4)`

	_, err := tx.Exec(ctx, q, key, customerUserID, requestHash, expiresAt)
	if err != nil {
		return fmt.Errorf("idempotency insert inflight: %w", err)
	}
	return nil
}

// Finalize sets status=COMPLETED (or ABANDONED) with the result envelope.
// Must run inside the caller's transaction.
func (s *IdempotencyStorage) Finalize(ctx context.Context, tx pgx.Tx, key, customerUserID string, status IdempotencyStatus, envelope json.RawMessage, httpStatus int) error {
	const q = `
		UPDATE checkout.idempotency_keys
		SET status = $3, response_envelope = $4, http_status = $5, updated_at = NOW()
		WHERE key = $1 AND customer_user_id = $2`

	_, err := tx.Exec(ctx, q, key, customerUserID, string(status), []byte(envelope), httpStatus)
	if err != nil {
		return fmt.Errorf("idempotency finalize: %w", err)
	}
	return nil
}

// CleanupExpired deletes rows past their TTL in batches.
// Runs outside any transaction — called from the background janitor goroutine.
func (s *IdempotencyStorage) CleanupExpired(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	const q = `
		DELETE FROM checkout.idempotency_keys
		WHERE expires_at < NOW()
		  AND status IN ('COMPLETED','ABANDONED')
		LIMIT 500`

	tag, err := pool.Exec(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("idempotency cleanup: %w", err)
	}
	return tag.RowsAffected(), nil
}

// MarkStaleInflightAbandoned sweeps INFLIGHT rows older than 60s and marks them ABANDONED.
// Per spec janitor: every 60s, LIMIT 100, SKIP LOCKED.
func (s *IdempotencyStorage) MarkStaleInflightAbandoned(ctx context.Context, pool *pgxpool.Pool, abandonedEnvelope json.RawMessage, httpStatus int) (int64, error) {
	const q = `
		UPDATE checkout.idempotency_keys
		SET status = 'ABANDONED', response_envelope = $1, http_status = $2, updated_at = NOW()
		WHERE key IN (
		    SELECT key FROM checkout.idempotency_keys
		    WHERE status = 'INFLIGHT' AND created_at < NOW() - INTERVAL '60 seconds'
		    LIMIT 100
		    FOR UPDATE SKIP LOCKED
		)`

	tag, err := pool.Exec(ctx, q, []byte(abandonedEnvelope), httpStatus)
	if err != nil {
		return 0, fmt.Errorf("idempotency mark abandoned: %w", err)
	}
	return tag.RowsAffected(), nil
}
