package product

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
)

type CreateProductRequest struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Price          float64 `json:"price" binding:"required,gt=0"`
	OrganizationID string  `json:"organizationId" binding:"required"`
}

type CreateProductResponse struct {
	ProductID uuid.UUID `json:"productId"`
}

func (h *handler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[CreateProductRequest](c, slog.String("handler", "CreateProduct"))
	if !ok {
		return
	}

	product := access.Product{
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
		OrganizationID: req.OrganizationID,
		Status:         access.ProductStatusActive,
	}

	created, err := h.productStorage.CreateProduct(ctx, product)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err: serror.Wrap(err).With(
				slog.String("product_name", req.Name),
				slog.String("organization_id", req.OrganizationID),
			),
		})
		return
	}

	productID, err := created.GetID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err: serror.Wrap(err).With(
				slog.String("product_name", req.Name),
				slog.String("organization_id", req.OrganizationID),
			),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
		HTTPStatus: http.StatusCreated,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data:       &CreateProductResponse{ProductID: productID},
	})
}
