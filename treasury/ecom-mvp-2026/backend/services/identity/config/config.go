package config

import (
	"log"
	"sync"
	"time"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"

	env "github.com/caarlos0/env/v11"
)

var logFatal = log.Fatal

type Config struct {
	Server  Server
	DB      DB
	JWT     JWT
	Hash    Hash
}

type Server struct {
	Port string `env:"PORT,notEmpty"`
}

type DB struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,notEmpty"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME,notEmpty"`
}

type JWT struct {
	Issuer              string        `env:"JWT_ISSUER,notEmpty"`
	Audience            string        `env:"JWT_AUDIENCE,notEmpty"`
	AccessTokenTTL      time.Duration `env:"JWT_ACCESS_TOKEN_TTL"  envDefault:"15m"`
	RefreshTokenTTL     time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"720h"`
	PrivateKeyPEM       string        `env:"JWT_PRIVATE_KEY_PEM,notEmpty"`
	PublicKeyPEM        string        `env:"JWT_PUBLIC_KEY_PEM,notEmpty"`
}

type Hash struct {
	Pepper string `env:"HASH_PEPPER,notEmpty"`
}

var (
	once   sync.Once
	config Config
)

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
