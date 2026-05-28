package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/auth/access"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResolveIdentityRequest struct {
	GoogleAccessToken string `json:"googleAccessToken" binding:"required"`
}

type ResolveIdentityResponse struct {
	IsMember       bool                              `json:"isMember"`
	MemberID       uuid.UUID                         `json:"memberId,omitempty"`
	Username       string                            `json:"username,omitempty"`
	Email          string                            `json:"email"`
	HashedEmail    string                            `json:"hashedEmail,omitempty"`
	Status         access.MemberStatusType           `json:"status,omitempty"`
	OrganizationID uuid.UUID                         `json:"organizationID,omitempty"`
	Role           access.OrganizationMemberRoleType `json:"role,omitempty"`
	ProfileImage   string                            `json:"profileImage,omitempty"`
}

func (h *handler) ResolveIdentify(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[ResolveIdentityRequest](c, slog.String("handler", "ResolveIdentity"))
	if !ok {
		return
	}

	googleData, err := h.authenticateGoogle(ctx, req.GoogleAccessToken)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			wrapper.Respond(c, wrapper.ResponseOption[ResolveIdentityResponse]{
				HTTPStatus: http.StatusUnauthorized,
				Code:       app.CodeUnauthorized,
				Message:    "Invalid or expired Google access token",
				Err:        serror.Wrap(err).With(slog.String("handler", "ResolveIdentity")),
			})
			return
		}
		wrapper.Respond(c, wrapper.ResponseOption[ResolveIdentityResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("handler", "ResolveIdentity")),
		})
		return
	}

	hashedEmail := h.hash.HashSha256EncodePepper(googleData.email)

	memberInfo, err := h.memberStorage.GetMemberByEmail(ctx, hashedEmail)
	if err != nil {
		if errors.Is(err, access.ErrMemberNotFound) {
			slog.Info("member not found, returning Google profile data", slog.String("email", googleData.email), slog.String("tag", "resolve identify"))
			wrapper.Respond(c, wrapper.ResponseOption[ResolveIdentityResponse]{
				HTTPStatus: http.StatusOK,
				Code:       app.CodeSuccess,
				Message:    app.MessageSuccess,
				Data: &ResolveIdentityResponse{
					IsMember:     false,
					Email:        googleData.email,
					Username:     googleData.username,
					ProfileImage: googleData.profileImage,
				},
			})
			return
		}
		wrapper.Respond(c, wrapper.ResponseOption[ResolveIdentityResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("hashed_email", hashedEmail)),
		})
		return
	}

	resp, err := h.resolveExistingMember(ctx, memberInfo, googleData.profileImage)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[ResolveIdentityResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("handler", "ResolveIdentity")),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[ResolveIdentityResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data:       resp,
	})
}
