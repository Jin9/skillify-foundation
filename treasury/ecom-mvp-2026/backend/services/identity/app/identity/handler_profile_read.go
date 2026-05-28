package identity

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type profileReadResponse struct {
	UserID string  `json:"userId"`
	Email  string  `json:"email"`
	Name   string  `json:"name"`
	Phone  *string `json:"phone,omitempty"`
	Role   string  `json:"role"`
}

// ProfileRead returns the authenticated customer's own profile (CUST-001).
//
// Auth: customer JWT required. The middleware places the parsed claims in the
// request context via token.WithClaims; we extract claims.Sub as the user id.
//
// Implements requirement CUST-001 and unblocks Checkout's `buyerEmail`
// hand-off (REV-L2-001 from dry-run #1).
func (h *handler) ProfileRead(c *gin.Context) {
	claims, err := token.ClaimsFromContext(c.Request.Context())
	if err != nil || claims == nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       "AUTH_INVALID",
			Message:    "Authentication required.",
		})
		return
	}

	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		slog.WarnContext(c.Request.Context(), "profile.read invalid sub",
			slog.String("sub", claims.Sub),
		)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       "AUTH_INVALID",
			Message:    "Authentication required.",
		})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       "USER_NOT_FOUND",
				Message:    "User not found.",
			})
			return
		}
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       wrapper.CodeInternalError,
			Message:    wrapper.MessageInternalError,
			Err:        serror.Wrap(err),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[profileReadResponse]{
		HTTPStatus: http.StatusOK,
		Code:       wrapper.CodeSuccess,
		Message:    wrapper.MessageSuccess,
		Data: &profileReadResponse{
			UserID: user.ID.String(),
			Email:  user.Email,
			Name:   user.Name,
			Phone:  user.Phone,
			Role:   user.Role,
		},
	})
}
