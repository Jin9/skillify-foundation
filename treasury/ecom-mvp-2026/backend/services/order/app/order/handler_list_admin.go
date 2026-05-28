package order

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// HandleListAdmin implements POST /api/v1/order/order/list-admin.
// STUB — returns 501 Not Implemented.
//
// Full implementation requires:
//   - ADMIN role check (403 AUTH_FORBIDDEN for non-ADMIN)
//   - Paginated query with optional filters: status, dateFrom, dateTo, search (order_number prefix OR buyer_email_snapshot exact)
//   - Write admin_action_log row in-tx (audit fallback per td.json)
func (s *Service) HandleListAdmin(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       app.CodeInternalError,
		Message:    "not implemented",
	})
}
