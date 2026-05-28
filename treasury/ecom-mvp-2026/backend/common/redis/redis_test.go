package redis

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeRedisConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := normalizeRedisConfig(RedisUniversalCfgs{})
		assert.Equal(t, defaultReadTimeout, cfg.ReadTimeout)
		assert.Equal(t, defaultMaxRetries, cfg.MaxRetries)
		assert.NotNil(t, cfg.TLSConfig)
		assert.True(t, cfg.TLSConfig.InsecureSkipVerify)
	})

	t.Run("custom overrides", func(t *testing.T) {
		customTLS := &tls.Config{InsecureSkipVerify: false}
		cfg := normalizeRedisConfig(RedisUniversalCfgs{
			ReadTimeout:   10 * time.Second,
			MaxRetries:    5,
			TLSConfig:     customTLS,
		})
		assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
		assert.Equal(t, 5, cfg.MaxRetries)
		assert.Equal(t, customTLS, cfg.TLSConfig)
	})

	t.Run("insecure skip verify override", func(t *testing.T) {
		cfg := normalizeRedisConfig(RedisUniversalCfgs{
			TLSInsecureSkipVerifySet: true,
			TLSInsecureSkipVerify:    false,
		})
		assert.NotNil(t, cfg.TLSConfig)
		assert.False(t, cfg.TLSConfig.InsecureSkipVerify)
	})
}

func TestNormalizeRedisClusterConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := normalizeRedisClusterConfig(RedisClusterCfgs{})
		assert.Equal(t, defaultReadTimeout, cfg.ReadTimeout)
		assert.Equal(t, defaultMaxRetries, cfg.MaxRetries)
		assert.True(t, cfg.RouteRandomly)
		assert.NotNil(t, cfg.TLSConfig)
		assert.True(t, cfg.TLSConfig.InsecureSkipVerify)
	})

	t.Run("insecure skip verify override", func(t *testing.T) {
		cfg := normalizeRedisClusterConfig(RedisClusterCfgs{
			TLSInsecureSkipVerifySet: true,
			TLSInsecureSkipVerify:    false,
		})
		assert.NotNil(t, cfg.TLSConfig)
		assert.False(t, cfg.TLSConfig.InsecureSkipVerify)
	})
}

func TestNormalizeRedisFailOverConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := normalizeRedisFailOverConfig(RedisFailOverCfgs{})
		assert.Equal(t, defaultReadTimeout, cfg.ReadTimeout)
		assert.Equal(t, defaultMaxRetries, cfg.MaxRetries)
		assert.True(t, cfg.RouteRandomly)
		assert.NotNil(t, cfg.TLSConfig)
		assert.True(t, cfg.TLSConfig.InsecureSkipVerify)
	})

	t.Run("insecure skip verify override", func(t *testing.T) {
		cfg := normalizeRedisFailOverConfig(RedisFailOverCfgs{
			TLSInsecureSkipVerifySet: true,
			TLSInsecureSkipVerify:    false,
		})
		assert.NotNil(t, cfg.TLSConfig)
		assert.False(t, cfg.TLSConfig.InsecureSkipVerify)
	})
}

func TestConnect_Errors(t *testing.T) {
	t.Run("universal client ping failure", func(t *testing.T) {
		_, err := Connect("invalid:address", "pass")
		assert.Error(t, err)
	})

	t.Run("cluster client ping failure", func(t *testing.T) {
		_, err := ConnectCluster([]string{"invalid:address"})
		assert.Error(t, err)
	})
}

func TestMustNewPanics(t *testing.T) {
	t.Run("universal", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNew("invalid:address", "pass")
		})
	})
	t.Run("universal with config", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewWithConfig(RedisUniversalCfgs{Addr: "invalid:address"})
		})
	})
	t.Run("cluster", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewCluster([]string{"invalid:address"})
		})
	})
}
