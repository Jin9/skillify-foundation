package router

import (
	"context"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/config"
)

// Deps holds all shared infrastructure clients for the order service.
type Deps struct {
	cfg      config.Config
	pool     *pgxpool.Pool
	producer kafka.Producer
	service  *order.Service
	jwtMW    func() (token.JWTParser, token.JWTVerifier)
}

// NewDeps constructs all infrastructure clients from config.
// The returned cleanup func must be deferred by the caller.
func NewDeps(ctx context.Context, cfg config.Config) (Deps, func()) {
	// PostgreSQL pool (pgx).
	pool, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		slog.ErrorContext(ctx, "pgxpool.New failed", "error", err)
		panic(err)
	}

	// Kafka producer (optional — skipped when KAFKA_PRODUCER_BROKERS is empty).
	producer := newProducer(cfg)

	// Storage adapters.
	orderStorage := access.NewOrderStorage(pool)
	itemStorage := access.NewOrderItemStorage(pool)
	historyStorage := access.NewStatusHistoryStorage(pool)
	outboxStorage := access.NewOutboxStorage(pool)
	consumedStorage := access.NewConsumedEventStorage(pool)

	// Domain service.
	svc := order.NewService(pool, orderStorage, itemStorage, historyStorage, outboxStorage, consumedStorage)

	// JWT verifier for common/middleware.JWT.
	jwtMW := func() (token.JWTParser, token.JWTVerifier) {
		parser := token.MustNewJWTParser(token.JWTParserConfig{
			Issuer:   cfg.JWT.Issuer,
			Audience: cfg.JWT.Audience,
		})
		verifier := token.MustNewJWTVerifier(token.JWTVerifierConfig{
			PublicKey: cfg.JWT.PublicKey,
			Alg:       string(token.ES256),
		})
		return parser, verifier
	}

	d := Deps{
		cfg:      cfg,
		pool:     pool,
		producer: producer,
		service:  svc,
		jwtMW:    jwtMW,
	}

	cleanup := func() {
		pool.Close()
		if producer != nil {
			_ = producer.Close()
		}
	}

	return d, cleanup
}

func newProducer(cfg config.Config) kafka.Producer {
	brokers := strings.TrimSpace(cfg.Producer.Brokers)
	if brokers == "" {
		slog.Info("kafka producer skipped (KAFKA_PRODUCER_BROKERS is empty)")
		return nil
	}
	producer := kafka.MustNewProducer(kafka.ProducerConfig{
		KafkaConf: kafka.NewSyncProducerGuarantee(),
		Brokers:   splitCSV(brokers),
	})
	return kafka.WithLoggingFromEnv(producer, cfg.Producer.Env)
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
