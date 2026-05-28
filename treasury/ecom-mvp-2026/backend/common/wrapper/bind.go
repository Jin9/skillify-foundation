package wrapper

import (
	"log/slog"
	"net/http"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"

	"github.com/gin-gonic/gin"
)

// BindJSON attempts to decode the request body into *T.
// On success it returns the populated struct and true.
// On failure it writes a 400 Bad Request response and returns nil, false,
// so the caller can simply return early.
// Optional slog.Attr values are attached to the error for structured logging.
func BindJSON[T any](c *gin.Context, attrs ...slog.Attr) (*T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		Respond(c, ResponseOption[T]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeBadRequest,
			Message:    MessageBadRequest,
			Err:        serror.Wrap(err).With(attrs...),
		})
		return nil, false
	}
	return &req, true
}
