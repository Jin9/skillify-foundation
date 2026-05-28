package organization

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access"
)

type CreateOrganizationMessage struct {
	Name string `json:"name" binding:"required"`
}

func (h *handler) OnOrganizationCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {

	var payload CreateOrganizationMessage
	if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
		return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
	}

	org := access.Organization{
		Name:   payload.Name,
		Status: access.OrganizationStatusActive,
	}

	_, err := h.organizationStorage.CreateOrganization(ctx, org)
	if err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("organization_name", payload.Name),
		)
	}
	return nil
}
