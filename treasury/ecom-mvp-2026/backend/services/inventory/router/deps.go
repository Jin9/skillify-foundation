package router

import (
	"context"
	"log/slog"
	"strings"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/config"
)

// Deps holds all shared infrastructure clients for the inventory service.
// The Postgres pool is managed in main.go and passed to the service directly.
// Deps carries the Kafka producer, JWT middleware factory, and config for the router/HTTP layer.
type Deps struct {
	cfg      config.Config
	producer kafka.Producer
	jwtMW    func() (token.JWTParser, token.JWTVerifier)
}

// NewDeps constructs the router-layer infrastructure from config.
// The returned cleanup func must be deferred by the caller.
func NewDeps(ctx context.Context, cfg config.Config) (Deps, func()) {
	var producer kafka.Producer
	if cfg.Kafka.Enabled && strings.TrimSpace(cfg.Kafka.Producer.Brokers) != "" {
		p, err := kafka.NewProducer(kafka.ProducerConfig{
			KafkaConf: kafka.NewSyncProducerGuarantee(),
			Brokers:   splitCSV(cfg.Kafka.Producer.Brokers),
		})
		if err != nil {
			slog.Error("inventory: failed to create kafka producer", slog.String("error", err.Error()))
			panic(err)
		}
		producer = p
		slog.Info("inventory: kafka producer started")
	} else {
		slog.Info("inventory: kafka producer disabled (KAFKA_ENABLED=false or no brokers)")
	}

	_ = ctx // reserved for future use (e.g. tracing setup)

	// JWT parser + verifier for customer and admin route groups.
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
		producer: producer,
		jwtMW:    jwtMW,
	}

	cleanup := func() {
		if producer != nil {
			_ = producer.Close()
		}
	}

	return d, cleanup
}

// Producer returns the Kafka producer (may be nil when KAFKA_ENABLED=false).
func (d Deps) Producer() kafka.Producer {
	return d.producer
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
