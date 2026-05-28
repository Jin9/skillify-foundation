package payment

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// Callback handles POST /api/v1/payment/intent/callback.
//
// Auth: HMAC-SHA256 signature verification against CALLBACK_HMAC_SECRET.
//   - Header X-Mock-Provider-Signature: hex(HMAC-SHA256(rawBody, secret))
//   - Header X-Mock-Provider-Timestamp: RFC3339 (reject if skew > 5 min)
//
// On valid signature → parse body → invoke process_callback.
// Returns the stored envelope verbatim (200 or 409, same structure on dedup hit).
//
// SECURITY: rawBody is read BEFORE any JSON parsing to ensure the HMAC is
// computed over the exact wire bytes. Full callback body is NEVER logged.
func (h *handler) Callback(c *gin.Context) {
	// ── 1. Read raw body for HMAC (must happen before any JSON parse) ─────────
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       wrapper.CodeBadRequest,
			Message:    msgOf("Failed to read request body."),
		})
		return
	}

	// ── 2. Timestamp skew check (±5 min) ─────────────────────────────────────
	tsHeader := c.GetHeader("X-Mock-Provider-Timestamp")
	if tsHeader == "" {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       codeOf("AUTH_INVALID"),
			Message:    msgOf("Missing X-Mock-Provider-Timestamp header."),
		})
		return
	}
	providerTimestamp, err := time.Parse(time.RFC3339, tsHeader)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       codeOf("AUTH_INVALID"),
			Message:    msgOf("Invalid X-Mock-Provider-Timestamp; expected RFC3339."),
		})
		return
	}
	if skew := time.Since(providerTimestamp); skew > 5*time.Minute || skew < -5*time.Minute {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       codeOf("AUTH_INVALID"),
			Message:    msgOf("Timestamp skew exceeds ±5 minutes."),
		})
		return
	}

	// ── 3. HMAC-SHA256 signature verification (constant-time) ─────────────────
	sigHeader := c.GetHeader("X-Mock-Provider-Signature")
	if sigHeader == "" {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       codeOf("AUTH_INVALID"),
			Message:    msgOf("Missing X-Mock-Provider-Signature header."),
		})
		return
	}
	if !VerifyHMACSignature(h.hmacSecret, rawBody, sigHeader) {
		// NEVER log sigHeader or rawBody — security redaction.
		slog.WarnContext(c.Request.Context(), "callback signature mismatch",
			"action", "SIGNATURE_INVALID",
		)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       codeOf("AUTH_INVALID"),
			Message:    msgOf("Signature verification failed."),
		})
		return
	}

	// ── 4. Parse body (raw bytes already verified) ────────────────────────────
	var req CallbackRequest
	if err := bindRawJSON(rawBody, &req); err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       wrapper.CodeBadRequest,
			Message:    msgOf("Invalid request body."),
		})
		return
	}

	// Validate required fields manually (binding tags not applied by bindRawJSON).
	if req.PaymentIntentID == "" || req.ProviderStatus == "" || req.MockPaymentRef == "" || req.Amount <= 0 {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       wrapper.CodeBadRequest,
			Message:    msgOf("Missing required fields in callback body."),
		})
		return
	}

	intentID, err := uuid.Parse(req.PaymentIntentID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       wrapper.CodeBadRequest,
			Message:    msgOf("Invalid paymentIntentId."),
		})
		return
	}

	// Validate providerStatus enum.
	switch req.ProviderStatus {
	case "SUCCEEDED", "FAILED", "EXPIRED":
		// valid
	default:
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       wrapper.CodeBadRequest,
			Message:    msgOf("Invalid providerStatus; must be SUCCEEDED, FAILED, or EXPIRED."),
		})
		return
	}

	// ── 5. Invoke canonical process_callback ─────────────────────────────────
	in := ProcessCallbackInput{
		IntentID:          intentID,
		ProviderStatus:    req.ProviderStatus,
		MockPaymentRef:    req.MockPaymentRef,
		AmountFromCaller:  req.Amount,
		ProviderTimestamp: providerTimestamp,
	}

	out, err := h.svc.Callback(c.Request.Context(), in)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Callback process error", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       wrapper.CodeInternalError,
			Message:    wrapper.MessageInternalError,
		})
		return
	}

	// Return process_callback output verbatim (200 or 409 with stored envelope).
	c.Data(out.HTTPStatus, "application/json", out.Body)
}
