package organization

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access"
)

type AddOrganizationMemberMessage struct {
	OrganizationID string                            `json:"organizationId" binding:"required"`
	MemberID       string                            `json:"memberId" binding:"required"`
	Role           access.OrganizationMemberRoleType `json:"role" binding:"required"`
}

func (h *handler) OnOrganizationMemberAdded(ctx context.Context, msg kafka.Message[json.RawMessage]) error {

	var payload AddOrganizationMemberMessage
	if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
		return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
	}

	orgMember := access.OrganizationMember{
		OrganizationID: payload.OrganizationID,
		MemberID:       payload.MemberID,
		Role:           payload.Role,
		Status:         access.OrganizationMemberStatusActive,
	}

	if err := h.organizationStorage.CreateOrganizationMember(ctx, orgMember); err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("organization_id", payload.OrganizationID),
			slog.String("member_id", payload.MemberID),
		)
	}
	return nil
}
