package identity

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,max=128"`
}

type loginResponse struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiresAt  string `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt string `json:"refreshTokenExpiresAt"`
	Role                  string `json:"role"`
}

func (h *handler) Login(c *gin.Context) {
	req, ok := wrapper.BindJSON[loginRequest](c, slog.String("handler", "Login"))
	if !ok {
		return
	}

	accessToken, refreshToken, accessExp, refreshExp, role, err := h.svc.Login(
		c.Request.Context(), req.Email, req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrAuthInvalid):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnauthorized,
				Code:       "AUTH_INVALID",
				Message:    "Invalid email or password.",
			})
		case errors.Is(err, ErrAuthSuspended):
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusForbidden,
				Code:       "AUTH_SUSPENDED",
				Message:    "Account suspended.",
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

	wrapper.Respond(c, wrapper.ResponseOption[loginResponse]{
		HTTPStatus: http.StatusOK,
		Code:       wrapper.CodeSuccess,
		Message:    wrapper.MessageSuccess,
		Data: &loginResponse{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  accessExp.Format("2006-01-02T15:04:05Z07:00"),
			RefreshTokenExpiresAt: refreshExp.Format("2006-01-02T15:04:05Z07:00"),
			Role:                  role,
		},
	})
}
