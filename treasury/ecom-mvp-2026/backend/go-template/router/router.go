package router

import (
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/auth"
	authaccess "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/auth/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member"
	memberaccess "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization"
	organizationaccess "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product"
	productaccess "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"

	"github.com/gin-gonic/gin"
)

// New constructs a gin.Engine with routes and middleware configured.
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

	registerAuthRoutes(r, d)
	registerOrganizationRoutes(r, d)
	registerMemberRoutes(r, d)
	registerProductRoutes(r, d)

	return r
}

func registerAuthRoutes(r *gin.Engine, d Deps) {
	googleClient := authaccess.NewGoogleClient(
		d.cfg.GoogleClient.VerifyTokenURL,
		d.cfg.GoogleClient.GetUserProfileURL,
		d.cfg.GoogleClient.RevokeTokenURL,
		d.httpClient,
	)
	memberStorage := authaccess.NewMemberStorage(d.firestoreClient.Inner())
	organizationStorage := authaccess.NewOrganizationStorage(d.firestoreClient.Inner())

	authHandlerCfg := auth.HandlerConfig{
		GoogleClient:        googleClient,
		MemberStorage:       memberStorage,
		OrganizationStorage: organizationStorage,
		Hash:                d.hash,
		Cipher:              d.cipher,
		Token:               d.token,
	}
	authHandler := auth.NewHandler(authHandlerCfg)

	authGroup := r.Group("/api/v1/platform/auth")
	{
		authGroup.POST("/resolve-identity", authHandler.ResolveIdentify)
		authGroup.POST("/issue-token", authHandler.IssueToken)
	}
}

func registerOrganizationRoutes(r *gin.Engine, d Deps) {
	organizationStorage := organizationaccess.NewOrganizationStorage(d.firestoreClient.Inner())
	organizationHandlerCfg := organization.HandlerConfig{
		OrganizationStorage: organizationStorage,
	}
	organizationHandler := organization.NewHandler(organizationHandlerCfg)

	organizationGroup := r.Group("/api/v1/platform/organization")
	{
		organizationGroup.POST("/current", organizationHandler.GetCurrent)
		organizationGroup.POST("/register", organizationHandler.RegisterOrganization)
	}
}

func registerMemberRoutes(r *gin.Engine, d Deps) {
	memberStorage := memberaccess.NewMemberStorage(d.firestoreClient.Inner())
	memberHandlerCfg := member.HandlerConfig{
		MemberStorage: memberStorage,
		Cipher:        d.cipher,
		Hash:          d.hash,
	}
	memberHandler := member.NewHandler(memberHandlerCfg)

	memberGroup := r.Group("/api/v1/platform/member")
	{
		memberGroup.POST("/me", memberHandler.GetInfo)
		memberGroup.POST("/register", memberHandler.RegisterMember)
	}
}

func registerProductRoutes(r *gin.Engine, d Deps) {
	productStorage := productaccess.NewProductStorage(d.firestoreClient.Inner())
	productHandlerCfg := product.HandlerConfig{
		ProductStorage: productStorage,
	}
	productHandler := product.NewHandler(productHandlerCfg)

	productGroup := r.Group("/api/v1/platform/product")
	{
		productGroup.POST("/create", productHandler.CreateProduct)
		productGroup.POST("/detail", productHandler.GetProduct)
		productGroup.POST("/list", productHandler.ListProducts)
		productGroup.POST("/update", productHandler.UpdateProduct)
		productGroup.POST("/delete", productHandler.DeleteProduct)
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
