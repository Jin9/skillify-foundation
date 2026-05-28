package access

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryStorage is the pgx-backed repository for categories.
// Methods are stubbed — full implementation deferred to post-MVP.
type CategoryStorage struct {
	db *pgxpool.Pool
}

// NewCategoryStorage constructs a CategoryStorage.
func NewCategoryStorage(db *pgxpool.Pool) *CategoryStorage {
	return &CategoryStorage{db: db}
}

// TODO: implement category.create storage
// TODO: implement category.update storage
// TODO: implement category.deactivate storage
// TODO: implement category.list storage
