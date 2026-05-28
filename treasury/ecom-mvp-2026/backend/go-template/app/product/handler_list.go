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

type ListProductsRequest struct {
	OrganizationID uuid.UUID `json:"organizationId" binding:"required"`
}

type ListProductItem struct {
	ProductID      uuid.UUID                `json:"productId"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Price          float64                  `json:"price"`
	OrganizationID string                   `json:"organizationId"`
	Status         access.ProductStatusType `json:"status"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
}

type ListProductsResponse struct {
	Products []ListProductItem `json:"products"`
}

func (h *handler) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[ListProductsRequest](c, slog.String("handler", "ListProducts"))
	if !ok {
		return
	}

	products, err := h.productStorage.ListProducts(ctx, req.OrganizationID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[ListProductsResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("organization_id", req.OrganizationID.String())),
		})
		return
	}

	items := make([]ListProductItem, 0, len(products))
	for _, p := range products {
		pid, err := p.GetID()
		if err != nil {
			continue
		}
		items = append(items, ListProductItem{
			ProductID:      pid,
			Name:           p.Name,
			Description:    p.Description,
			Price:          p.Price,
			OrganizationID: p.OrganizationID,
			Status:         p.Status,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}

	wrapper.Respond(c, wrapper.ResponseOption[ListProductsResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data:       &ListProductsResponse{Products: items},
	})
}
