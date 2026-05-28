package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
)

var logFatal = log.Fatal

// Config holds all environment-based configuration for the cart service.
type Config struct {
	Server        Server
	AccessControl AccessControl
	Header        Header
	JWT           JWT
	DB            DB
	Catalog       Catalog
	Inventory     Inventory
	Fanout        Fanout
	HttpClient    HttpClient
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
	PublicKey string `env:"SECRET_JWT_PUBLIC_KEY,notEmpty"`
}

// DB holds PostgreSQL connection config.
// DB_URL takes precedence if set; otherwise host/port/user/pass/name are used.
type DB struct {
	URL      string `env:"DB_URL"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"SECRET_DB_USERNAME" envDefault:"cart"`
	Password string `env:"SECRET_DB_PASSWORD"`
	Name     string `env:"DB_NAME" envDefault:"cart"`
}

type Catalog struct {
	BaseURL string `env:"CATALOG_BASE_URL,notEmpty"`
}

type Inventory struct {
	BaseURL string `env:"INVENTORY_BASE_URL,notEmpty"`
}

type Fanout struct {
	TimeoutMS   int `env:"FANOUT_TIMEOUT_MS" envDefault:"3000"`
	MaxItems    int `env:"FANOUT_MAX_ITEMS" envDefault:"50"`
	Concurrency int `env:"FANOUT_CONCURRENCY" envDefault:"10"`
}

type HttpClient struct {
	EnableLogDebug bool `env:"HTTP_CLIENT_ENABLE_LOG_DEBUG" envDefault:"false"`
}

var once sync.Once
var config Config

func prefix(e string) string {
	if e == "" {
		return ""
	}
	return fmt.Sprintf("%s_", e)
}

// C returns the singleton Config, parsed once from environment variables.
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
