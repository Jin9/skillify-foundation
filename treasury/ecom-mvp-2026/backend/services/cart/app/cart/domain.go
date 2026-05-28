package cart

import (
	"time"

	"github.com/google/uuid"
)

// Cart is the aggregate root — one per customer.
type Cart struct {
	CartID    uuid.UUID
	UserID    uuid.UUID
	Items     []CartItem
	UpdatedAt time.Time
}

// CartItem is a persisted line in cart_items.
type CartItem struct {
	CartItemID uuid.UUID
	CartID     uuid.UUID
	ProductID  uuid.UUID
	Qty        int
	AddedAt    time.Time
}

// ProductStatus mirrors the catalog service's product status enum.
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
	ProductStatusDeleted  ProductStatus = "DELETED"
)

// CartItemView is an enriched line — DB row merged with live catalog + inventory data.
type CartItemView struct {
	CartItemID    uuid.UUID     `json:"cartItemId"`
	ProductID     uuid.UUID     `json:"productId"`
	SKU           string        `json:"sku"`
	Name          string        `json:"name"`
	Image         string        `json:"image"`
	CurrentPrice  float64       `json:"currentPrice"`
	Qty           int           `json:"qty"`
	LineSubtotal  float64       `json:"lineSubtotal"`
	AvailableQty  int           `json:"availableQty"`
	ProductStatus ProductStatus `json:"productStatus"`
	Checkoutable  bool          `json:"checkoutable"`
}

// CartView is the full enriched response shape for cart.read (and returned by write operations).
type CartView struct {
	CartID         uuid.UUID      `json:"cartId"`
	CustomerUserID uuid.UUID      `json:"customerUserId"`
	Items          []CartItemView `json:"items"`
	Subtotal       float64        `json:"subtotal"`
	Currency       string         `json:"currency"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
