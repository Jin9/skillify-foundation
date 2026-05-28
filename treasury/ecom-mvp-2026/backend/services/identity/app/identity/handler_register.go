package identity

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type registerRequest struct {
	Email    string  `json:"email"    binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8,max=128"`
	Name     string  `json:"name"     binding:"required,min=1,max=80"`
	Phone    *string `json:"phone"`
}

type registerResponse struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func (h *handler) Register(c *gin.Context) {
	req, ok := wrapper.BindJSON[registerRequest](c, slog.String("handler", "Register"))
	if !ok {
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Email, req.Password, req.Name, req.Phone)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyRegistered) {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusConflict,
				Code:       "EMAIL_ALREADY_REGISTERED",
				Message:    "Email is already registered.",
			})
			return
		}
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       wrapper.CodeInternalError,
			Message:    wrapper.MessageInternalError,
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[registerResponse]{
		HTTPStatus: http.StatusCreated,
		Code:       wrapper.CodeSuccess,
		Message:    wrapper.MessageSuccess,
		Data: &registerResponse{
			UserID: user.ID.String(),
			Email:  user.Email,
			Name:   user.Name,
		},
	})
}
