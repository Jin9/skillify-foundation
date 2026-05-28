package inventory

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// StockAdjustHandler is a STUB handler for POST /api/v1/inventory/stock/adjust.
// Full implementation:
//   - Auth: ADMIN only (reject CUSTOMER with AUTH_FORBIDDEN).
//   - BEGIN TX; SELECT FOR UPDATE on stock_levels; check available_qty + delta >= 0
//     else STOCK_NEGATIVE_INVARIANT; UPDATE stock_levels; INSERT stock_adjustments; COMMIT.
//   - Optional Idempotency-Key header per cross-cutting.idempotency.
//
// TODO: implement per td.json endpoint inventory.stock.adjust (INV-007, INV-008).
func (svc *Service) StockAdjustHandler(c *gin.Context) {
	traceID := middleware.RefID(c)
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       CodeNotImplementedMVP,
		Message:    MessageNotImplementedMVP,
		TraceID:    traceID,
	})
}
