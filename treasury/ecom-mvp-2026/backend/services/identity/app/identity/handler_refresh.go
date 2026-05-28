package identity

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type refreshResponse struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiresAt  string `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt string `json:"refreshTokenExpiresAt"`
}

func (h *handler) Refresh(c *gin.Context) {
	req, ok := wrapper.BindJSON[refreshRequest](c, slog.String("handler", "Refresh"))
	if !ok {
		return
	}

	accessToken, refreshToken, accessExp, refreshExp, err := h.svc.Refresh(
		c.Request.Context(), req.RefreshToken,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrAuthInvalid):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnauthorized,
				Code:       "AUTH_INVALID",
				Message:    "Token invalid or expired.",
			})
		case errors.Is(err, ErrAuthRevoked):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnauthorized,
				Code:       "AUTH_REVOKED",
				Message:    "Refresh token has been revoked.",
			})
		default:
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusInternalServerError,
				Code:       wrapper.CodeInternalError,
				Message:    wrapper.MessageInternalError,
			})
		}
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[refreshResponse]{
		HTTPStatus: http.StatusOK,
		Code:       wrapper.CodeSuccess,
		Message:    wrapper.MessageSuccess,
		Data: &refreshResponse{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  accessExp.Format("2006-01-02T15:04:05Z07:00"),
			RefreshTokenExpiresAt: refreshExp.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}
