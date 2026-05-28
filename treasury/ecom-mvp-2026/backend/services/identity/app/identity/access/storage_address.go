package access

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity/app/identity"
)

// addressStorage implements identity.AddressStorage.
// Only SetDefault is fully implemented for MVP; remaining operations (create,
// list, update, delete) are deferred — see README NOT_IMPLEMENTED_MVP list.
type addressStorage struct {
	db *pgxpool.Pool
}

var _ identity.AddressStorage = (*addressStorage)(nil)

func NewAddressStorage(db *pgxpool.Pool) identity.AddressStorage {
	return &addressStorage{db: db}
}

// SetDefault runs the atomic two-step UPDATE described in the address set-default
// spec (identity-address/set-default.md):
//
//  1. Clear is_default on all other live addresses for the user (zero rows
//     affected here is fine — there may not be a prior default).
//  2. Set is_default = true on the target; capture RowsAffected.
//
// Returns RowsAffected of step 2. The caller (service) maps 0 → ErrAddressNotFound.
func (s *addressStorage) SetDefault(ctx context.Context, userID, addressID uuid.UUID) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Step 1: clear is_default on every other live address for this user.
	clearQ := `
		UPDATE addresses
		SET    is_default = false, updated_at = NOW()
		WHERE  user_id = $1
		  AND  deleted_at IS NULL
		  AND  is_default = true
		  AND  id <> $2
	`
	if _, err = tx.Exec(ctx, clearQ, userID, addressID); err != nil {
		return 0, err
	}

	// Step 2: set is_default = true on the target.
	setQ := `
		UPDATE addresses
		SET    is_default = true, updated_at = NOW()
		WHERE  id = $1
		  AND  user_id = $2
		  AND  deleted_at IS NULL
	`
	tag, err := tx.Exec(ctx, setQ, addressID, userID)
	if err != nil {
		return 0, err
	}

	rows := tag.RowsAffected()
	if rows == 0 {
		// Target not found, soft-deleted, or wrong user — rollback and signal to caller.
		_ = tx.Rollback(ctx)
		return 0, nil
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	slog.InfoContext(ctx, "address.set_default",
		slog.String("user_id", userID.String()),
		slog.String("address_id", addressID.String()),
	)

	return rows, nil
}
