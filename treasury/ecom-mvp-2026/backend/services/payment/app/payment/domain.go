package payment

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ─── Payment Intent ───────────────────────────────────────────────────────────

// IntentStatus is the lifecycle state of a PaymentIntent row.
type IntentStatus string

const (
	StatusRequiresPayment IntentStatus = "REQUIRES_PAYMENT"
	StatusSucceeded       IntentStatus = "SUCCEEDED"
	StatusFailed          IntentStatus = "FAILED"
	StatusExpired         IntentStatus = "EXPIRED"
)

// TerminalStatuses is the set of states from which no further transitions are
// allowed unless via the same dedup key (which would be a REPLAY_HIT).
var TerminalStatuses = map[IntentStatus]bool{
	StatusSucceeded: true,
	StatusFailed:    true,
	StatusExpired:   true,
}

// PaymentIntent is the single row in payment_intents for one order.
// Columns match persistence.tables[0] in td.json.
type PaymentIntent struct {
	IntentID        uuid.UUID    `db:"intent_id"`
	OrderID         uuid.UUID    `db:"order_id"`
	OwnerUserID     uuid.UUID    `db:"owner_user_id"`
	AmountMinor     int64        `db:"amount_minor"`
	Currency        string       `db:"currency"`
	Status          IntentStatus `db:"status"`
	MockProviderRef *string      `db:"mock_provider_ref"`
	ProviderStatus  *string      `db:"provider_status"`
	PaidAt          *time.Time   `db:"paid_at"`
	ExpiresAt       time.Time    `db:"expires_at"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at"`
	Version         int          `db:"version"`
}

// ─── Callback Dedup ───────────────────────────────────────────────────────────

// CallbackDedup is a row in payment_callback_dedup.
// PK = sha256(intent_id || '|' || provider_status).
// The envelope is the exact HTTP response returned to the original caller.
type CallbackDedup struct {
	DedupKey       string          `db:"dedup_key"`
	IntentID       uuid.UUID       `db:"intent_id"`
	ProviderStatus string          `db:"provider_status"`
	Envelope       json.RawMessage `db:"envelope"`
	HTTPStatus     int             `db:"http_status"`
	CreatedAt      time.Time       `db:"created_at"`
	ExpiresAt      time.Time       `db:"expires_at"`
}

// ─── Outbox ───────────────────────────────────────────────────────────────────

// OutboxEvent is a row in the payment.outbox table.
// Written in the SAME tx as the intent state transition.
// A separate publisher goroutine drains PENDING rows to Kafka.
type OutboxEvent struct {
	ID            uuid.UUID       `db:"id"`
	AggregateType string          `db:"aggregate_type"` // "payment_intent"
	AggregateID   uuid.UUID       `db:"aggregate_id"`   // intent_id
	EventID       uuid.UUID       `db:"event_id"`       // UUID v4, for consumer dedup
	EventType     string          `db:"event_type"`     // "payment.completed" | "payment.failed"
	Topic         string          `db:"topic"`          // "ecom.payment.events"
	PartitionKey  string          `db:"partition_key"`  // order_id (string)
	Payload       json.RawMessage `db:"payload"`
	Status        string          `db:"status"`      // "PENDING" | "PUBLISHED" | "FAILED"
	CreatedAt     time.Time       `db:"created_at"`
	PublishedAt   *time.Time      `db:"published_at"`
	Attempts      int             `db:"attempts"`
	LastError     *string         `db:"last_error"`
}

// ─── Outbox Payloads ──────────────────────────────────────────────────────────

// PaymentCompletedPayload matches events.payment.completed contract.
type PaymentCompletedPayload struct {
	EventID          string    `json:"eventId"`
	OccurredAt       time.Time `json:"occurredAt"`
	OrderID          string    `json:"orderId"`
	PaymentIntentID  string    `json:"paymentIntentId"`
	Amount           int64     `json:"amount"`
	MockPaymentRef   string    `json:"mockPaymentRef"`
}

// PaymentFailedPayload matches events.payment.failed contract.
type PaymentFailedPayload struct {
	EventID         string    `json:"eventId"`
	OccurredAt      time.Time `json:"occurredAt"`
	OrderID         string    `json:"orderId"`
	PaymentIntentID string    `json:"paymentIntentId"`
	Reason          string    `json:"reason"` // "PROVIDER_DECLINED"
}

// ─── Response Envelopes (stored in dedup rows) ────────────────────────────────

// StoredEnvelope is the shape persisted in payment_callback_dedup.envelope.
// On a REPLAY_HIT the stored HTTP status + body are returned verbatim.
type StoredEnvelope struct {
	HTTPStatus int             `json:"httpStatus"`
	Body       json.RawMessage `json:"body"`
}

// ─── Request / Response shapes ────────────────────────────────────────────────

// IntentCreateRequest is the body for POST /api/v1/payment/intent/create.
type IntentCreateRequest struct {
	OrderID     string    `json:"orderId"     binding:"required,uuid"`
	OwnerUserID string    `json:"ownerUserId" binding:"required,uuid"`
	Amount      int64     `json:"amount"      binding:"required,min=1"`
	Currency    string    `json:"currency"    binding:"omitempty,len=3"`
}

// IntentCreateResponse is the data field of the create response envelope.
type IntentCreateResponse struct {
	PaymentIntentID string    `json:"paymentIntentId"`
	Status          string    `json:"status"`
	Amount          int64     `json:"amount"`
	ExpiresAt       time.Time `json:"expiresAt"`
}

// SimulateRequest is the body for POST /api/v1/payment/intent/simulate.
type SimulateRequest struct {
	PaymentIntentID string `json:"paymentIntentId" binding:"required,uuid"`
	Outcome         string `json:"outcome"         binding:"required,oneof=success failed timeout"`
}

// SimulateResponse is the data field of the simulate response.
type SimulateResponse struct {
	PaymentIntentID    string `json:"paymentIntentId"`
	ProviderStatus     string `json:"providerStatus"`
	MockPaymentRef     string `json:"mockPaymentRef"`
	WillTriggerCallback bool  `json:"willTriggerCallback"`
}

// CallbackRequest is the body for POST /api/v1/payment/intent/callback.
type CallbackRequest struct {
	PaymentIntentID   string    `json:"paymentIntentId"   binding:"required,uuid"`
	ProviderStatus    string    `json:"providerStatus"    binding:"required,oneof=SUCCEEDED FAILED EXPIRED"`
	MockPaymentRef    string    `json:"mockPaymentRef"    binding:"required,max=64"`
	Amount            int64     `json:"amount"            binding:"required,min=1"`
	ProviderTimestamp time.Time `json:"providerTimestamp" binding:"required"`
}

// CallbackResponse is the data field of the callback success response.
type CallbackResponse struct {
	PaymentIntentID string `json:"paymentIntentId"`
	Applied         bool   `json:"applied"`
	ProviderStatus  string `json:"providerStatus"`
}

// ProcessCallbackInput is the internal argument to process_callback.
type ProcessCallbackInput struct {
	IntentID          uuid.UUID
	ProviderStatus    string // "SUCCEEDED" | "FAILED" | "EXPIRED"
	MockPaymentRef    string
	AmountFromCaller  int64
	ProviderTimestamp time.Time
}

// ProcessCallbackOutput is the result of process_callback (returned to callers).
type ProcessCallbackOutput struct {
	HTTPStatus int
	Body       json.RawMessage
}

// ─── Errors ───────────────────────────────────────────────────────────────────

// Sentinel errors used in service / storage layers.
// Handlers switch on these to produce the correct HTTP error envelope.
var (
	ErrIntentNotFound      = errSentinel("PAYMENT_INTENT_NOT_FOUND")
	ErrAmountMismatch      = errSentinel("PAYMENT_AMOUNT_MISMATCH")
	ErrIntentTerminal      = errSentinel("PAYMENT_INTENT_TERMINAL")
	ErrAuthForbidden       = errSentinel("AUTH_FORBIDDEN")
	ErrAuthInvalid         = errSentinel("AUTH_INVALID")
	ErrConflict            = errSentinel("CONFLICT")
	ErrTimeoutNotSupported = errSentinel("TIMEOUT_OUTCOME_NOT_SUPPORTED")
	// ErrDuplicateOrderID is returned by IntentStorage.Insert when a row for
	// that order_id already exists. Defined here so the service layer can check
	// it without importing the access package (breaks import cycle).
	ErrDuplicateOrderID = errSentinel("DUPLICATE_ORDER_ID")
)

type sentinel string

func errSentinel(s string) error { return sentinel(s) }
func (s sentinel) Error() string { return string(s) }
