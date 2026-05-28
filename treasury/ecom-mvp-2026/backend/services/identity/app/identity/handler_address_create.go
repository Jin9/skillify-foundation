package identity

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

func (h *handler) AddressCreate(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       "NOT_IMPLEMENTED_MVP",
		Message:    "Endpoint stubbed for MVP — see README",
	})
}
