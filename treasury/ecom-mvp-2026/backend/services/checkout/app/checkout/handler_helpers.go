package checkout

// handler_helpers.go — shared utilities used by both handler_preview.go and handler_commit.go.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// ---------------------------------------------------------------------------
// Application-level codes for checkout
// ---------------------------------------------------------------------------

// These extend wrapper.Code and align with cross-cutting.error-codes.registry.
const (
	CodeSuccess    = wrapper.CodeSuccess
	CodeCreated    wrapper.Code = "CREATED"

	CodeValidationError        wrapper.Code    = "VALIDATION_ERROR"
	CodeAuthMissing            wrapper.Code    = "AUTH_MISSING"
	CodeAuthForbidden          wrapper.Code    = "AUTH_FORBIDDEN"
	CodeAddressNotOwned        wrapper.Code    = "ADDRESS_NOT_OWNED"
	CodeAddressIncomplete      wrapper.Code    = "ADDRESS_INCOMPLETE"
	CodeInsufficientStock      wrapper.Code    = "INSUFFICIENT_STOCK"
	CodeProductInactive        wrapper.Code    = "PRODUCT_INACTIVE"
	CodeProductDeleted         wrapper.Code    = "PRODUCT_DELETED"
	CodeIdempotencyKeyReused   wrapper.Code    = "IDEMPOTENCY_KEY_REUSED"
	CodeIdempotencyKeyInflight wrapper.Code    = "IDEMPOTENCY_KEY_INFLIGHT"
	CodeDBUnavailable          wrapper.Code    = "DATABASE_UNAVAILABLE"
	CodeUpstreamTimeout        wrapper.Code    = "UPSTREAM_TIMEOUT"
)

// ---------------------------------------------------------------------------
// Forbidden field lists (CHK-005, CHK-AMBIG-002)
// ---------------------------------------------------------------------------

// previewForbiddenFields: any pricing or coupon field is rejected.
var previewForbiddenFields = []string{
	"total", "subtotal", "shippingFee", "shippingFeeMinor",
	"couponCode", "couponDiscount",
	"grandTotal", "currentPrice", "priceSnapshot",
}

// commitForbiddenFields: same plus commit-specific price fields.
var commitForbiddenFields = []string{
	"total", "subtotal", "shippingFee", "shippingFeeMinor",
	"couponCode", "couponDiscount",
	"grandTotal", "currentPrice", "priceSnapshot",
}

// rejectForbiddenFields returns an error if any forbidden key is present in body.
func rejectForbiddenFields(body map[string]json.RawMessage, forbidden []string) error {
	for _, f := range forbidden {
		if _, ok := body[f]; ok {
			if f == "couponCode" {
				return errors.New("couponCode is not supported in MVP")
			}
			return errors.New("client-supplied price/total fields are not accepted; server is authoritative (CHK-005)")
		}
	}
	return nil
}

// requireStringField extracts a required non-empty string from the raw body map.
func requireStringField(body map[string]json.RawMessage, field string, maxLen int) (string, error) {
	raw, ok := body[field]
	if !ok {
		return "", fmt.Errorf("%s is required", field)
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", fmt.Errorf("%s must be a string", field)
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("%s must not be empty", field)
	}
	if len(s) > maxLen {
		return "", fmt.Errorf("%s must not exceed %d characters", field, maxLen)
	}
	return s, nil
}

// ---------------------------------------------------------------------------
// Auth helpers
// ---------------------------------------------------------------------------

// requireCustomerClaims extracts JWT claims and enforces role=CUSTOMER.
func requireCustomerClaims(c *gin.Context) (*token.Claims, bool) {
	claims, err := token.ClaimsFromContext(c.Request.Context())
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       CodeAuthMissing,
			Message:    "authentication required",
		})
		c.Abort()
		return nil, false
	}

	role, _ := claims.Extra["role"].(string)
	if role != "CUSTOMER" {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusForbidden,
			Code:       CodeAuthForbidden,
			Message:    "CUSTOMER role required",
		})
		c.Abort()
		return nil, false
	}

	return claims, true
}

// extractBearer returns the raw Bearer token from the Authorization header.
func extractBearer(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	return strings.TrimPrefix(auth, "Bearer ")
}

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

func respondValidationError(c *gin.Context, msg string) {
	traceID := traceIDFromContext(c)
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusBadRequest,
		Code:       CodeValidationError,
		Message:    wrapper.Message(msg),
		TraceID:    traceID,
	})
}

func respondDBUnavailable(c *gin.Context, err error, traceID string) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusServiceUnavailable,
		Code:       CodeDBUnavailable,
		Message:    "database unavailable",
		TraceID:    traceID,
		Err:        err,
	})
}

// respondUpstreamError maps UpstreamError codes to appropriate HTTP responses.
func respondUpstreamError(c *gin.Context, err error, traceID string) {
	var u *UpstreamError
	if errors.As(err, &u) {
		switch u.Code {
		case "DATABASE_UNAVAILABLE":
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusServiceUnavailable,
				Code:       CodeDBUnavailable,
				Message:    "upstream database unavailable",
				TraceID:    traceID,
				Err:        err,
			})
			return
		}
	}
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusGatewayTimeout,
		Code:       CodeUpstreamTimeout,
		Message:    "upstream service timeout",
		TraceID:    traceID,
		Err:        err,
	})
}

// traceIDFromContext extracts the W3C trace ID injected by the traceparent middleware.
func traceIDFromContext(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
