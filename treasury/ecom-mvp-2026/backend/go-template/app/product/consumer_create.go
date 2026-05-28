package product

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
)

type CreateProductMessage struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Price          float64 `json:"price" binding:"required,gt=0"`
	OrganizationID string  `json:"organizationId" binding:"required"`
}

func (h *handler) OnProductCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {

	var payload CreateProductMessage
	if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
		return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
	}

	product := access.Product{
		Name:           payload.Name,
		Description:    payload.Description,
		Price:          payload.Price,
		OrganizationID: payload.OrganizationID,
		Status:         access.ProductStatusActive,
	}

	_, err := h.productStorage.CreateProduct(ctx, product)
	if err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("product_name", payload.Name),
			slog.String("organization_id", payload.OrganizationID),
		)
	}
	return nil
}
