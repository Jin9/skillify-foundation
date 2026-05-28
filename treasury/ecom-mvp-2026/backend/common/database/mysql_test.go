package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQLConfig_DSN(t *testing.T) {
	t.Run("valid config with all fields", func(t *testing.T) {
		cfg := MySQLConfig{
			Host:     "localhost",
			Port:     "3307",
			User:     "root",
			Password: "secret",
			DBName:   "testdb",
		}
		dsn, err := cfg.DSN()
		require.NoError(t, err)
		assert.Contains(t, dsn, "root:secret@tcp(localhost:3307)/testdb")
		assert.Contains(t, dsn, "parseTime=True")
		assert.Contains(t, dsn, "charset=utf8")
	})

	t.Run("valid config with empty password", func(t *testing.T) {
		cfg := MySQLConfig{
			Host:   "localhost",
			Port:   "3306",
			User:   "root",
			DBName: "testdb",
		}
		dsn, err := cfg.DSN()
		require.NoError(t, err)
		assert.Contains(t, dsn, "root:@tcp(localhost:3306)/testdb")
	})

	t.Run("default port when empty", func(t *testing.T) {
		cfg := MySQLConfig{
			Host:   "localhost",
			User:   "root",
			DBName: "testdb",
		}
		dsn, err := cfg.DSN()
		require.NoError(t, err)
		assert.Contains(t, dsn, ":3306)")
	})

	t.Run("error when host is empty", func(t *testing.T) {
		cfg := MySQLConfig{User: "root", DBName: "db"}
		_, err := cfg.DSN()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "host is required")
	})

	t.Run("error when user is empty", func(t *testing.T) {
		cfg := MySQLConfig{Host: "localhost", DBName: "db"}
		_, err := cfg.DSN()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user is required")
	})

	t.Run("error when dbname is empty", func(t *testing.T) {
		cfg := MySQLConfig{Host: "localhost", User: "root"}
		_, err := cfg.DSN()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dbname is required")
	})

	t.Run("error when port is non-numeric", func(t *testing.T) {
		cfg := MySQLConfig{Host: "localhost", User: "root", DBName: "db", Port: "abc"}
		_, err := cfg.DSN()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "port must be numeric")
	})
}

func TestConnectMySQLDB_ValidationError(t *testing.T) {
	_, err := ConnectMySQLDB(MySQLConfig{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")
}

func TestMustNewMySQLWithConfig_PanicsOnInvalidConfig(t *testing.T) {
	assert.Panics(t, func() {
		MustNewMySQLWithConfig(MySQLConfig{})
	})
}

func TestNewMySQL_PanicsOnInvalidURL(t *testing.T) {
	assert.Panics(t, func() {
		NewMySQL("invalid://url")
	})
}

func TestIsMysqlReady(t *testing.T) {
	assert.True(t, IsMysqlReady())
}
