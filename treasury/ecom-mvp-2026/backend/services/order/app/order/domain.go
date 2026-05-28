package order

import (
	"encoding/json"
	"time"
)

// OrderStatus is the set of valid states for an order.
// Order is the SOLE WRITER of orders.status — no other service may UPDATE this column.
type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusPaymentFailed  OrderStatus = "PAYMENT_FAILED"
	StatusPaymentExpired OrderStatus = "PAYMENT_EXPIRED"
	StatusPacking        OrderStatus = "PACKING"
	StatusShipped        OrderStatus = "SHIPPED"
	StatusDelivered      OrderStatus = "DELIVERED"
	StatusCancelled      OrderStatus = "CANCELLED"
)

// ActorRole identifies who drives a state transition.
type ActorRole string

const (
	RoleCustomer      ActorRole = "CUSTOMER"
	RoleAdmin         ActorRole = "ADMIN"
	RoleSystem        ActorRole = "SYSTEM"
	RolePaymentSystem ActorRole = "PAYMENT_SYSTEM" // alias for event-driven SYSTEM transitions recorded in history
)

// AddressSnapshot is the frozen address captured at order-creation time.
// It is stored as JSONB in orders.address_snapshot and never re-resolved.
type AddressSnapshot struct {
	ReceiverName string `json:"receiverName"`
	Phone        string `json:"phone"`
	Line         string `json:"line"`
	Province     string `json:"province"`
	District     string `json:"district"`
	PostalCode   string `json:"postalCode"`
}

// Order is the aggregate root for the order bounded context.
type Order struct {
	ID                 string          `db:"id"`
	OrderNumber        string          `db:"order_number"`
	UserID             string          `db:"user_id"`
	Status             OrderStatus     `db:"status"`
	Subtotal           int64           `db:"subtotal"`
	ShippingFee        int64           `db:"shipping_fee"`
	CouponDiscount     int64           `db:"coupon_discount"`
	GrandTotal         int64           `db:"grand_total"`
	Currency           string          `db:"currency"`
	AddressSnapshot    json.RawMessage `db:"address_snapshot"`
	BuyerEmailSnapshot string          `db:"buyer_email_snapshot"`
	TrackingNumber     *string         `db:"tracking_number"`
	IdempotencyKey     *string         `db:"idempotency_key"`
	Version            int             `db:"version"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}

// OrderItem is immutable per ORD-007. The repo layer exposes no Update method.
type OrderItem struct {
	ID              string  `db:"id"`
	OrderID         string  `db:"order_id"`
	ProductID       string  `db:"product_id"`
	SKU             string  `db:"sku"`
	Qty             int     `db:"qty"`
	NameSnapshot    string  `db:"name_snapshot"`
	ImageURLSnapshot *string `db:"image_url_snapshot"`
	PriceSnapshot   int64   `db:"price_snapshot"`
	LineSubtotal    int64   `db:"line_subtotal"`
	CreatedAt       time.Time `db:"created_at"`
}

// OrderStatusHistory is the append-only audit log per ORD-009.
// Written in the SAME transaction as the orders.status UPDATE.
type OrderStatusHistory struct {
	ID          string      `db:"id"`
	OrderID     string      `db:"order_id"`
	FromStatus  *OrderStatus `db:"from_status"` // NULL for the initial PENDING_PAYMENT row
	ToStatus    OrderStatus  `db:"to_status"`
	ActorUserID *string     `db:"actor_user_id"` // NULL when actor_role=SYSTEM
	ActorRole   ActorRole   `db:"actor_role"`
	Reason      *string     `db:"reason"`
	OccurredAt  time.Time   `db:"occurred_at"`
}

// OutboxEvent represents a transactional outbox row for events.order.cancelled.
// Written in the SAME transaction as the status mutation.
type OutboxEvent struct {
	ID          string          `db:"id"`
	AggregateID string          `db:"aggregate_id"` // = orders.id, used as Kafka partition key
	EventType   string          `db:"event_type"`   // "order.cancelled"
	PayloadJSON json.RawMessage `db:"payload_json"`
	CreatedAt   time.Time       `db:"created_at"`
	PublishedAt *time.Time      `db:"published_at"` // NULL until outbox publisher ships to Kafka
}

// OrderCancelledPayload conforms to events.order.cancelled payload shape.
type OrderCancelledPayload struct {
	EventID     string  `json:"eventId"`
	OccurredAt  string  `json:"occurredAt"` // RFC3339
	OrderID     string  `json:"orderId"`
	CancelActor string  `json:"cancelActor"` // CUSTOMER | ADMIN | SYSTEM
	FromStatus  string  `json:"fromStatus"`  // PENDING_PAYMENT | PAID
	Reason      *string `json:"reason,omitempty"`
}

// ConsumedEvent is the per-service idempotency record for inbound Kafka events.
// PK on event_id absorbs Kafka redelivery (layer 1 of two-layer dedup).
type ConsumedEvent struct {
	EventID      string    `db:"event_id"`
	ConsumerName string    `db:"consumer_name"`
	EventType    string    `db:"event_type"`
	OrderID      string    `db:"order_id"`
	ProcessedAt  time.Time `db:"processed_at"`
}

// --- Request / Response DTOs ---

// CreateFromCheckoutRequest is the inbound payload for order.create-from-checkout.
// The caller is Checkout service; fields are populated from checkout.commit's data.
type CreateFromCheckoutRequest struct {
	IdempotencyKey  string              `json:"idempotencyKey" binding:"required"`
	CustomerUserID  string              `json:"customerUserId" binding:"required,uuid"`
	BuyerEmail      string              `json:"buyerEmail" binding:"required,email"`
	AddressSnapshot AddressSnapshot     `json:"addressSnapshot" binding:"required"`
	Items           []CreateItemRequest `json:"items" binding:"required,min=1,dive"`
	Subtotal        int64               `json:"subtotal" binding:"min=0"`
	ShippingFee     int64               `json:"shippingFee" binding:"min=0"`
	CouponDiscount  int64               `json:"couponDiscount"`
	GrandTotal      int64               `json:"grandTotal" binding:"min=0"`
}

// CreateItemRequest is a single line-item in CreateFromCheckoutRequest.
type CreateItemRequest struct {
	ProductID        string  `json:"productId" binding:"required,uuid"`
	SKU              string  `json:"sku" binding:"required"`
	Qty              int     `json:"qty" binding:"required,min=1"`
	PriceSnapshot    int64   `json:"priceSnapshot" binding:"min=0"`
	NameSnapshot     string  `json:"nameSnapshot" binding:"required"`
	ImageURLSnapshot *string `json:"imageUrlSnapshot"`
	LineSubtotal     int64   `json:"lineSubtotal" binding:"min=0"`
}

// CreateFromCheckoutResponse is returned by order.create-from-checkout.
type CreateFromCheckoutResponse struct {
	OrderID     string `json:"orderId"`
	OrderNumber string `json:"orderNumber"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

// DetailRequest is the inbound payload for order.detail.
type DetailRequest struct {
	OrderID string `json:"orderId" binding:"required,uuid"`
}

// DetailResponse is the full order view returned by order.detail.
type DetailResponse struct {
	OrderID         string              `json:"orderId"`
	OrderNumber     string              `json:"orderNumber"`
	Status          OrderStatus         `json:"status"`
	Items           []ItemResponse      `json:"items"`
	AddressSnapshot AddressSnapshot     `json:"addressSnapshot"`
	Subtotal        int64               `json:"subtotal"`
	ShippingFee     int64               `json:"shippingFee"`
	CouponDiscount  int64               `json:"couponDiscount"`
	GrandTotal      int64               `json:"grandTotal"`
	Currency        string              `json:"currency"`
	TrackingNumber  *string             `json:"trackingNumber,omitempty"`
	StatusHistory   []HistoryResponse   `json:"statusHistory"`
	CreatedAt       string              `json:"createdAt"`
	UpdatedAt       string              `json:"updatedAt"`
}

// ItemResponse is a line item in DetailResponse, including all ORD-007 snapshots.
type ItemResponse struct {
	ItemID           string  `json:"itemId"`
	ProductID        string  `json:"productId"`
	SKU              string  `json:"sku"`
	Qty              int     `json:"qty"`
	NameSnapshot     string  `json:"nameSnapshot"`
	ImageURLSnapshot *string `json:"imageUrlSnapshot,omitempty"`
	PriceSnapshot    int64   `json:"priceSnapshot"`
	LineSubtotal     int64   `json:"lineSubtotal"`
}

// HistoryResponse is a single status-history entry in DetailResponse.
type HistoryResponse struct {
	FromStatus  *OrderStatus `json:"fromStatus"`
	ToStatus    OrderStatus  `json:"toStatus"`
	ActorRole   ActorRole    `json:"actorRole"`
	ActorUserID *string      `json:"actorUserId,omitempty"`
	Reason      *string      `json:"reason,omitempty"`
	OccurredAt  string       `json:"occurredAt"`
}

// ListMineRequest is the inbound payload for order.list-mine.
type ListMineRequest struct {
	Page   int         `json:"page"`
	Limit  int         `json:"limit"`
	Status *OrderStatus `json:"status,omitempty"`
}

// OrderSummary is a single row in list responses.
type OrderSummary struct {
	OrderID     string      `json:"orderId"`
	OrderNumber string      `json:"orderNumber"`
	Status      OrderStatus `json:"status"`
	GrandTotal  int64       `json:"grandTotal"`
	Currency    string      `json:"currency"`
	CreatedAt   string      `json:"createdAt"`
}

// ListMineResponse is returned by order.list-mine.
type ListMineResponse struct {
	Orders []OrderSummary `json:"orders"`
	Total  int            `json:"total"`
	Page   int            `json:"page"`
}
