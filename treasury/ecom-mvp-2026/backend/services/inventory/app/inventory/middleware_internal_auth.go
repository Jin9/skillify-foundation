package inventory

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

const (
	// HeaderInternalSecret is the request header checked by InternalAuthMiddleware.
	// Checkout service must supply the value of $INTERNAL_SHARED_SECRET on every
	// call to reservation.create and reservation.release.
	HeaderInternalSecret = "X-Internal-Secret"

	// CodeAuthInvalid is the application-level code returned on header mismatch.
	CodeAuthInvalid    app.Code    = "INV401"
	MessageAuthInvalid app.Message = "Internal auth invalid"
)

// InternalAuthMiddleware gates endpoints that are only callable by internal
// services (e.g. Checkout → reservation.create / reservation.release).
//
// MVP implementation: shared-secret header check. The header value must equal
// the INTERNAL_SHARED_SECRET environment variable supplied via config.Internal.SharedSecret.
//
// Hardening path (deferred): replace with mTLS or SERVICE-role JWT as noted in
// td.json and cross-cutting.route-convention.trailing_rules.
func InternalAuthMiddleware(sharedSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		supplied := c.GetHeader(HeaderInternalSecret)
		// Use constant-time comparison to prevent timing-oracle attacks (REV-L2-003).
		// An empty supplied value must still go through the comparison so the
		// response time does not leak whether the header was present.
		if subtle.ConstantTimeCompare([]byte(supplied), []byte(sharedSecret)) != 1 {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnauthorized,
				Code:       CodeAuthInvalid,
				Message:    MessageAuthInvalid,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
