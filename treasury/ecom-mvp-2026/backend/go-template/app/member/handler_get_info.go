package member

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access"
)

type GetMemberRequest struct {
	MemberID uuid.UUID `json:"memberId" binding:"required"`
}

type GetMemberResponse struct {
	MemberID    uuid.UUID               `json:"memberId"`
	Username    string                  `json:"username"`
	Email       string                  `json:"email"`
	HashedEmail string                  `json:"hashedEmail"`
	Status      access.MemberStatusType `json:"status"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

func (h *handler) GetInfo(c *gin.Context) {

	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[GetMemberRequest](c, slog.String("handler", "GetMemberInfo"))
	if !ok {
		return
	}

	memberInfo, err := h.memberStorage.GetMemberByID(ctx, req.MemberID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("member_id", req.MemberID.String())),
		})
		return
	}

	email, err := h.cipher.Decrypt(memberInfo.EncryptedEmail)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	memberID, err := memberInfo.GetID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("member_id", req.MemberID.String())),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[GetMemberResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &GetMemberResponse{
			MemberID:    memberID,
			Username:    memberInfo.Username,
			Email:       email,
			HashedEmail: memberInfo.HashedEmail,
			Status:      memberInfo.Status,
			CreatedAt:   memberInfo.CreatedAt,
			UpdatedAt:   memberInfo.UpdatedAt,
		},
	})

}
