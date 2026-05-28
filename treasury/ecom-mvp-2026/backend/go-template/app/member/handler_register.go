package member

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access"
)

type RegisterMemberRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

type RegisterMemberResponse struct {
	MemberID uuid.UUID `json:"memberId"`
}

func (h *handler) RegisterMember(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[RegisterMemberRequest](c, slog.String("handler", "RegisterMember"))
	if !ok {
		return
	}

	// Hash email
	hashedEmail := h.hash.HashSha256EncodePepper(req.Email)

	// Generate IV for encryption
	iv, err := h.cipher.GenerateIV()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("handler", "RegisterMember")),
		})
		return
	}

	// Encrypt email
	encryptedEmail, err := h.cipher.Encrypt(req.Email, iv)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("handler", "RegisterMember")),
		})
		return
	}

	member := access.Member{
		Username:       req.Name,
		EncryptedEmail: encryptedEmail,
		HashedEmail:    hashedEmail,
		Status:         access.MemberStatusActive,
	}

	createdMember, err := h.memberStorage.CreateMember(ctx, member)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("username", req.Name)),
		})
		return
	}

	memberID, err := createdMember.GetID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterMemberResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("username", req.Name)),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[RegisterMemberResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &RegisterMemberResponse{
			MemberID: memberID,
		},
	})
}
