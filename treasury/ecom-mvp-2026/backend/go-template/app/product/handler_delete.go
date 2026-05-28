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

type DeleteProductRequest struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
}

type DeleteProductResponse struct{}

func (h *handler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[DeleteProductRequest](c, slog.String("handler", "DeleteProduct"))
	if !ok {
		return
	}

	if err := h.productStorage.DeleteProduct(ctx, req.ProductID); err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[DeleteProductResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", req.ProductID.String())),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[DeleteProductResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
	})
}
