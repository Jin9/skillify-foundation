package checkout

import (
	"encoding/json"
	"time"
)

// ---------------------------------------------------------------------------
// Idempotency
// ---------------------------------------------------------------------------

// IdempotencyStatus enumerates the lifecycle states of a checkout commit attempt.
type IdempotencyStatus string

const (
	StatusInflight   IdempotencyStatus = "INFLIGHT"
	StatusCompleted  IdempotencyStatus = "COMPLETED"
	StatusAbandoned  IdempotencyStatus = "ABANDONED"
)

// IdempotencyEntry is the persisted row in idempotency_keys.
type IdempotencyEntry struct {
	Key              string
	CustomerUserID   string
	RequestHash      string
	Status           IdempotencyStatus
	ResponseEnvelope json.RawMessage // NULL until COMPLETED|ABANDONED
	HTTPStatus       *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExpiresAt        time.Time
	FinalizedAt      *time.Time
}

// ---------------------------------------------------------------------------
// Preview
// ---------------------------------------------------------------------------

// PreviewRequest is the decoded request body for POST /checkout/checkout/preview.
// Strict body: only shippingAddressId is accepted. Any pricing or coupon field
// causes a 400 VALIDATION_ERROR per CHK-005 / CHK-AMBIG-002.
type PreviewRequest struct {
	ShippingAddressID string `json:"shippingAddressId"`
}

// PreviewItem is one line in the preview response.
type PreviewItem struct {
	ProductID    string  `json:"productId"`
	SKU          string  `json:"sku"`
	Qty          int     `json:"qty"`
	CurrentPrice float64 `json:"currentPrice"` // THB, whole units
	LineSubtotal float64 `json:"lineSubtotal"`
	AvailableQty int     `json:"availableQty"`
	Checkoutable bool    `json:"checkoutable"`
}

// AddressSnapshot is the resolved address embedded in preview and commit responses.
type AddressSnapshot struct {
	AddressID    string `json:"addressId"`
	FullName     string `json:"fullName"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	Province     string `json:"province"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country"`
	IsComplete   bool   `json:"isComplete"`
}

// PreviewResponse is the data field of the preview envelope.
type PreviewResponse struct {
	Items           []PreviewItem   `json:"items"`
	Subtotal        float64         `json:"subtotal"`
	ShippingFee     float64         `json:"shippingFee"`
	CouponDiscount  float64         `json:"couponDiscount"`
	Total           float64         `json:"total"`
	Blockers        []string        `json:"blockers"`
	AddressSnapshot AddressSnapshot `json:"addressSnapshot"`
}

// ---------------------------------------------------------------------------
// Commit
// ---------------------------------------------------------------------------

// CommitRequest is the decoded request body for POST /checkout/checkout/commit.
// Server is the sole pricing authority (CHK-005).
type CommitRequest struct {
	ShippingAddressID string `json:"shippingAddressId"`
}

// PaymentIntentResult is the nested paymentIntent object in CommitResponse.
type PaymentIntentResult struct {
	PaymentIntentID string  `json:"paymentIntentId"`
	Status          string  `json:"status"` // always "REQUIRES_PAYMENT"
	Amount          float64 `json:"amount"` // THB whole units
}

// CommitResponse is the data field of the commit envelope (HTTP 201).
type CommitResponse struct {
	OrderID        string              `json:"orderId"`
	OrderNumber    string              `json:"orderNumber"`
	Status         string              `json:"status"` // always "PENDING_PAYMENT"
	Subtotal       float64             `json:"subtotal"`
	ShippingFee    float64             `json:"shippingFee"`
	CouponDiscount float64             `json:"couponDiscount"`
	Total          float64             `json:"total"`
	PaymentIntent  PaymentIntentResult `json:"paymentIntent"`
}

// ---------------------------------------------------------------------------
// Internal line-item snapshot (used during orchestration, not serialised to client)
// ---------------------------------------------------------------------------

// LineItem holds the server-side price snapshot for one cart line.
type LineItem struct {
	ProductID             string
	SKU                   string
	CartItemID            string
	Qty                   int
	PriceSnapshotMinor    int64  // catalog price in THB minor units (1 THB = 100 minor)
	ProductName           string
	ProductImageURL       string
}
