package payment

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// IntentCreate handles POST /api/v1/payment/intent/create.
//
// Auth: X-Internal-Secret header (checked by internalSecretMW in router).
// Caller: checkout service only.
// Idempotent on order_id: returns existing intent if REQUIRES_PAYMENT.
// Returns 409 CONFLICT if the intent is already in a terminal state.
func (h *handler) IntentCreate(c *gin.Context) {
	req, ok := wrapper.BindJSON[IntentCreateRequest](c, slog.String("handler", "IntentCreate"))
	if !ok {
		return
	}

	resp, err := h.svc.CreateIntent(c.Request.Context(), *req)
	if err != nil {
		switch {
		case errors.Is(err, ErrConflict):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusConflict,
				Code:       codeOf("CONFLICT"),
				Message:    msgOf("A terminal payment intent already exists for this order."),
			})
		default:
			slog.ErrorContext(c.Request.Context(), "CreateIntent error", "error", err)
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusInternalServerError,
				Code:       wrapper.CodeInternalError,
				Message:    wrapper.MessageInternalError,
			})
		}
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[IntentCreateResponse]{
		HTTPStatus: http.StatusOK,
		Code:       wrapper.CodeSuccess,
		Message:    wrapper.MessageSuccess,
		Data:       &resp,
	})
}
