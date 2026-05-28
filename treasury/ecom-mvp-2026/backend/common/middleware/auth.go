package middleware

import (
	"net/http"
	"strings"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

	"github.com/gin-gonic/gin"
)

const (
	bearerScheme = "Bearer "
)

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func JWT(jwtManager token.JWTParser, jwtVerifier token.JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := new(authHeader)
		if err := c.ShouldBindHeader(header); err != nil {
			abortUnauthorized(c)
			return
		}

		rawToken := extractBearerToken(header.Authorization)

		err := jwtVerifier.VerifySignatureES256(rawToken)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		claim, err := jwtManager.ParseToken(rawToken)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		if err := jwtManager.ValidateClaims(claim); err != nil {
			abortUnauthorized(c)
			return
		}

		// REV-L2-102: reject tokens that are not access tokens.
		// Refresh tokens (tokenType:"refresh") must not authenticate Bearer routes.
		// Default-deny: absent or non-"access" tokenType is rejected.
		if tokenType, _ := claim.Extra["tokenType"].(string); tokenType != "access" {
			abortUnauthorized(c)
			return
		}

		ctx := token.WithClaims(c.Request.Context(), claim)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusUnauthorized,
		Code:       app.CodeUnauthorized,
		Message:    app.MessageUnauthorized,
	})

	c.Abort()
}

func extractBearerToken(authorization string) string {
	trimmed := strings.TrimPrefix(authorization, bearerScheme)
	return strings.TrimSpace(trimmed)
}
