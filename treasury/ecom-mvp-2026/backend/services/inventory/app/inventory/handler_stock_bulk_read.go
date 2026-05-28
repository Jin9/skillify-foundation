package inventory

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// StockBulkReadHandler is a STUB handler for POST /api/v1/inventory/stock/bulk-read.
// Full implementation: validate len(skus) in [1,200]; SELECT WHERE sku = ANY($1);
// reorder result to match input order; drop missing; return envelope SUCCESS.
// TODO: implement per td.json endpoint inventory.stock.bulk-read.
func (svc *Service) StockBulkReadHandler(c *gin.Context) {
	traceID := middleware.RefID(c)
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       CodeNotImplementedMVP,
		Message:    MessageNotImplementedMVP,
		TraceID:    traceID,
	})
}
