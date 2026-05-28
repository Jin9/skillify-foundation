package config

import (
	"log"
	"sync"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"

	env "github.com/caarlos0/env/v11"
)

var logFatal = log.Fatal

// Config holds all environment-based configuration for the payment service.
type Config struct {
	Server   Server
	Database Database
	Kafka    Kafka
	JWT      JWT
	Payment  Payment
	Sweeper  Sweeper
}

type Server struct {
	Port string `env:"PORT,notEmpty"`
}

type Database struct {
	// DB_URL is the full Postgres connection URL.
	// Example: postgres://user:pass@host:5432/payment
	URL string `env:"DB_URL,notEmpty"`
}

type Kafka struct {
	// KAFKA_ENABLED controls whether the outbox relay publishes to Kafka.
	// When false, the outbox relay marks rows published_at=NOW() without
	// actually publishing and logs once per batch.
	Enabled         bool   `env:"KAFKA_ENABLED" envDefault:"false"`
	Brokers         string `env:"KAFKA_BROKERS"`
	ProducerBrokers string `env:"KAFKA_PRODUCER_BROKERS"`
}

type JWT struct {
	PublicKey string `env:"JWT_PUBLIC_KEY,notEmpty"`
	Issuer    string `env:"JWT_ISSUER" envDefault:"shoppilot-identity"`
	Audience  string `env:"JWT_AUDIENCE" envDefault:"shoppilot-api"`
}

type Payment struct {
	// MOCK_PROVIDER_SECRET is the shared secret for HMAC-SHA256 webhook signature verification.
	// Header: X-Mock-Provider-Signature: hex(HMAC-SHA256(rawBody, secret)).
	// Env name LOCKED per td.json config_env_vars (ADR-008 alignment).
	CallbackHMACSecret string `env:"MOCK_PROVIDER_SECRET,notEmpty"`

	// INTERNAL_SHARED_SECRET is the shared secret for X-Internal-Secret header on
	// payment.intent.create (called by checkout). MVP internal auth.
	// Renamed from INTERNAL_SECRET to align with Order/Inventory (REV-L2-002).
	InternalSecret string `env:"INTERNAL_SHARED_SECRET,notEmpty"`

	// INTENT_TTL_MINUTES is the payment-intent lifetime. Must match reservation TTL.
	IntentTTLMinutes int `env:"INTENT_TTL_MINUTES" envDefault:"15"`

	// PAYMENT_EMIT_EXPIRED_EVENT gates future events.payment.expired emission.
	// Default false in MVP per expired_event_decision.
	EmitExpiredEvent bool `env:"PAYMENT_EMIT_EXPIRED_EVENT" envDefault:"false"`
}

type Sweeper struct {
	// SWEEPER_CADENCE_SECONDS is how often the expiry-sweep goroutine fires.
	CadenceSeconds int `env:"SWEEPER_CADENCE_SECONDS" envDefault:"60"`

	// SWEEPER_BATCH_SIZE is the max number of expired intents processed per tick.
	BatchSize int `env:"SWEEPER_BATCH_SIZE" envDefault:"200"`
}

var once sync.Once
var config Config

// C parses and returns the singleton Config. Fatals on misconfiguration.
func C() Config {
	once.Do(func() {
		opts := env.Options{}
		var err error
		config, err = commonconfig.ParseEnv[Config](opts)
		if err != nil {
			logFatal(err)
		}
	})
	return config
}
