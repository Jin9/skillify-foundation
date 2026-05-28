package catalog

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the lifecycle state of a product.
type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "DRAFT"
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
	ProductStatusDeleted  ProductStatus = "DELETED"
)

// StockStatus represents inventory availability returned on product.detail.
type StockStatus string

const (
	StockStatusInStock    StockStatus = "IN_STOCK"
	StockStatusLow        StockStatus = "LOW"
	StockStatusOutOfStock StockStatus = "OUT_OF_STOCK"
)

// Category is a product classification node (domain entity).
type Category struct {
	ID               uuid.UUID  `json:"categoryId"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	ParentCategoryID *uuid.UUID `json:"parentCategoryId,omitempty"`
	Level            int        `json:"level"`
	Active           bool       `json:"active"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// OutboxEvent is the persisted outbox record used by the relay goroutine.
// published_at IS NULL means the event is pending dispatch.
type OutboxEvent struct {
	ID          uuid.UUID  `json:"id"`
	AggregateID uuid.UUID  `json:"aggregateId"`
	EventType   string     `json:"eventType"`
	PayloadJSON []byte     `json:"payloadJson"`
	CreatedAt   time.Time  `json:"createdAt"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

// ProductListFilter holds optional filter parameters for product.list.
// Defined here for use in the handler/service layer; mapped to access.ProductListFilter in service.
type ProductListFilter struct {
	CategoryID *uuid.UUID
	MinPrice   *string // NUMERIC string (e.g. "10.00")
	MaxPrice   *string
	InStock    *bool // handled post-SQL via inventory client stub
	Search     *string
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

// ValidSortOptions is the allowed set for validation in the handler layer.
var ValidSortOptions = map[SortOption]struct{}{
	SortPriceAsc:    {},
	SortPriceDesc:   {},
	SortNameAsc:     {},
	SortNameDesc:    {},
	SortCreatedDesc: {},
}
