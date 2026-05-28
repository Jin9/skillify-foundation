package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/shoppilot/catalog/app/catalog/access"
)

// CatalogService exposes all product and category business operations.
type CatalogService struct {
	db            *pgxpool.Pool
	productStore  *access.ProductStorage
	categoryStore *access.CategoryStorage
	imageStore    *access.ImageStorage
	outboxStore   *access.OutboxStorage
}

// NewCatalogService constructs the service with its access-layer dependencies.
func NewCatalogService(
	db *pgxpool.Pool,
	productStore *access.ProductStorage,
	categoryStore *access.CategoryStorage,
	imageStore *access.ImageStorage,
	outboxStore *access.OutboxStorage,
) *CatalogService {
	return &CatalogService{
		db:            db,
		productStore:  productStore,
		categoryStore: categoryStore,
		imageStore:    imageStore,
		outboxStore:   outboxStore,
	}
}

// ListProductsResult is the paginated result for product.list.
type ListProductsResult struct {
	Items []ProductListItem
	Total int
	Page  int
	Limit int
}

// ProductListItem is the projected view returned in product.list.
type ProductListItem struct {
	ProductID    uuid.UUID
	SKU          string
	Name         string
	Price        string
	CategoryID   uuid.UUID
	CategoryName string
	CategorySlug string
	ThumbnailURL string // first image URL, empty string if none
	Status       string
}

// ProductDetail is the full product record returned by product.detail.
type ProductDetail struct {
	ProductID    uuid.UUID
	SKU          string
	Name         string
	Description  string
	Images       []string
	Price        string
	CategoryID   uuid.UUID
	CategoryName string
	CategorySlug string
	Status       string
}

// CreateProductInput carries validated data for product.create.
type CreateProductInput struct {
	SKU         string
	Name        string
	Description string
	Images      []string
	Price       string // validated NUMERIC string
	CategoryID  uuid.UUID
	Status      ProductStatus
	ActorUserID uuid.UUID
}

// ListProducts queries products with optional filters and pagination.
func (s *CatalogService) ListProducts(
	ctx context.Context,
	filter ProductListFilter,
	sort SortOption,
	page, limit int,
) (ListProductsResult, error) {
	// Map domain filter to access-layer filter (no catalog import in access)
	storeFilter := access.ProductListFilter{
		CategoryID: filter.CategoryID,
		MinPrice:   filter.MinPrice,
		MaxPrice:   filter.MaxPrice,
		Search:     filter.Search,
	}

	rows, total, err := s.productStore.List(ctx, storeFilter, access.SortOption(sort), page, limit)
	if err != nil {
		return ListProductsResult{}, fmt.Errorf("catalog: list products: %w", err)
	}

	items := make([]ProductListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, ProductListItem{
			ProductID:    r.ID,
			SKU:          r.SKU,
			Name:         r.Name,
			Price:        r.Price,
			CategoryID:   r.CategoryID,
			CategoryName: r.CategoryName,
			CategorySlug: r.CategorySlug,
			ThumbnailURL: r.ThumbnailURL,
			Status:       r.Status,
		})
	}

	// TODO: wire inventory.stock.bulk-read once Inventory service is up.
	// When filter.InStock != nil && *filter.InStock == true:
	//   skus := extractSKUs(items)
	//   stockMap, err := client_inventory.BulkReadStock(ctx, skus)  // returns map[sku]availableQty
	//   items = filterInStock(items, stockMap)
	//   On inventory 503/timeout: return 504 UPSTREAM_TIMEOUT (do NOT return unfiltered list)

	return ListProductsResult{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetProduct fetches a single product by ID, enforcing visibility rules based on role.
func (s *CatalogService) GetProduct(ctx context.Context, productID uuid.UUID, isAdmin bool) (ProductDetail, error) {
	row, err := s.productStore.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, access.ErrProductNotFound) {
			return ProductDetail{}, ErrProductNotFound
		}
		return ProductDetail{}, fmt.Errorf("catalog: get product: %w", err)
	}

	// DELETED and DRAFT products are hidden from non-admin callers.
	// INACTIVE products are returned (for cart compatibility per TD edge_behavior).
	if !isAdmin {
		if row.Status == string(ProductStatusDeleted) || row.Status == string(ProductStatusDraft) {
			return ProductDetail{}, ErrProductNotFound
		}
	}

	return ProductDetail{
		ProductID:    row.ID,
		SKU:          row.SKU,
		Name:         row.Name,
		Description:  row.Description,
		Images:       row.Images,
		Price:        row.Price,
		CategoryID:   row.CategoryID,
		CategoryName: row.CategoryName,
		CategorySlug: row.CategorySlug,
		Status:       row.Status,
	}, nil
}

// CreateProduct inserts a new product and an outbox event in a single transaction.
func (s *CatalogService) CreateProduct(ctx context.Context, input CreateProductInput) (uuid.UUID, error) {
	productID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog: generate product id: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog: begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// 1. Insert product
	if err = s.productStore.Insert(ctx, tx, access.InsertProductParams{
		ID:          productID,
		SKU:         input.SKU,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		Status:      string(input.Status),
	}); err != nil {
		if errors.Is(err, access.ErrDuplicateSKU) {
			return uuid.Nil, ErrDuplicateSKU
		}
		if errors.Is(err, access.ErrCategoryNotFound) {
			return uuid.Nil, ErrCategoryNotFound
		}
		return uuid.Nil, fmt.Errorf("catalog: insert product: %w", err)
	}

	// 2. Insert product images (ordered batch)
	for i, imgURL := range input.Images {
		if err = s.imageStore.Insert(ctx, tx, access.InsertImageParams{
			ProductID: productID,
			URL:       imgURL,
			SortOrder: i,
		}); err != nil {
			return uuid.Nil, fmt.Errorf("catalog: insert image %d: %w", i, err)
		}
	}

	// 3. Insert product_status_history (from_status=NULL → initial creation)
	if err = s.productStore.InsertStatusHistory(ctx, tx, access.InsertStatusHistoryParams{
		ProductID:   productID,
		FromStatus:  nil, // NULL = initial creation
		ToStatus:    string(input.Status),
		ActorUserID: input.ActorUserID,
	}); err != nil {
		return uuid.Nil, fmt.Errorf("catalog: insert status history: %w", err)
	}

	// 4. Build outbox payload — per TD events.produced spec
	eventID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog: generate event id: %w", err)
	}

	payload := map[string]any{
		"eventId":    eventID.String(),
		"occurredAt": time.Now().UTC().Format(time.RFC3339),
		"productId":  productID.String(),
		"sku":        input.SKU,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog: marshal outbox payload: %w", err)
	}

	outboxID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog: generate outbox id: %w", err)
	}

	// 5. Insert outbox event (same tx — both commit or neither)
	if err = s.outboxStore.Insert(ctx, tx, access.InsertOutboxParams{
		ID:          outboxID,
		AggregateID: productID,
		EventType:   "product.created",
		PayloadJSON: payloadBytes,
	}); err != nil {
		return uuid.Nil, fmt.Errorf("catalog: insert outbox: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("catalog: commit tx: %w", err)
	}

	slog.Info("product created", "product_id", productID, "sku", input.SKU, "status", input.Status)
	return productID, nil
}
