package identity

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type addressSetDefaultRequest struct {
	AddressID string `json:"addressId" binding:"required,uuid"`
}

type addressSetDefaultResponse struct {
	AddressID string `json:"addressId"`
	IsDefault bool   `json:"isDefault"`
}

// AddressSetDefault atomically promotes one live address to the user's default
// delivery address (spec: identity-address/set-default.md).
//
// Steps:
//  1. Validate Bearer JWT → extract user_id from context claims.
//  2. Bind and validate request (addressId must be a valid UUID).
//  3. Delegate the two-step UPDATE transaction to service.SetAddressDefault.
//  4. Map ErrAddressNotFound → 404 NOT_FOUND; other errors → 500.
func (h *handler) AddressSetDefault(c *gin.Context) {
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
		slog.WarnContext(c.Request.Context(), "address.set_default invalid sub",
			slog.String("sub", claims.Sub),
		)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       "AUTH_INVALID",
			Message:    "Authentication required.",
		})
		return
	}

	req, ok := wrapper.BindJSON[addressSetDefaultRequest](c, slog.String("handler", "AddressSetDefault"))
	if !ok {
		return
	}

	// uuid binding tag already enforces format; Parse is safe here.
	addressID, _ := uuid.Parse(req.AddressID)

	if err := h.svc.SetAddressDefault(c.Request.Context(), userID, addressID); err != nil {
		if errors.Is(err, ErrAddressNotFound) {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       "NOT_FOUND",
				Message:    "Address not found.",
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

	wrapper.Respond(c, wrapper.ResponseOption[addressSetDefaultResponse]{
		HTTPStatus: http.StatusOK,
		Code:       wrapper.CodeSuccess,
		Message:    wrapper.MessageSuccess,
		Data: &addressSetDefaultResponse{
			AddressID: addressID.String(),
			IsDefault: true,
		},
	})
}
