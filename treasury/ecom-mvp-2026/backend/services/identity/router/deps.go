package router

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	commondb "gitlab.com/b2c-e-commerce-platform/platform/backend/common/database"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity/config"
)

type Deps struct {
	cfg           config.Config
	db            *pgxpool.Pool
	accessSigner  token.JWTSigner // TTL = AccessTokenTTL (15m)
	refreshSigner token.JWTSigner // TTL = RefreshTokenTTL (30d)
	verifier      token.JWTVerifier
	parser        token.JWTParser
}

func NewDeps(ctx context.Context, cfg config.Config) (Deps, func()) {
	db, err := commondb.ConnectPostgresDBWithContext(ctx, commondb.PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
	})
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		panic(err)
	}

	accessSigner, err := token.NewJWTSigner(token.JWTSignerConfig{
		PrivateKey: cfg.JWT.PrivateKeyPEM,
		Alg:        string(token.ES256),
		Issuer:     cfg.JWT.Issuer,
		Audience:   cfg.JWT.Audience,
		Expire:     cfg.JWT.AccessTokenTTL,
	})
	if err != nil {
		panic(err)
	}

	refreshSigner, err := token.NewJWTSigner(token.JWTSignerConfig{
		PrivateKey: cfg.JWT.PrivateKeyPEM,
		Alg:        string(token.ES256),
		Issuer:     cfg.JWT.Issuer,
		Audience:   cfg.JWT.Audience,
		Expire:     cfg.JWT.RefreshTokenTTL,
	})
	if err != nil {
		panic(err)
	}

	verifier, err := token.NewJWTVerifier(token.JWTVerifierConfig{
		PublicKey: cfg.JWT.PublicKeyPEM,
		Alg:       string(token.ES256),
	})
	if err != nil {
		panic(err)
	}

	parser, err := token.NewJWTParser(token.JWTParserConfig{
		Issuer:   cfg.JWT.Issuer,
		Audience: cfg.JWT.Audience,
	})
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		db.Close()
	}

	return Deps{
		cfg:           cfg,
		db:            db,
		accessSigner:  accessSigner,
		refreshSigner: refreshSigner,
		verifier:      verifier,
		parser:        parser,
	}, cleanup
}
