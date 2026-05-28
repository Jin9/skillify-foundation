package database

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresConfig_DatabaseURL(t *testing.T) {
	t.Run("valid config with all fields", func(t *testing.T) {
		cfg := PostgresConfig{
			Host:     "localhost",
			Port:     "5433",
			User:     "admin",
			Password: "secret",
			DBName:   "mydb",
		}
		url, err := cfg.DatabaseURL()
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(url, "postgres://"))
		assert.Contains(t, url, "admin:secret@localhost:5433")
		assert.Contains(t, url, "/mydb")
	})

	t.Run("valid config without password", func(t *testing.T) {
		cfg := PostgresConfig{
			Host:   "localhost",
			Port:   "5432",
			User:   "admin",
			DBName: "mydb",
		}
		url, err := cfg.DatabaseURL()
		require.NoError(t, err)
		assert.Contains(t, url, "admin@localhost:5432")
		assert.NotContains(t, url, "admin:")
	})

	t.Run("default port when empty", func(t *testing.T) {
		cfg := PostgresConfig{
			Host:   "localhost",
			User:   "admin",
			DBName: "mydb",
		}
		url, err := cfg.DatabaseURL()
		require.NoError(t, err)
		assert.Contains(t, url, ":5432")
	})

	t.Run("error when host is empty", func(t *testing.T) {
		cfg := PostgresConfig{User: "admin", DBName: "db"}
		_, err := cfg.DatabaseURL()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "host is required")
	})

	t.Run("error when user is empty", func(t *testing.T) {
		cfg := PostgresConfig{Host: "localhost", DBName: "db"}
		_, err := cfg.DatabaseURL()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user is required")
	})

	t.Run("error when dbname is empty", func(t *testing.T) {
		cfg := PostgresConfig{Host: "localhost", User: "admin"}
		_, err := cfg.DatabaseURL()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dbname is required")
	})

	t.Run("error when port is non-numeric", func(t *testing.T) {
		cfg := PostgresConfig{Host: "localhost", User: "admin", DBName: "db", Port: "xyz"}
		_, err := cfg.DatabaseURL()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "port must be numeric")
	})
}

func TestParseConfig(t *testing.T) {
	t.Run("valid config sets pool options", func(t *testing.T) {
		cfg := PostgresConfig{
			Host:   "localhost",
			User:   "admin",
			DBName: "mydb",
		}
		poolCfg, err := ParseConfig(cfg)
		require.NoError(t, err)
		assert.Equal(t, int32(maxOpenConns), poolCfg.MaxConns)
		assert.Equal(t, connMaxLifetime, poolCfg.MaxConnLifetime)
		assert.Equal(t, connMaxLifetimeJitter, poolCfg.MaxConnLifetimeJitter)
		assert.Equal(t, connMaxIdleTime, poolCfg.MaxConnIdleTime)
		assert.Equal(t, healthCheckPeriod, poolCfg.HealthCheckPeriod)
	})

	t.Run("propagates validation error", func(t *testing.T) {
		_, err := ParseConfig(PostgresConfig{})
		require.Error(t, err)
	})
}

func TestMustParseConfig(t *testing.T) {
	t.Run("panics on invalid config", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseConfig(PostgresConfig{})
		})
	})

	t.Run("returns config on valid input", func(t *testing.T) {
		assert.NotPanics(t, func() {
			cfg := MustParseConfig(PostgresConfig{
				Host:   "localhost",
				User:   "admin",
				DBName: "mydb",
			})
			assert.NotNil(t, cfg)
		})
	})
}

func TestConnectPostgresDB_ValidationError(t *testing.T) {
	_, err := ConnectPostgresDB(PostgresConfig{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")
}

func TestMustNewPostgresDB_PanicsOnInvalidConfig(t *testing.T) {
	assert.Panics(t, func() {
		MustNewPostgresDB(PostgresConfig{})
	})
}
