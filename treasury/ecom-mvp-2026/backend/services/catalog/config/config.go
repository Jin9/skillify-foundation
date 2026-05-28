package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"

	env "github.com/caarlos0/env/v11"
)

var logFatal = log.Fatal

// Config holds all environment-based configuration for the catalog service.
type Config struct {
	Server   Server
	Header   Header
	JWT      JWT
	Postgres Postgres
	Kafka    Kafka
	HttpClient HttpClient
	Migrations Migrations
}

type Server struct {
	Hostname string `env:"HOSTNAME"`
	Port     string `env:"PORT,notEmpty"`
}

type Header struct {
	RefIDHeaderKey string `env:"REF_ID_HEADER_KEY,notEmpty"`
}

type AccessControl struct {
	AllowOrigin string `env:"ACCESS_CONTROL_ALLOW_ORIGIN"`
}

type JWT struct {
	Issuer    string `env:"JWT_ISSUER,notEmpty"`
	Audience  string `env:"JWT_AUDIENCE,notEmpty"`
	PublicKey string `env:"JWT_PUBLIC_KEY,notEmpty"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST,notEmpty"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER,notEmpty"`
	Password string `env:"SECRET_POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB,notEmpty"`
}

type Kafka struct {
	Enabled         bool          `env:"KAFKA_ENABLED" envDefault:"false"`
	Brokers         string        `env:"KAFKA_PRODUCER_BROKERS"`
	Topic           string        `env:"KAFKA_CATALOG_TOPIC" envDefault:"ecom.catalog.events"`
	RelayIntervalSec int          `env:"OUTBOX_RELAY_INTERVAL_SECONDS" envDefault:"5"`
}

type HttpClient struct {
	InventoryBaseURL string        `env:"INVENTORY_BASE_URL" envDefault:"http://inventory-svc"`
	Timeout          time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"2s"`
}

type Migrations struct {
	AutoApply bool `env:"MIGRATIONS_AUTO_APPLY" envDefault:"false"`
}

var (
	once   sync.Once
	config Config
)

func prefix(e string) string {
	if e == "" {
		return ""
	}
	return fmt.Sprintf("%s_", e)
}

func C(envPrefix string) Config {
	once.Do(func() {
		opts := env.Options{
			Prefix: prefix(envPrefix),
		}

		var err error
		config, err = commonconfig.ParseEnv[Config](opts)
		if err != nil {
			logFatal(err)
		}
	})
	return config
}
