package product

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type UpdateProductRequest struct {
	ProductID   uuid.UUID `json:"productId" binding:"required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"omitempty,gt=0"`
}

type UpdateProductResponse struct {
	ProductID uuid.UUID `json:"productId"`
}

func (h *handler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[UpdateProductRequest](c, slog.String("handler", "UpdateProduct"))
	if !ok {
		return
	}

	existing, err := h.productStorage.GetProductByID(ctx, req.ProductID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[UpdateProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", req.ProductID.String())),
		})
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}

	updated, err := h.productStorage.UpdateProduct(ctx, existing)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[UpdateProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", req.ProductID.String())),
		})
		return
	}

	productID, err := updated.GetID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[UpdateProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", req.ProductID.String())),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[UpdateProductResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data:       &UpdateProductResponse{ProductID: productID},
	})
}
