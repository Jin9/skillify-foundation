package config

import (
	"log"
	"sync"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"

	env "github.com/caarlos0/env/v11"
)

var logFatal = log.Fatal

// Config holds all environment-based configuration for the inventory service.
type Config struct {
	Server   Server
	Database Database
	Kafka    Kafka
	JWT      JWT
	Header   Header
	Sweeper  Sweeper
	Internal Internal
}

type Server struct {
	Port string `env:"PORT,notEmpty"`
}

type Database struct {
	// DB_URL is the full Postgres connection URL.
	// Example: postgres://user:pass@host:5432/inventory
	URL string `env:"DB_URL,notEmpty"`
}

type Kafka struct {
	// KAFKA_ENABLED controls whether the Kafka consumer/producer is started.
	// When false, the outbox relay marks rows published_at=NOW() without actually
	// publishing and logs once per batch. Sweeper still runs.
	Enabled  bool   `env:"KAFKA_ENABLED" envDefault:"false"`
	Brokers  string `env:"KAFKA_BROKERS"`
	GroupID  string `env:"KAFKA_GROUP_ID"`
	Producer struct {
		Brokers string `env:"KAFKA_PRODUCER_BROKERS"`
	}
}

type JWT struct {
	PublicKey string `env:"JWT_PUBLIC_KEY,notEmpty"`
	Issuer    string `env:"JWT_ISSUER" envDefault:"shoppilot-identity"`
	Audience  string `env:"JWT_AUDIENCE" envDefault:"shoppilot-api"`
}

type Header struct {
	RefIDHeaderKey string `env:"REF_ID_HEADER_KEY" envDefault:"X-Ref-Id"`
}

// Internal holds secrets for service-to-service authentication.
// INTERNAL_SHARED_SECRET is checked by InternalAuthMiddleware on the
// reservation.create and reservation.release endpoints (Checkout → Inventory).
type Internal struct {
	SharedSecret string `env:"INTERNAL_SHARED_SECRET,notEmpty"`
}

type Sweeper struct {
	// SWEEPER_CADENCE_SECONDS is how often the expiry-sweep goroutine fires.
	// Source-of-truth for cadence; default 30s per TD spec.
	CadenceSeconds int `env:"SWEEPER_CADENCE_SECONDS" envDefault:"30"`

	// SWEEPER_BATCH_SIZE is the max number of expired reservations processed per tick.
	// Default 200 per TD spec.
	BatchSize int `env:"SWEEPER_BATCH_SIZE" envDefault:"200"`

	// RESERVATION_TTL_MINUTES matches cross-cutting.pricing.constants.
	// The sweeper reads expires_at (written by reservation.create); this constant
	// stays in lockstep so the two agree on the window.
	ReservationTTLMinutes int `env:"RESERVATION_TTL_MINUTES" envDefault:"15"`
}

var once sync.Once
var config Config

// C parses and returns the singleton Config. Fatals on misconfiguration.
func C(envPrefix string) Config {
	once.Do(func() {
		opts := env.Options{}
		if envPrefix != "" {
			opts.Prefix = envPrefix + "_"
		}
		var err error
		config, err = commonconfig.ParseEnv[Config](opts)
		if err != nil {
			logFatal(err)
		}
	})
	return config
}
