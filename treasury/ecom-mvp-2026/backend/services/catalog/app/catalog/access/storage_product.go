package access

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"
const fkViolationCode = "23503"

// Sentinel errors from the access layer — mapped to domain errors by the service.
var (
	ErrProductNotFound  = errors.New("access: product not found")
	ErrDuplicateSKU     = errors.New("access: duplicate sku")
	ErrCategoryNotFound = errors.New("access: category not found or inactive")
)

// ProductStorage is the pgx-backed repository for products.
type ProductStorage struct {
	db *pgxpool.Pool
}

// NewProductStorage constructs a ProductStorage.
func NewProductStorage(db *pgxpool.Pool) *ProductStorage {
	return &ProductStorage{db: db}
}

// InsertProductParams carries the values for a new product row.
type InsertProductParams struct {
	ID          uuid.UUID
	SKU         string
	Name        string
	Description string
	Price       string // NUMERIC as string
	CategoryID  uuid.UUID
	Status      string
}

// InsertStatusHistoryParams carries values for a product_status_history row.
type InsertStatusHistoryParams struct {
	ProductID   uuid.UUID
	FromStatus  *string // nil = initial creation
	ToStatus    string
	ActorUserID uuid.UUID
}

// ProductListFilter holds optional filter parameters for product.list.
// Mirrors the domain filter — access layer owns this to avoid circular import.
type ProductListFilter struct {
	CategoryID *uuid.UUID
	MinPrice   *string // NUMERIC string (e.g. "10.00")
	MaxPrice   *string
	Search     *string
	// InStock is not applied at DB level — handled by service after inventory call
}

// SortOption is a validated sort enum for product.list.
type SortOption string

const (
	SortPriceAsc    SortOption = "price_asc"
	SortPriceDesc   SortOption = "price_desc"
	SortNameAsc     SortOption = "name_asc"
	SortNameDesc    SortOption = "name_desc"
	SortCreatedDesc SortOption = "created_desc"
)

// ValidSortOptions is the allowed set for handler-level validation.
var ValidSortOptions = map[SortOption]struct{}{
	SortPriceAsc:    {},
	SortPriceDesc:   {},
	SortNameAsc:     {},
	SortNameDesc:    {},
	SortCreatedDesc: {},
}

// ProductListRow is the DB-level projection for product.list items.
type ProductListRow struct {
	ID           uuid.UUID
	SKU          string
	Name         string
	Price        string
	CategoryID   uuid.UUID
	CategoryName string
	CategorySlug string
	ThumbnailURL string // empty string if no image
	Status       string
}

// ProductDetailRow is the DB-level projection for product.detail.
type ProductDetailRow struct {
	ID           uuid.UUID
	SKU          string
	Name         string
	Description  string
	Price        string
	Images       []string
	CategoryID   uuid.UUID
	CategoryName string
	CategorySlug string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Querier is satisfied by both *pgxpool.Pool and pgx.Tx, allowing the same
// storage functions to participate in a caller-managed transaction.
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Insert performs a single-row INSERT for a new product.
// Callers pass a pgx.Tx so this participates in a broader transaction.
// Returns ErrDuplicateSKU on unique-constraint violation on sku.
// Returns ErrCategoryNotFound on foreign-key violation on category_id.
func (s *ProductStorage) Insert(ctx context.Context, q Querier, p InsertProductParams) error {
	const query = `
		INSERT INTO catalog.products
			(id, sku, name, description, price, category_id, status, visible, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5::numeric, $6, $7, TRUE, NOW(), NOW())`

	_, err := q.Exec(ctx, query,
		p.ID, p.SKU, p.Name, p.Description, p.Price, p.CategoryID, p.Status,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case uniqueViolationCode:
				return ErrDuplicateSKU
			case fkViolationCode:
				return ErrCategoryNotFound
			}
		}
		return fmt.Errorf("storage: insert product: %w", err)
	}
	return nil
}

// InsertStatusHistory writes one audit row to product_status_history.
func (s *ProductStorage) InsertStatusHistory(ctx context.Context, q Querier, p InsertStatusHistoryParams) error {
	const query = `
		INSERT INTO catalog.product_status_history
			(id, product_id, from_status, to_status, actor_user_id, created_at)
		VALUES
			(gen_random_uuid(), $1, $2, $3, $4, NOW())`

	_, err := q.Exec(ctx, query, p.ProductID, p.FromStatus, p.ToStatus, p.ActorUserID)
	if err != nil {
		return fmt.Errorf("storage: insert status history: %w", err)
	}
	return nil
}

// List executes the parameterized product list query with optional filters and pagination.
// Only ACTIVE + visible=TRUE products are returned.
func (s *ProductStorage) List(
	ctx context.Context,
	filter ProductListFilter,
	sort SortOption,
	page, limit int,
) ([]ProductListRow, int, error) {
	offset := (page - 1) * limit

	// Build ORDER BY from validated enum — not user-controlled string concat
	orderClause := buildOrderClause(sort)

	const baseWhere = `
		WHERE p.status = 'ACTIVE'
		  AND p.visible = TRUE
		  AND ($1::uuid IS NULL OR p.category_id = $1::uuid)
		  AND ($2::numeric IS NULL OR p.price >= $2::numeric)
		  AND ($3::numeric IS NULL OR p.price <= $3::numeric)
		  AND ($4::text IS NULL OR p.name ILIKE '%' || $4 || '%')`

	countQuery := `
		SELECT COUNT(*)
		FROM catalog.products p
		` + baseWhere

	listQuery := `
		SELECT
			p.id, p.sku, p.name, p.price::text,
			p.category_id, c.name AS category_name, c.slug AS category_slug,
			COALESCE(
				(SELECT pi.url FROM catalog.product_images pi
				 WHERE pi.product_id = p.id ORDER BY pi.sort_order ASC LIMIT 1),
				''
			) AS thumbnail_url,
			p.status
		FROM catalog.products p
		JOIN catalog.categories c ON c.id = p.category_id
		` + baseWhere + `
		` + orderClause + `
		LIMIT $6 OFFSET $7`

	var catID interface{} = nil
	var minPrice interface{} = nil
	var maxPrice interface{} = nil
	var search interface{} = nil

	if filter.CategoryID != nil {
		catID = filter.CategoryID
	}
	if filter.MinPrice != nil {
		minPrice = *filter.MinPrice
	}
	if filter.MaxPrice != nil {
		maxPrice = *filter.MaxPrice
	}
	if filter.Search != nil && strings.TrimSpace(*filter.Search) != "" {
		trimmed := strings.TrimSpace(*filter.Search)
		search = trimmed
	}

	whereArgs := []any{catID, minPrice, maxPrice, search}

	// Execute count query
	var total int
	if err := s.db.QueryRow(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("storage: count products: %w", err)
	}

	// Execute list query
	listArgs := append(whereArgs, limit, offset)
	rows, err := s.db.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("storage: list products: %w", err)
	}
	defer rows.Close()

	var items []ProductListRow
	for rows.Next() {
		var r ProductListRow
		if err := rows.Scan(
			&r.ID, &r.SKU, &r.Name, &r.Price,
			&r.CategoryID, &r.CategoryName, &r.CategorySlug,
			&r.ThumbnailURL,
			&r.Status,
		); err != nil {
			return nil, 0, fmt.Errorf("storage: scan product row: %w", err)
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("storage: iterate product rows: %w", err)
	}

	return items, total, nil
}

// GetByID fetches a single product with its category and images.
// Does NOT apply visibility rules — that is the service layer's responsibility.
// Returns ErrProductNotFound if the product does not exist.
func (s *ProductStorage) GetByID(ctx context.Context, productID uuid.UUID) (ProductDetailRow, error) {
	const productQuery = `
		SELECT
			p.id, p.sku, p.name, p.description, p.price::text,
			p.category_id, c.name AS cat_name, c.slug AS cat_slug,
			p.status, p.created_at, p.updated_at
		FROM catalog.products p
		JOIN catalog.categories c ON c.id = p.category_id
		WHERE p.id = $1`

	var row ProductDetailRow
	err := s.db.QueryRow(ctx, productQuery, productID).Scan(
		&row.ID, &row.SKU, &row.Name, &row.Description, &row.Price,
		&row.CategoryID, &row.CategoryName, &row.CategorySlug,
		&row.Status, &row.CreatedAt, &row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProductDetailRow{}, ErrProductNotFound
		}
		return ProductDetailRow{}, fmt.Errorf("storage: get product by id: %w", err)
	}

	// Fetch image URLs ordered by sort_order
	const imgQuery = `
		SELECT url
		FROM catalog.product_images
		WHERE product_id = $1
		ORDER BY sort_order ASC`

	imgRows, err := s.db.Query(ctx, imgQuery, productID)
	if err != nil {
		return ProductDetailRow{}, fmt.Errorf("storage: get product images: %w", err)
	}
	defer imgRows.Close()

	var images []string
	for imgRows.Next() {
		var url string
		if scanErr := imgRows.Scan(&url); scanErr != nil {
			return ProductDetailRow{}, fmt.Errorf("storage: scan image row: %w", scanErr)
		}
		images = append(images, url)
	}
	if err := imgRows.Err(); err != nil {
		return ProductDetailRow{}, fmt.Errorf("storage: iterate image rows: %w", err)
	}

	if images == nil {
		images = []string{}
	}
	row.Images = images

	return row, nil
}

// buildOrderClause converts a validated SortOption into a safe ORDER BY fragment.
// sort has already been validated against ValidSortOptions at the handler layer.
func buildOrderClause(sort SortOption) string {
	switch sort {
	case SortPriceAsc:
		return "ORDER BY p.price ASC, p.id ASC"
	case SortPriceDesc:
		return "ORDER BY p.price DESC, p.id ASC"
	case SortNameAsc:
		return "ORDER BY p.name ASC, p.id ASC"
	case SortNameDesc:
		return "ORDER BY p.name DESC, p.id ASC"
	default: // SortCreatedDesc is the default
		return "ORDER BY p.created_at DESC, p.id ASC"
	}
}
