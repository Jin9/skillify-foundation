package router

import (
	"time"

	"github.com/gin-gonic/gin"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

	"github.com/example/shoppilot/catalog/app/catalog"
	"github.com/example/shoppilot/catalog/app/catalog/access"
	"github.com/example/shoppilot/catalog/config"
)

// New constructs a gin.Engine with all catalog routes registered.
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
		middleware.AccessControl("*", allowedHeaders(d.Cfg.Header.RefIDHeaderKey)),
		middleware.TraceContextTraceIDMiddleware(""),
		middleware.RefIDMiddleware(d.Cfg.Header.RefIDHeaderKey),
		middleware.Timeout(timeoutDuration),
		middleware.AccessLog(),
	)

	// Wire access-layer storage
	productStorage := access.NewProductStorage(d.DB)
	categoryStorage := access.NewCategoryStorage(d.DB)
	imageStorage := access.NewImageStorage(d.DB)
	outboxStorage := access.NewOutboxStorage(d.DB)

	// Wire service
	svc := catalog.NewCatalogService(d.DB, productStorage, categoryStorage, imageStorage, outboxStorage)
	h := catalog.NewHandler(svc)

	// --- Public routes (no JWT required) ---
	publicGroup := r.Group("/api/v1/catalog")
	{
		productGroup := publicGroup.Group("/product")
		{
			productGroup.POST("/list", h.ListProducts)
			productGroup.POST("/detail", h.GetProduct)
		}

		categoryGroup := publicGroup.Group("/category")
		{
			categoryGroup.POST("/list", h.ListCategories)
		}
	}

	// --- Admin routes (JWT required: ADMIN role) ---
	jwtParser, jwtVerifier := buildJWTMiddleware(d.Cfg)
	adminGroup := r.Group("/api/v1/catalog")
	adminGroup.Use(middleware.JWT(jwtParser, jwtVerifier))
	{
		adminProductGroup := adminGroup.Group("/product")
		{
			adminProductGroup.POST("/create", h.CreateProduct)
			adminProductGroup.POST("/update", h.UpdateProduct)
			adminProductGroup.POST("/soft-delete", h.SoftDeleteProduct)
		}

		adminCategoryGroup := adminGroup.Group("/category")
		{
			adminCategoryGroup.POST("/create", h.CreateCategory)
			adminCategoryGroup.POST("/update", h.UpdateCategory)
			adminCategoryGroup.POST("/deactivate", h.DeactivateCategory)
		}
	}

	return r
}

func buildJWTMiddleware(cfg config.Config) (token.JWTParser, token.JWTVerifier) {
	parser, err := token.NewJWTParser(token.JWTParserConfig{
		Issuer:   cfg.JWT.Issuer,
		Audience: cfg.JWT.Audience,
	})
	if err != nil {
		panic("catalog: failed to build JWT parser: " + err.Error())
	}

	verifier, err := token.NewJWTVerifier(token.JWTVerifierConfig{
		PublicKey: cfg.JWT.PublicKey,
		Alg:       string(token.ES256),
	})
	if err != nil {
		panic("catalog: failed to build JWT verifier: " + err.Error())
	}

	return parser, verifier
}

func allowedHeaders(refIDHeaderKey string) []string {
	return []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		refIDHeaderKey,
	}
}
