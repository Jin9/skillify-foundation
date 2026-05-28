package inventory

import (
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// Code and Message are type aliases from the wrapper package.
type Code = wrapper.Code
type Message = wrapper.Message

// Re-export common codes for convenience.
const (
	CodeSuccess    = wrapper.CodeSuccess
	MessageSuccess = wrapper.MessageSuccess

	CodeBadRequest    = wrapper.CodeBadRequest
	MessageBadRequest = wrapper.MessageBadRequest

	CodeNotFound    = wrapper.CodeNotFound
	MessageNotFound = wrapper.MessageNotFound

	CodeInternalError     = wrapper.CodeInternalError
	MessageInternalError  = wrapper.MessageInternalError

	CodeServiceUnavail    = wrapper.CodeServiceUnavail
	MessageServiceUnavail = wrapper.MessageServiceUnavail
)

// Inventory domain-specific error codes per contracts.json:cross-cutting.error-codes.
const (
	CodeValidationError         Code    = "VALIDATION_ERROR"
	MessageValidationError      Message = "Validation error"

	CodeInsufficientStock         Code    = "INSUFFICIENT_STOCK"
	MessageInsufficientStock      Message = "One or more items have insufficient stock"

	CodeStockNegativeInvariant         Code    = "STOCK_NEGATIVE_INVARIANT"
	MessageStockNegativeInvariant      Message = "Stock negative invariant breached"

	CodeIdempotencyKeyReused         Code    = "IDEMPOTENCY_KEY_REUSED"
	MessageIdempotencyKeyReused      Message = "Idempotency key reused with different request body"

	CodeReservationNotFound         Code    = "RESERVATION_NOT_FOUND"
	MessageReservationNotFound      Message = "Reservation not found"

	CodeDatabaseUnavailable         Code    = "DATABASE_UNAVAILABLE"
	MessageDatabaseUnavailable      Message = "Database unavailable"

	CodeNotImplementedMVP         Code    = "NOT_IMPLEMENTED_MVP"
	MessageNotImplementedMVP      Message = "Not implemented in MVP"

	CodeAuthForbidden         Code    = "AUTH_FORBIDDEN"
	MessageAuthForbidden      Message = "Forbidden"
)
