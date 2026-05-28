package router

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/app/checkout"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/app/checkout/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// Deps holds all shared infrastructure clients for the checkout service.
type Deps struct {
	cfg         config.Config
	pool        *pgxpool.Pool
	jwtParser   token.JWTParser
	jwtVerifier token.JWTVerifier
	service     *checkout.Service
}

// NewDeps constructs all infrastructure clients from config.
// The returned cleanup func must be deferred by the caller.
func NewDeps(ctx context.Context, cfg config.Config) (Deps, func()) {
	// Postgres pool
	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		panic("checkout: postgres pool error: " + err.Error())
	}

	// JWT verifier (ES256 public key)
	jwtVerifier := token.MustNewJWTVerifier(token.JWTVerifierConfig{
		PublicKey: cfg.JWT.PublicKey,
		Alg:       string(token.ES256),
	})
	jwtParser := token.MustNewJWTParser(token.JWTParserConfig{
		Issuer:   cfg.JWT.Issuer,
		Audience: cfg.JWT.Audience,
	})

	// Downstream clients
	cartClient := checkout.NewCartClient(cfg.Downstream.CartBaseURL, cfg.Downstream.CartReadTimeout, cfg.Downstream.CartClearTimeout)
	catalogClient := checkout.NewCatalogClient(cfg.Downstream.CatalogBaseURL, cfg.Downstream.CatalogDetailTimeout)
	inventoryClient := checkout.NewInventoryClient(
		cfg.Downstream.InventoryBaseURL,
		cfg.InternalAuth.SharedSecret,
		cfg.Downstream.InventoryBulkReadTimeout,
		cfg.Downstream.InventoryReservationTimeout,
		cfg.Downstream.InventoryReleaseTimeout,
	)
	paymentClient := checkout.NewPaymentClient(cfg.Downstream.PaymentBaseURL, cfg.InternalAuth.SharedSecret, cfg.Downstream.PaymentIntentTimeout)
	orderClient := checkout.NewOrderClient(
		cfg.Downstream.OrderBaseURL,
		cfg.InternalAuth.SharedSecret,
		cfg.Downstream.OrderCreateTimeout,
		cfg.Downstream.OrderCancelTimeout,
	)
	identityClient := checkout.NewIdentityClient(cfg.Downstream.IdentityBaseURL, cfg.Downstream.IdentityAddressTimeout)

	svc := checkout.NewService(checkout.ServiceConfig{
		DB:               pool,
		IdempotencyStore: access.NewIdempotencyStorage(),
		SagaLogStore:     access.NewSagaLogStorage(),
		Cart:             cartClient,
		Catalog:          catalogClient,
		Inventory:        inventoryClient,
		Payment:          paymentClient,
		Order:            orderClient,
		Identity:         identityClient,
	})

	d := Deps{
		cfg:         cfg,
		pool:        pool,
		jwtParser:   jwtParser,
		jwtVerifier: jwtVerifier,
		service:     svc,
	}

	cleanup := func() {
		pool.Close()
		slog.Info("checkout: postgres pool closed")
	}

	return d, cleanup
}

// Pool exposes the pgxpool for background goroutines in main.go.
func (d Deps) Pool() *pgxpool.Pool { return d.pool }
