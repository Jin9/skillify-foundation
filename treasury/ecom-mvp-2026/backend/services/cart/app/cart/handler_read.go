package cart

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// ReadCart handles POST /api/v1/cart/cart/read.
// Auth: customer JWT — user_id extracted from claims.sub.
// Logic:
//   - Fetch cart + items from DB.
//   - Fan-out to catalog (per-item, parallel, semaphore) and inventory (bulk) using errgroup.
//   - Cap at FANOUT_MAX_ITEMS (50) to bound the fan-out.
//   - Per-item: catalog timeout → line carries available:false, reason:"catalog_unreachable" (read continues).
//   - Compute subtotal = sum(currentPrice * qty) for checkoutable lines only.
//   - Return {cartId, items, subtotal, currency:"THB", updatedAt}.
func (h *handler) ReadCart(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Extract claims.
	claims, err := token.ClaimsFromContext(ctx)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       app.CodeUnauthorized,
			Message:    app.MessageUnauthorized,
		})
		return
	}

	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       app.CodeUnauthorized,
			Message:    app.MessageUnauthorized,
		})
		return
	}

	// 2. Read cart with parallel fan-out (errgroup + semaphore inside service).
	view, err := h.service.ReadCart(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "cart.read failed", slog.Any("error", err))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        err,
		})
		return
	}

	type response struct {
		Cart CartView `json:"cart"`
	}
	wrapper.Respond(c, wrapper.ResponseOption[response]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data:       &response{Cart: view},
	})
}
