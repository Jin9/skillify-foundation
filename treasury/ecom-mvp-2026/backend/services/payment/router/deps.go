package router

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/app/payment"
	paymentaccess "gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/app/payment/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/config"
)

// Deps holds all runtime dependencies for the payment service.
type Deps struct {
	cfg      config.Config
	db       *pgxpool.Pool
	svc      payment.Service
	sweeper  *payment.ExpirySweeper
	verifier token.JWTVerifier
	parser   token.JWTParser
}

// NewDeps constructs all runtime dependencies and returns a cleanup function.
func NewDeps(ctx context.Context, cfg config.Config) (Deps, func()) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		slog.Error("failed to parse postgres config", "error", err)
		panic(err)
	}
	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		panic(err)
	}
	if pingErr := db.Ping(ctx); pingErr != nil {
		slog.Error("failed to ping postgres", "error", pingErr)
		panic(pingErr)
	}

	verifier, err := token.NewJWTVerifier(token.JWTVerifierConfig{
		PublicKey: cfg.JWT.PublicKey,
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

	intentStore := paymentaccess.NewIntentStorage(db)
	dedupStore := paymentaccess.NewCallbackDedupStorage(db)
	outboxStore := paymentaccess.NewOutboxStorage(db)

	intentTTL := time.Duration(cfg.Payment.IntentTTLMinutes) * time.Minute

	svc := payment.NewService(payment.ServiceConfig{
		DB:               db,
		IntentStore:      intentStore,
		DedupStore:       dedupStore,
		OutboxStore:      outboxStore,
		IntentTTL:        intentTTL,
		EmitExpiredEvent: cfg.Payment.EmitExpiredEvent,
	})

	sweepCadence := time.Duration(cfg.Sweeper.CadenceSeconds) * time.Second
	sweeper := payment.NewExpirySweeper(db, intentStore, sweepCadence, cfg.Sweeper.BatchSize)

	cleanup := func() {
		db.Close()
	}

	return Deps{
		cfg:      cfg,
		db:       db,
		svc:      svc,
		sweeper:  sweeper,
		verifier: verifier,
		parser:   parser,
	}, cleanup
}

// Sweeper returns the ExpirySweeper so main.go can call Start(ctx) in a goroutine.
func (d Deps) Sweeper() *payment.ExpirySweeper {
	return d.sweeper
}
