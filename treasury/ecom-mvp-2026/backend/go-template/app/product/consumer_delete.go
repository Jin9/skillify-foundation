package product

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
)

type DeleteProductMessage struct {
	ProductID string `json:"productId" binding:"required"`
}

func (h *handler) OnProductDeleted(ctx context.Context, msg kafka.Message[json.RawMessage]) error {

	var payload DeleteProductMessage
	if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
		return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
	}

	id, err := (&access.Product{ProductID: payload.ProductID}).GetID()
	if err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("product_id", payload.ProductID),
		)
	}

	if err := h.productStorage.DeleteProduct(ctx, id); err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("product_id", payload.ProductID),
		)
	}
	return nil
}
