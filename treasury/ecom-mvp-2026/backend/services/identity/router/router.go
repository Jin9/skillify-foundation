package router

import (
	"time"

	"github.com/gin-gonic/gin"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity/app/identity"
	identityaccess "gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity/app/identity/access"
)

func New(d Deps, timeoutDuration time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if commonconfig.IsLocalEnv() {
		r.Use(gin.Logger())
	}

	r.GET("/liveness", health.Liveness("identity", ""))
	r.GET("/readiness", health.Readiness())

	r.Use(
		middleware.SecurityHeaders(),
		middleware.AccessLog(),
		middleware.Timeout(timeoutDuration),
	)

	registerIdentityRoutes(r, d)

	return r
}

func registerIdentityRoutes(r *gin.Engine, d Deps) {
	userStorage := identityaccess.NewUserStorage(d.db)
	rtStorage := identityaccess.NewRefreshTokenStorage(d.db)
	addrStorage := identityaccess.NewAddressStorage(d.db)

	svc := identity.NewService(identity.ServiceConfig{
		UserStorage:         userStorage,
		RefreshTokenStorage: rtStorage,
		AddressStorage:      addrStorage,
		AccessSigner:        d.accessSigner,
		RefreshSigner:       d.refreshSigner,
		Verifier:            d.verifier,
		Parser:              d.parser,
		AccessTokenTTL:      d.cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL:     d.cfg.JWT.RefreshTokenTTL,
	})

	h := identity.NewHandler(identity.HandlerConfig{
		Service: svc,
	})

	jwtMW := middleware.JWT(d.parser, d.verifier)

	// Public auth routes
	auth := r.Group("/api/v1/identity/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", jwtMW, h.Logout)
	}

	// Protected profile routes
	profile := r.Group("/api/v1/identity/profile", jwtMW)
	{
		profile.GET("/me", h.ProfileRead)
		profile.PATCH("/me", h.ProfileUpdate)
	}

	// Protected address routes
	addr := r.Group("/api/v1/identity/address", jwtMW)
	{
		addr.POST("/create", h.AddressCreate)
		addr.GET("/list", h.AddressList)
		addr.PATCH("/update", h.AddressUpdate)
		addr.DELETE("/delete", h.AddressDelete)
		addr.POST("/set-default", h.AddressSetDefault)
	}

}
