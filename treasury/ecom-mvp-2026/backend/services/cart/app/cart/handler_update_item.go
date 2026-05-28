package cart

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// updateItemRequest is the POST body for cart.update-item.
type updateItemRequest struct {
	CartItemID string `json:"cartItemId" binding:"required"`
	// qty >= 0: 0 triggers remove path (CART-003 / AMB-004). Negative → VALIDATION_ERROR.
	Qty int `json:"qty" binding:"min=0"`
}

// UpdateItem handles POST /api/v1/cart/cart/update-item.
// Auth: customer JWT.
// Logic:
//   - qty == 0 → remove (delegate to remove path, log).
//   - qty > 0  → stock-check via inventory.stock.bulk-read. If qty > availableQty → 409 STOCK_INSUFFICIENT.
//   - UPDATE cart_items SET qty, added_at=NOW() + touch carts.updated_at.
//   - Return enriched cart.
func (h *handler) UpdateItem(c *gin.Context) {
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

	// 2. Bind and validate.
	req, ok := wrapper.BindJSON[updateItemRequest](c, slog.String("handler", "UpdateItem"))
	if !ok {
		return
	}

	cartItemID, err := uuid.Parse(req.CartItemID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "cartItemId must be a valid UUID",
		})
		return
	}

	// Log qty=0 remove intent (CART-003).
	if req.Qty == 0 {
		slog.InfoContext(ctx, "cart.update-item: qty=0, treating as remove",
			slog.String("user_id", userID.String()),
			slog.String("cart_item_id", cartItemID.String()))
	}

	// 3. Delegate to service.
	view, err := h.service.UpdateItem(ctx, userID, cartItemID, req.Qty)
	if err != nil {
		var withDetail *ErrWithDetail
		switch {
		case errors.As(err, &withDetail) && errors.Is(withDetail, ErrStockInsufficient):
			detail, _ := withDetail.Detail.(StockInsufficientDetail)
			type stockErr struct {
				AvailableQty int `json:"availableQty"`
			}
			wrapper.Respond(c, wrapper.ResponseOption[stockErr]{
				HTTPStatus: http.StatusConflict,
				Code:       "STOCK_INSUFFICIENT",
				Message:    "requested qty exceeds available stock",
				Data:       &stockErr{AvailableQty: detail.AvailableQty},
				Err:        err,
			})
		case errors.Is(err, ErrCartItemNotFound):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       app.CodeNotFound,
				Message:    "cart item not found",
				Err:        err,
			})
		case isUpstreamTimeout(err):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusGatewayTimeout,
				Code:       app.CodeTimeout,
				Message:    "upstream service timeout",
				Err:        err,
			})
		default:
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusInternalServerError,
				Code:       app.CodeInternalError,
				Message:    app.MessageInternalError,
				Err:        err,
			})
		}
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
