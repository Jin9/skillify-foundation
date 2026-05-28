package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"

	cartdomain "gitlab.com/b2c-e-commerce-platform/platform/backend/cart/app/cart"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/cart/app/cart/access"
)

// New constructs a gin.Engine with all cart routes and middleware configured.
func New(d Deps, version, commit string, timeoutDuration time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if commonconfig.IsLocalEnv() {
		r.Use(gin.Logger())
	}

	r.GET("/liveness", health.Liveness(version, commit))
	r.GET("/metrics", health.Metrics())
	r.GET("/readiness", health.Readiness())

	r.Use(
		middleware.SecurityHeaders(),
		middleware.AccessControl(d.cfg.AccessControl.AllowOrigin, allowedHeaders(d.cfg.Header.RefIDHeaderKey)),
		middleware.TraceContextTraceIDMiddleware(""),
		middleware.RefIDMiddleware(d.cfg.Header.RefIDHeaderKey),
		middleware.AutoLoggingMiddleware(app.CodeSuccess),
		middleware.Timeout(timeoutDuration),
		middleware.AccessLog(),
	)

	registerCartRoutes(r, d)

	return r
}

func registerCartRoutes(r *gin.Engine, d Deps) {
	storage := access.NewCartStorage(d.pool)

	catalogClient := cartdomain.NewCatalogClient(d.httpClient, d.cfg.Catalog.BaseURL)
	inventoryClient := cartdomain.NewInventoryClient(d.httpClient, d.cfg.Inventory.BaseURL)

	fanoutCfg := cartdomain.FanoutConfig{
		TimeoutMS:   d.cfg.Fanout.TimeoutMS,
		MaxItems:    d.cfg.Fanout.MaxItems,
		Concurrency: d.cfg.Fanout.Concurrency,
	}

	svc := cartdomain.NewService(storage, catalogClient, inventoryClient, fanoutCfg)

	h := cartdomain.NewHandler(cartdomain.HandlerConfig{Service: svc})

	// Protected group — all cart endpoints require a valid customer JWT.
	protected := r.Group("/api/v1/cart/cart")
	protected.Use(middleware.JWT(d.jwtParser, d.jwtVerifier))
	{
		protected.POST("/add-item",          h.AddItem)
		protected.POST("/read",              h.ReadCart)
		protected.POST("/update-item",       h.UpdateItem)
		protected.POST("/remove-item",       h.RemoveItem)          // STUB 501
		protected.POST("/clear-on-checkout", h.ClearOnCheckout)     // STUB 501
	}
}

func allowedHeaders(refIDHeaderKey string) []string {
	return []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
		refIDHeaderKey,
	}
}
