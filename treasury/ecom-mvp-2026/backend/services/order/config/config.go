package config

import (
	"log"
	"sync"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"

	env "github.com/caarlos0/env/v11"
)

var logFatal = log.Fatal

// Config holds all environment-based configuration for the order service.
type Config struct {
	Server        Server
	AccessControl AccessControl
	Header        Header
	JWT           JWT
	Postgres      Postgres
	Consumer      Consumer
	Producer      Producer
	Internal      Internal
}

type Server struct {
	Hostname string `env:"HOSTNAME"`
	Port     string `env:"PORT,notEmpty"`
}

type AccessControl struct {
	AllowOrigin string `env:"ACCESS_CONTROL_ALLOW_ORIGIN"`
}

type Header struct {
	RefIDHeaderKey string `env:"REF_ID_HEADER_KEY,notEmpty"`
}

type JWT struct {
	Issuer    string `env:"JWT_ISSUER,notEmpty"`
	Audience  string `env:"JWT_AUDIENCE,notEmpty"`
	PublicKey string `env:"JWT_PUBLIC_KEY,notEmpty"`
}

// Postgres holds the PostgreSQL DSN (pgx connection URL).
// Example: postgres://user:pass@host:5432/orderdb?sslmode=disable&search_path=order
type Postgres struct {
	URL string `env:"DB_URL,notEmpty"`
}

// Consumer holds Kafka consumer configuration.
type Consumer struct {
	Enabled           bool   `env:"KAFKA_ENABLED" envDefault:"false"`
	Brokers           string `env:"KAFKA_BROKERS"`
	GroupID           string `env:"KAFKA_GROUP_ID"`
	Topics            string `env:"KAFKA_TOPICS"`
	OffsetsInitial    string `env:"KAFKA_OFFSETS_INITIAL" envDefault:"latest"`
	RebalanceStrategy string `env:"KAFKA_REBALANCE_STRATEGY" envDefault:"range"`
}

// Producer holds Kafka producer configuration.
type Producer struct {
	Brokers string `env:"KAFKA_PRODUCER_BROKERS"`
	Env     string `env:"ENV" envDefault:"LOCAL"`
}

// Internal holds secrets for service-to-service authentication.
type Internal struct {
	SharedSecret string `env:"INTERNAL_SHARED_SECRET,notEmpty"`
}

var once sync.Once
var config Config

// C initialises the config once (singleton) using the given env prefix.
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
