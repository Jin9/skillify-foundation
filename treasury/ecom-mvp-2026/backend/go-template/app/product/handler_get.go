package product

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
)

type GetProductRequest struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
}

type GetProductResponse struct {
	ProductID      uuid.UUID                `json:"productId"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Price          float64                  `json:"price"`
	OrganizationID string                   `json:"organizationId"`
	Status         access.ProductStatusType `json:"status"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
}

func (h *handler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[GetProductRequest](c, slog.String("handler", "GetProduct"))
	if !ok {
		return
	}

	product, err := h.productStorage.GetProductByID(ctx, req.ProductID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", req.ProductID.String())),
		})
		return
	}

	productID, err := product.GetID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", req.ProductID.String())),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[GetProductResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &GetProductResponse{
			ProductID:      productID,
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			OrganizationID: product.OrganizationID,
			Status:         product.Status,
			CreatedAt:      product.CreatedAt,
			UpdatedAt:      product.UpdatedAt,
		},
	})
}
