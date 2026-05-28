package router

import (
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/app/payment"
)

// New builds and returns the Gin engine for the payment service.
func New(d Deps, timeoutDuration time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if commonconfig.IsLocalEnv() {
		r.Use(gin.Logger())
	}

	r.GET("/liveness", health.Liveness("payment", ""))
	r.GET("/readiness", health.Readiness())

	r.Use(
		middleware.SecurityHeaders(),
		middleware.AccessLog(),
		middleware.Timeout(timeoutDuration),
	)

	registerPaymentRoutes(r, d)

	return r
}

func registerPaymentRoutes(r *gin.Engine, d Deps) {
	h := payment.NewHandler(payment.HandlerConfig{
		Service:    d.svc,
		HMACSecret: []byte(d.cfg.Payment.CallbackHMACSecret),
	})

	jwtMW := middleware.JWT(d.parser, d.verifier)
	internalMW := internalSecretMW(d.cfg.Payment.InternalSecret)

	api := r.Group("/api/v1/payment/intent")
	{
		// Internal: called by checkout service (X-Internal-Secret header auth).
		api.POST("/create", internalMW, h.IntentCreate)

		// Customer-facing: JWT required, role=CUSTOMER.
		api.POST("/simulate", jwtMW, h.Simulate)

		// Webhook: HMAC signature verified inside the handler (raw body needed).
		api.POST("/callback", h.Callback)
	}
}

// internalSecretMW is a minimal middleware that checks the X-Internal-Secret header.
// MVP internal auth for service-to-service calls (checkout → payment.intent.create).
// Uses constant-time comparison to prevent timing-oracle attacks (REV-L2-003).
func internalSecretMW(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		supplied := c.GetHeader("X-Internal-Secret")
		// Always run ConstantTimeCompare even on empty header so response time
		// does not leak whether the header was present.
		if subtle.ConstantTimeCompare([]byte(supplied), []byte(secret)) != 1 {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusUnauthorized,
				Code:       "AUTH_INVALID",
				Message:    "Invalid or missing internal secret.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
