package router

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/database"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/cart/config"
)

// Deps holds all shared infrastructure clients for the cart service.
// Initialized once in main; passed to router wiring.
type Deps struct {
	cfg         config.Config
	httpClient  *http.Client
	pool        *pgxpool.Pool
	jwtParser   token.JWTParser
	jwtVerifier token.JWTVerifier
}

// NewDeps constructs all infrastructure clients from config.
// The returned cleanup func must be deferred by the caller.
func NewDeps(ctx context.Context, cfg config.Config) (Deps, func()) {
	httpClient := httpclient.NewHTTPClient(
		middleware.ForwardRefIDOption,
		httpclient.DebugOption(cfg.HttpClient.EnableLogDebug),
	)

	pool := database.MustNewPostgresDB(newPostgresConfig(cfg))

	// Cart service only VERIFIES tokens — use JWTVerifier (public key) + JWTParser (iss/aud/exp).
	jwtVerifier := token.MustNewJWTVerifier(token.JWTVerifierConfig{
		PublicKey: cfg.JWT.PublicKey,
		Alg:       string(token.ES256),
	})
	jwtParser := token.MustNewJWTParser(token.JWTParserConfig{
		Issuer:   cfg.JWT.Issuer,
		Audience: cfg.JWT.Audience,
	})

	d := Deps{
		cfg:         cfg,
		httpClient:  httpClient,
		pool:        pool,
		jwtParser:   jwtParser,
		jwtVerifier: jwtVerifier,
	}

	cleanup := func() {
		pool.Close()
	}

	return d, cleanup
}

func newPostgresConfig(cfg config.Config) database.PostgresConfig {
	return database.PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
	}
}
