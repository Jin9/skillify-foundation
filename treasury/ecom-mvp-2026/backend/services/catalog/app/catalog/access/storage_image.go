package access

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ImageStorage is the pgx-backed repository for product_images.
type ImageStorage struct {
	db *pgxpool.Pool
}

// NewImageStorage constructs an ImageStorage.
func NewImageStorage(db *pgxpool.Pool) *ImageStorage {
	return &ImageStorage{db: db}
}

// InsertImageParams carries the values for a single product_images row.
type InsertImageParams struct {
	ProductID uuid.UUID
	URL       string
	SortOrder int
}

// Insert inserts a single image row inside the caller-managed transaction q.
func (s *ImageStorage) Insert(ctx context.Context, q Querier, p InsertImageParams) error {
	const query = `
		INSERT INTO catalog.product_images (id, product_id, url, sort_order, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW())`

	if _, err := q.Exec(ctx, query, p.ProductID, p.URL, p.SortOrder); err != nil {
		return fmt.Errorf("storage: insert image: %w", err)
	}
	return nil
}

// TODO: implement DeleteByProductID for product.update (replace image set atomically)
