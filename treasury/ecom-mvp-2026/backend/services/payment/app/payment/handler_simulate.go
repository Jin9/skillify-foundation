package payment

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// Simulate handles POST /api/v1/payment/intent/simulate.
//
// Auth: JWT required, role=CUSTOMER (enforced by middleware.JWT before reaching handler).
// Ownership: intent.owner_user_id must match the JWT sub claim.
//
// Maps customer outcome (success|failed|timeout) to providerStatus and invokes
// process_callback in-process — the canonical idempotent path.
//
// Timeout outcome: refused (sweeper-driven; simulate does not replicate it).
func (h *handler) Simulate(c *gin.Context) {
	// Extract JWT claims set by middleware.JWT (stored in request context).
	claims, err := token.ClaimsFromContext(c.Request.Context())
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       wrapper.CodeUnauthorized,
			Message:    wrapper.MessageUnauthorized,
		})
		return
	}
	callerSub := claims.Sub
	if callerSub == "" {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       codeOf("AUTH_INVALID"),
			Message:    msgOf("Invalid token claims."),
		})
		return
	}

	req, ok := wrapper.BindJSON[SimulateRequest](c, slog.String("handler", "Simulate"))
	if !ok {
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

	// Ownership check: load intent and compare owner_user_id to JWT sub.
	intent, err := h.svc.GetIntentByID(c.Request.Context(), intentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       codeOf("PAYMENT_INTENT_NOT_FOUND"),
				Message:    msgOf("Payment intent not found."),
			})
			return
		}
		slog.ErrorContext(c.Request.Context(), "GetIntentByID error", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       wrapper.CodeInternalError,
			Message:    wrapper.MessageInternalError,
		})
		return
	}

	if intent.OwnerUserID.String() != callerSub {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusForbidden,
			Code:       codeOf("AUTH_FORBIDDEN"),
			Message:    msgOf("You do not own this payment intent."),
		})
		return
	}

	// Timeout outcome: refuse — sweeper-driven path only.
	if req.Outcome == "timeout" {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnprocessableEntity,
			Code:       wrapper.CodeBadRequest,
			Message:    msgOf("Timeout outcome is driven by the expiry sweeper and cannot be simulated directly."),
		})
		return
	}

	// Invoke the canonical process_callback in-process.
	out, err := h.svc.Simulate(c.Request.Context(), intentID, req.Outcome)
	if err != nil {
		switch {
		case errors.Is(err, ErrIntentNotFound):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       codeOf("PAYMENT_INTENT_NOT_FOUND"),
				Message:    msgOf("Payment intent not found."),
			})
		default:
			slog.ErrorContext(c.Request.Context(), "Simulate error", "error", err)
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusInternalServerError,
				Code:       wrapper.CodeInternalError,
				Message:    wrapper.MessageInternalError,
			})
		}
		return
	}

	// Return process_callback's output verbatim (may be 200 OK or 409 terminal).
	c.Data(out.HTTPStatus, "application/json", out.Body)
}
