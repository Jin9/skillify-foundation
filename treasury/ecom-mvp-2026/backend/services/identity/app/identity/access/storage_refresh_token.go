package access

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity/app/identity"
)

// refreshTokenStorage implements identity.RefreshTokenStorage.
type refreshTokenStorage struct {
	db *pgxpool.Pool
}

var _ identity.RefreshTokenStorage = (*refreshTokenStorage)(nil)

func NewRefreshTokenStorage(db *pgxpool.Pool) identity.RefreshTokenStorage {
	return &refreshTokenStorage{db: db}
}

func (s *refreshTokenStorage) Insert(ctx context.Context, rt identity.RefreshToken) error {
	q := `
		INSERT INTO refresh_tokens (jti, user_id, issued_at, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.db.Exec(ctx, q, rt.JTI, rt.UserID, rt.IssuedAt, rt.ExpiresAt)
	return err
}

// RotateInTx executes the single-tx rotation (ID-DEC-1):
//  1. UPDATE refresh_tokens SET revoked_at=NOW(), replaced_by_jti=$new WHERE jti=$old AND revoked_at IS NULL
//  2. Verify RowsAffected == 1
//  3. INSERT new refresh_tokens row
//  4. COMMIT
func (s *refreshTokenStorage) RotateInTx(ctx context.Context, oldJTI, newJTI uuid.UUID, newRT identity.RefreshToken) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	updateQ := `
		UPDATE refresh_tokens
		SET revoked_at = NOW(), replaced_by_jti = $1
		WHERE jti = $2 AND revoked_at IS NULL
	`
	tag, err := tx.Exec(ctx, updateQ, newJTI, oldJTI)
	if err != nil {
		return 0, err
	}

	rows := tag.RowsAffected()
	if rows == 0 {
		_ = tx.Rollback(ctx)
		return 0, nil
	}

	insertQ := `
		INSERT INTO refresh_tokens (jti, user_id, issued_at, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, insertQ, newRT.JTI, newRT.UserID, newRT.IssuedAt, newRT.ExpiresAt)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	slog.InfoContext(ctx, "refresh_token.rotated",
		slog.String("old_jti", oldJTI.String()),
		slog.String("new_jti", newJTI.String()),
	)

	return rows, nil
}

func (s *refreshTokenStorage) Revoke(ctx context.Context, jti uuid.UUID) error {
	q := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE jti = $1 AND revoked_at IS NULL`
	_, err := s.db.Exec(ctx, q, jti)
	return err
}
