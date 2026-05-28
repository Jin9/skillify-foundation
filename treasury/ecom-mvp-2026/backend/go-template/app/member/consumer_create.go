package member

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access"
)

type CreateMemberMessage struct {
	Username       string `json:"username" binding:"required"`
	EncryptedEmail string `json:"encryptedEmail" binding:"required"`
	HashedEmail    string `json:"hashedEmail" binding:"required"`
}

func (h *handler) OnMemberCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {

	var payload CreateMemberMessage
	if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
		return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
	}

	m := access.Member{
		Username:       payload.Username,
		EncryptedEmail: payload.EncryptedEmail,
		HashedEmail:    payload.HashedEmail,
		Status:         access.MemberStatusActive,
	}

	_, err := h.memberStorage.CreateMember(ctx, m)
	if err != nil {
		return serror.Wrap(err).With(
			slog.String("event_id", msg.EventID),
			slog.String("username", payload.Username),
		)
	}
	return nil
}
