package router

import (
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory"
)

// New constructs a gin.Engine with all inventory routes and middleware.
func New(d Deps, svc *inventory.Service, version, commit string, timeoutDuration time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if commonconfig.IsLocalEnv() {
		r.Use(gin.Logger())
	}

	// Public health endpoints (no auth, no timeout middleware).
	r.GET("/liveness", health.Liveness(version, commit))
	r.GET("/readiness", health.Readiness())
	r.GET("/metrics", health.Metrics())

	r.Use(
		middleware.SecurityHeaders(),
		middleware.AccessControl("*", allowedHeaders(d.cfg.Header.RefIDHeaderKey)),
		middleware.TraceContextTraceIDMiddleware(""),
		middleware.RefIDMiddleware(d.cfg.Header.RefIDHeaderKey),
		middleware.AutoLoggingMiddleware(app.CodeSuccess),
		middleware.Timeout(timeoutDuration),
		middleware.AccessLog(),
	)

	jwtParser, jwtVerifier := d.jwtMW()
	registerInventoryRoutes(r, svc, d.cfg.Internal.SharedSecret, jwtParser, jwtVerifier)

	return r
}

// registerInventoryRoutes registers all inventory HTTP routes with appropriate auth guards.
//
// Auth coverage (INV-L1-04):
//
//	stock.read, stock.bulk-read  — customer JWT required (cart reads need auth).
//	stock.adjust                 — admin JWT required (ADMIN-only operation).
//	reservation.create           — X-Internal-Secret required (Checkout → Inventory only).
//	reservation.release          — X-Internal-Secret required (Checkout → Inventory only).
//
// Routes match cross-cutting.route-convention: POST /api/v1/inventory/{aggregate}/{action}.
func registerInventoryRoutes(
	r *gin.Engine,
	svc *inventory.Service,
	internalSecret string,
	jwtParser token.JWTParser,
	jwtVerifier token.JWTVerifier,
) {
	// ── Stock: customer JWT required (stock.read, stock.bulk-read) ───────────────
	stockCustomerGroup := r.Group("/api/v1/inventory/stock")
	stockCustomerGroup.Use(middleware.JWT(jwtParser, jwtVerifier))
	{
		// FULL
		stockCustomerGroup.POST("/read", svc.StockReadHandler)

		// STUB — 501 NOT_IMPLEMENTED_MVP
		stockCustomerGroup.POST("/bulk-read", svc.StockBulkReadHandler)
	}

	// ── Stock: admin JWT required (stock.adjust — ADMIN only) ───────────────────
	stockAdminGroup := r.Group("/api/v1/inventory/stock")
	stockAdminGroup.Use(middleware.JWT(jwtParser, jwtVerifier))
	{
		// STUB — 501 NOT_IMPLEMENTED_MVP
		// JWT middleware validates the token; handler enforces ADMIN role internally
		// once the full implementation lands (td.json INV-007, INV-008).
		stockAdminGroup.POST("/adjust", svc.StockAdjustHandler)
	}

	// ── Reservation: internal-secret required (Checkout calls only) ─────────────
	reservationGroup := r.Group("/api/v1/inventory/reservation")
	reservationGroup.Use(inventory.InternalAuthMiddleware(internalSecret))
	{
		// FULL (multi-SKU lex-order lock, FOR UPDATE, no-negative invariant, idempotency)
		reservationGroup.POST("/create", svc.ReservationCreateHandler)

		// FULL (handles RESERVED→release and COMMITTED→sold-restock for PR-006)
		reservationGroup.POST("/release", svc.ReservationReleaseHandler)

		// reservation.commit is NOT exposed as HTTP in MVP per TD spec td_ambiguity_notes[0].
		// The commit function is internal, called from the events.payment.completed consumer.
		// Route preserved as a comment for future ops tooling:
		// reservationGroup.POST("/commit", adminOnly(svc.ReservationCommitHandler))
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
		inventory.HeaderInternalSecret,
		refIDHeaderKey,
	}
}
