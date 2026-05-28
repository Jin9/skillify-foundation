package inventory

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ─── Stock ───────────────────────────────────────────────────────────────────

// StockLevel represents a single SKU's inventory state.
// Columns: sku (PK), available_qty, reserved_qty, sold_qty, version, updated_at.
// CHECK constraints in schema enforce no-negative invariant at the DB layer
// (belt-and-braces; the service enforces it in-tx first).
type StockLevel struct {
	SKU          string    `db:"sku"`
	AvailableQty int       `db:"available_qty"`
	ReservedQty  int       `db:"reserved_qty"`
	SoldQty      int       `db:"sold_qty"`
	Version      int       `db:"version"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// StockAdjustment is the audit row written in the same tx as inventory.stock.adjust.
type StockAdjustment struct {
	ID          uuid.UUID `db:"id"`
	SKU         string    `db:"sku"`
	Delta       int       `db:"delta"`
	Reason      string    `db:"reason"`
	ActorUserID uuid.UUID `db:"actor_user_id"`
	CreatedAt   time.Time `db:"created_at"`
}

// ─── Reservation ─────────────────────────────────────────────────────────────

// ReservationStatus is the lifecycle state of a Reservation row.
type ReservationStatus string

const (
	StatusReserved  ReservationStatus = "RESERVED"
	StatusCommitted ReservationStatus = "COMMITTED"
	StatusReleased  ReservationStatus = "RELEASED"
	StatusExpired   ReservationStatus = "EXPIRED"
)

// ReleaseReason documents why a reservation was released/expired.
type ReleaseReason string

const (
	ReasonPaymentFailed   ReleaseReason = "PAYMENT_FAILED"
	ReasonPaymentExpired  ReleaseReason = "PAYMENT_EXPIRED"
	ReasonOrderCancelled  ReleaseReason = "ORDER_CANCELLED"
	ReasonAdminForce      ReleaseReason = "ADMIN_FORCE"
	ReasonExpired         ReleaseReason = "EXPIRED"
	ReasonPaidSoftCancel  ReleaseReason = "PAID_SOFT_CANCEL"
)

// Reservation is a single SKU reservation row.
// One Reservation is created per SKU per order; they share the same order_id.
type Reservation struct {
	ID            uuid.UUID          `db:"id"`
	OrderID       uuid.UUID          `db:"order_id"`
	SKU           string             `db:"sku"`
	Qty           int                `db:"qty"`
	Status        ReservationStatus  `db:"status"`
	ExpiresAt     time.Time          `db:"expires_at"`
	CreatedAt     time.Time          `db:"created_at"`
	ReleasedAt    *time.Time         `db:"released_at"`
	ReleaseReason *ReleaseReason     `db:"release_reason"`
}

// ─── Outbox ──────────────────────────────────────────────────────────────────

// OutboxEvent is a row in outbox_events. The relayer goroutine polls for rows
// where published_at IS NULL and publishes them to Kafka at-least-once.
type OutboxEvent struct {
	ID          uuid.UUID       `db:"id"`
	AggregateID string          `db:"aggregate_id"` // = orderId for reservation.expired
	EventType   string          `db:"event_type"`
	PayloadJSON json.RawMessage `db:"payload_json"`
	CreatedAt   time.Time       `db:"created_at"`
	PublishedAt *time.Time      `db:"published_at"`
}

// ReservationExpiredPayload is the payload written to the outbox when the sweeper
// transitions a reservation to EXPIRED. Shape per contracts.json:events.reservation.expired.
type ReservationExpiredPayload struct {
	EventID       string    `json:"eventId"`
	OccurredAt    time.Time `json:"occurredAt"`
	OrderID       string    `json:"orderId"`
	ReservationID string    `json:"reservationId"`
	SweptAt       time.Time `json:"sweptAt"`
}

// ─── Consumed Events (consumer dedup) ────────────────────────────────────────

// ConsumedEvent is a dedup row inserted inside the consumer's tx.
// PK is (event_id, consumer_name); ON CONFLICT DO NOTHING → idempotent skip.
type ConsumedEvent struct {
	EventID      uuid.UUID `db:"event_id"`
	ConsumerName string    `db:"consumer_name"`
	ProcessedAt  time.Time `db:"processed_at"`
}

// ─── Idempotency ─────────────────────────────────────────────────────────────

// IdempotencyKey is the server-side idempotency store for reservation.create.
// Key = orderId, endpoint = "inventory.reservation.create".
type IdempotencyKey struct {
	Key              string          `db:"key"`
	Endpoint         string          `db:"endpoint"`
	RequestHash      string          `db:"request_hash"`
	ResponseEnvelope json.RawMessage `db:"response_envelope"`
	HTTPStatus       int             `db:"http_status"`
	CreatedAt        time.Time       `db:"created_at"`
	ExpiresAt        time.Time       `db:"expires_at"`
}

// ─── Request/Response shapes ─────────────────────────────────────────────────

// ReservationItem is a {sku, qty} pair in the create request.
type ReservationItem struct {
	SKU string `json:"sku" binding:"required"`
	Qty int    `json:"qty" binding:"required,min=1"`
}

// ReservationCreateRequest is the body for POST /api/v1/inventory/reservation/create.
type ReservationCreateRequest struct {
	OrderID        string            `json:"orderId" binding:"required,uuid"`
	Items          []ReservationItem `json:"items" binding:"required,min=1"`
	ExpiresAt      time.Time         `json:"expiresAt"`
	IdempotencyKey string            `json:"idempotencyKey" binding:"required"`
}

// ReservationCreateResponseItem is one item in the create response.
type ReservationCreateResponseItem struct {
	SKU      string `json:"sku"`
	Qty      int    `json:"qty"`
	Reserved bool   `json:"reserved"`
}

// ReservationCreateResponse is the data field of the create response envelope.
type ReservationCreateResponse struct {
	ReservationID string                          `json:"reservationId"`
	Items         []ReservationCreateResponseItem `json:"items"`
	ExpiresAt     time.Time                       `json:"expiresAt"`
}

// InsufficientDetail is per-SKU detail returned on INSUFFICIENT_STOCK.
type InsufficientDetail struct {
	SKU          string `json:"sku"`
	RequestedQty int    `json:"requestedQty"`
	AvailableQty int    `json:"availableQty"`
}

// StockReadRequest is the body for POST /api/v1/inventory/stock/read.
type StockReadRequest struct {
	SKU string `json:"sku" binding:"required"`
}

// StockReadResponse is the data field of the read response.
type StockReadResponse struct {
	SKU          string `json:"sku"`
	AvailableQty int    `json:"availableQty"`
	ReservedQty  int    `json:"reservedQty"`
	SoldQty      int    `json:"soldQty"`
}

// ReservationReleaseRequest is the body for POST /api/v1/inventory/reservation/release.
type ReservationReleaseRequest struct {
	OrderID string `json:"orderId" binding:"required,uuid"`
	Reason  string `json:"reason" binding:"required,oneof=PAYMENT_FAILED PAYMENT_EXPIRED ORDER_CANCELLED ADMIN_FORCE"`
}

// ReservationReleaseResponse is the data field of the release response.
type ReservationReleaseResponse struct {
	ReservationID string `json:"reservationId"`
	Status        string `json:"status"`
}
