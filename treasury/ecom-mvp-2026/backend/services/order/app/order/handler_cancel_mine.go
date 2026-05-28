package order

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// HandleCancelMine implements POST /api/v1/order/order/cancel-mine.
// STUB — returns 501 Not Implemented.
//
// Full implementation requires:
//   - JWT auth (CUSTOMER scope)
//   - SELECT ... FOR UPDATE on owned order
//   - PENDING_PAYMENT → CANCELLED transition guard via ValidateTransition
//   - INSERT order_status_history + outbox_events (events.order.cancelled) in same tx
//   - Idempotent re-cancel by self (AMBIG-ORD-3: only if last history actor == claims.sub)
func (s *Service) HandleCancelMine(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       app.CodeInternalError,
		Message:    "not implemented",
	})
}
