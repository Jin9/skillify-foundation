package auth

import (
	"log/slog"
	"net/http"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/auth/access"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IssueTokenRequest struct {
	MemberID       uuid.UUID                         `json:"memberId" binding:"required"`
	OrganizationID uuid.UUID                         `json:"organizationId" binding:"required"`
	MemberRole     access.OrganizationMemberRoleType `json:"memberRole" binding:"required"`
}

type IssueTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func (h *handler) IssueToken(c *gin.Context) {
	req, ok := wrapper.BindJSON[IssueTokenRequest](c, slog.String("handler", "IssueToken"))
	if !ok {
		return
	}

	claims := token.Claims{
		Sub: req.MemberID.String(),
		Jti: uuid.New().String(),
		Extra: map[string]any{
			"organizationID": req.OrganizationID.String(),
			"memberRole":     string(req.MemberRole),
		},
	}

	accessToken, err := h.token.SignES256(claims)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[IssueTokenResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err: serror.Wrap(err).With(
				slog.String("member_id", req.MemberID.String()),
				slog.String("organization_id", req.OrganizationID.String()),
			),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[IssueTokenResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &IssueTokenResponse{
			AccessToken: accessToken,
		},
	})
}
