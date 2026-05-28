package inventory

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// StockReadHandler is the FULL handler for POST /api/v1/inventory/stock/read.
// Auth: JWT (CUSTOMER or ADMIN). Read-only, no tx required.
func (svc *Service) StockReadHandler(c *gin.Context) {
	traceID := middleware.RefID(c)
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[StockReadRequest](c)
	if !ok {
		return
	}

	if req.SKU == "" {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeValidationError,
			Message:    "sku must be non-empty",
			TraceID:    traceID,
		})
		return
	}

	sl, err := svc.StockStorage.GetBySKU(ctx, req.SKU)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       CodeNotFound,
				Message:    "sku not found",
				TraceID:    traceID,
			})
			return
		}
		slog.ErrorContext(ctx, "stock.read database error",
			slog.String("sku", req.SKU),
			slog.String("error", err.Error()),
		)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
			Err:        serror.Wrap(err),
		})
		return
	}

	data := StockReadResponse{
		SKU:          sl.SKU,
		AvailableQty: sl.AvailableQty,
		ReservedQty:  sl.ReservedQty,
		SoldQty:      sl.SoldQty,
	}

	wrapper.Respond(c, wrapper.ResponseOption[StockReadResponse]{
		HTTPStatus: http.StatusOK,
		Code:       CodeSuccess,
		Message:    MessageSuccess,
		Data:       &data,
		TraceID:    traceID,
	})
}
