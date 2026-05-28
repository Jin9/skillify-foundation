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

// addItemRequest is the POST body for cart.add-item.
type addItemRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Qty       int    `json:"qty"       binding:"required,min=1,max=100"`
}

// AddItem handles POST /api/v1/cart/cart/add-item.
// Auth: customer JWT — user_id extracted from claims.sub.
// Logic: validate → catalog check → upsert with ON CONFLICT merge → return enriched cart.
func (h *handler) AddItem(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Extract claims from JWT middleware.
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

	// 2. Bind and validate request body.
	req, ok := wrapper.BindJSON[addItemRequest](c, slog.String("handler", "AddItem"))
	if !ok {
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "productId must be a valid UUID",
		})
		return
	}

	// 3. Delegate to service (catalog check + upsert + enrich).
	view, err := h.service.AddItem(ctx, userID, productID, req.Qty)
	if err != nil {
		switch {
		case errors.Is(err, ErrProductNotFound):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       app.CodeNotFound,
				Message:    "product not found",
				Err:        err,
			})
		case errors.Is(err, ErrProductDeleted):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnprocessableEntity,
				Code:       "PRODUCT_DELETED",
				Message:    "product has been deleted",
				Err:        err,
			})
		case errors.Is(err, ErrProductInactive):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnprocessableEntity,
				Code:       "PRODUCT_INACTIVE",
				Message:    "product is not active",
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

	// 4. Return enriched cart view.
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
