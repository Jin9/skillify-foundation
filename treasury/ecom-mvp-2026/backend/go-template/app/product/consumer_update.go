package product

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
)

type UpdateProductMessage struct {
	ProductID   string  `json:"productId" binding:"required"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
}

func (h *handler) OnProductUpdated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {

	var payload UpdateProductMessage
	if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
		return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
	}

	product := access.Product{
		ProductID:   payload.ProductID,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
	}

	_, err := h.productStorage.UpdateProduct(ctx, product)
	if err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("product_id", payload.ProductID),
		)
	}
	return nil
}
