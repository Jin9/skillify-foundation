package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity/app/identity"
)

const pgUniqueViolation = "23505"

var ErrUserNotFound = errors.New("user not found")

// userStorage implements identity.UserStorage.
type userStorage struct {
	db *pgxpool.Pool
}

var _ identity.UserStorage = (*userStorage)(nil)

func NewUserStorage(db *pgxpool.Pool) identity.UserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) Insert(ctx context.Context, u identity.User) error {
	q := `
		INSERT INTO users (id, email, password_hash, name, phone, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.Exec(ctx, q,
		u.ID, u.Email, u.PasswordHash, u.Name, u.Phone, u.Role, u.Status,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return fmt.Errorf("%w: %s", identity.ErrEmailAlreadyRegistered, pgErr.Detail)
		}
		return err
	}
	return nil
}

func (s *userStorage) GetByEmail(ctx context.Context, email string) (identity.User, error) {
	q := `
		SELECT id, email, password_hash, name, phone, role, status, created_at, updated_at
		FROM users WHERE email = $1
	`
	row := s.db.QueryRow(ctx, q, email)
	return scanUser(row)
}

func (s *userStorage) GetByID(ctx context.Context, id uuid.UUID) (identity.User, error) {
	q := `
		SELECT id, email, password_hash, name, phone, role, status, created_at, updated_at
		FROM users WHERE id = $1
	`
	row := s.db.QueryRow(ctx, q, id)
	return scanUser(row)
}

func scanUser(row pgx.Row) (identity.User, error) {
	var u identity.User
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Phone,
		&u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return identity.User{}, ErrUserNotFound
	}
	return u, err
}
