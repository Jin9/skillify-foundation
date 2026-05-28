package catalog

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// catalogService is the minimal interface the Handler needs from the service layer.
// CatalogService satisfies this interface automatically.
type catalogService interface {
	ListProducts(ctx context.Context, filter ProductListFilter, sort SortOption, page, limit int) (ListProductsResult, error)
	GetProduct(ctx context.Context, productID uuid.UUID, isAdmin bool) (ProductDetail, error)
	CreateProduct(ctx context.Context, input CreateProductInput) (uuid.UUID, error)
}

// Handler groups all catalog HTTP handlers.
type Handler struct {
	svc catalogService
}

// NewHandler constructs a Handler.
func NewHandler(svc *CatalogService) *Handler {
	return &Handler{svc: svc}
}

// stubNotImplemented writes a uniform 501 response for unimplemented endpoints.
func stubNotImplemented(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       app.Code("NOT_IMPLEMENTED_MVP"),
		Message:    "Endpoint stubbed for MVP — see README",
	})
}
