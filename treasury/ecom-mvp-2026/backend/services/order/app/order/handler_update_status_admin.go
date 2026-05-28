package order

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// HandleUpdateStatusAdmin implements POST /api/v1/order/order/update-status-admin.
// STUB — returns 501 Not Implemented.
//
// Full implementation requires:
//   - ADMIN role guard (403 AUTH_FORBIDDEN)
//   - Request validation: trackingNumber required when toStatus=SHIPPED (ORD-006); reason required when toStatus=CANCELLED (ORD-005)
//   - SELECT ... FOR UPDATE; idempotent same-status → 200 no-op
//   - ValidateTransition(current, toStatus, ADMIN)
//   - UPDATE orders + INSERT order_status_history + (if CANCELLED) INSERT outbox_events in same tx
//   - INSERT admin_action_log in same tx
func (s *Service) HandleUpdateStatusAdmin(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       app.CodeInternalError,
		Message:    "not implemented",
	})
}
