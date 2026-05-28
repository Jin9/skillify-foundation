package database

import (
	"context"
	"errors"
	"log"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultConnectTimeout = 10 * time.Second

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (cfg PostgresConfig) DatabaseURL() (string, error) {
	if cfg.Host == "" {
		return "", errors.New("database: postgres host is required")
	}
	if cfg.User == "" {
		return "", errors.New("database: postgres user is required")
	}
	if cfg.DBName == "" {
		return "", errors.New("database: postgres dbname is required")
	}

	port := cfg.Port
	if port == "" {
		port = "5432"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return "", errors.New("database: postgres port must be numeric")
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   cfg.Host + ":" + port,
		Path:   path.Join("/", cfg.DBName),
	}
	if cfg.Password != "" {
		u.User = url.UserPassword(cfg.User, cfg.Password)
	} else {
		u.User = url.User(cfg.User)
	}

	return u.String(), nil
}

// ParseConfig parses and prepares a pgxpool config without panicking.
func ParseConfig(cfg PostgresConfig) (*pgxpool.Config, error) {
	dbURL, err := cfg.DatabaseURL()
	if err != nil {
		return nil, err
	}

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(maxOpenConns)
	dbConfig.MaxConnLifetime = connMaxLifetime
	dbConfig.MaxConnLifetimeJitter = connMaxLifetimeJitter
	dbConfig.MaxConnIdleTime = connMaxIdleTime
	dbConfig.HealthCheckPeriod = healthCheckPeriod

	return dbConfig, nil
}

// MustParseConfig panics on error. Use ParseConfig instead for graceful error handling.
func MustParseConfig(cfg PostgresConfig) *pgxpool.Config {
	dbConfig, err := ParseConfig(cfg)
	if err != nil {
		log.Panic("Failed to create a config, error: ", err)
	}
	return dbConfig
}

// ConnectPostgresDB creates a connection pool and pings the database.
// It returns an error instead of panicking.
func ConnectPostgresDB(cfg PostgresConfig) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
	return ConnectPostgresDBWithContext(ctx, cfg)
}

// ConnectPostgresDBWithContext creates a connection pool and pings the database using ctx.
func ConnectPostgresDBWithContext(ctx context.Context, cfg PostgresConfig) (*pgxpool.Pool, error) {
	poolCfg, err := ParseConfig(cfg)
	if err != nil {
		return nil, err
	}

	connPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := connPool.Ping(ctx); err != nil {
		connPool.Close()
		return nil, err
	}

	return connPool, nil
}

// MustNewPostgresDB panics on error. Use ConnectPostgresDB instead for graceful error handling.
func MustNewPostgresDB(cfg PostgresConfig) *pgxpool.Pool {
	connPool, err := ConnectPostgresDB(cfg)
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}
	return connPool
}
