package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultReadTimeout = 2 * time.Second
	defaultMaxRetries  = 2
	defaultPingTimeout = 5 * time.Second
)

type RedisUniversalCfgs struct {
	Addr     string
	Password string

	ReadTimeout time.Duration
	MaxRetries  int
	TLSConfig   *tls.Config

	// TLSInsecureSkipVerify is used only when TLSConfig is nil.
	// To override the default, set TLSInsecureSkipVerifySet=true.
	TLSInsecureSkipVerify    bool
	TLSInsecureSkipVerifySet bool

	DB int
}

type RedisClusterCfgs struct {
	Addrs []string

	ReadTimeout   time.Duration
	MaxRetries    int
	RouteRandomly bool
	TLSConfig     *tls.Config

	TLSInsecureSkipVerify    bool
	TLSInsecureSkipVerifySet bool
}

type RedisFailOverCfgs struct {
	MasterName string
	Addrs      []string

	ReadTimeout   time.Duration
	MaxRetries    int
	RouteRandomly bool
	TLSConfig     *tls.Config

	TLSInsecureSkipVerify    bool
	TLSInsecureSkipVerifySet bool
}

func applyConnectionDefaults(readTimeout *time.Duration, maxRetries *int, tlsCfg **tls.Config, insecure bool, insecureSet bool) {
	if *readTimeout <= 0 {
		*readTimeout = defaultReadTimeout
	}
	if *maxRetries <= 0 {
		*maxRetries = defaultMaxRetries
	}
	if *tlsCfg == nil {
		skip := true // preserve historical default
		if insecureSet {
			skip = insecure
		}
		*tlsCfg = &tls.Config{InsecureSkipVerify: skip}
	}
}

func normalizeRedisConfig(cfgs RedisUniversalCfgs) RedisUniversalCfgs {
	applyConnectionDefaults(&cfgs.ReadTimeout, &cfgs.MaxRetries, &cfgs.TLSConfig, cfgs.TLSInsecureSkipVerify, cfgs.TLSInsecureSkipVerifySet)
	return cfgs
}

func normalizeRedisClusterConfig(cfgs RedisClusterCfgs) RedisClusterCfgs {
	applyConnectionDefaults(&cfgs.ReadTimeout, &cfgs.MaxRetries, &cfgs.TLSConfig, cfgs.TLSInsecureSkipVerify, cfgs.TLSInsecureSkipVerifySet)
	if !cfgs.RouteRandomly {
		cfgs.RouteRandomly = true
	}
	return cfgs
}

func normalizeRedisFailOverConfig(cfgs RedisFailOverCfgs) RedisFailOverCfgs {
	applyConnectionDefaults(&cfgs.ReadTimeout, &cfgs.MaxRetries, &cfgs.TLSConfig, cfgs.TLSInsecureSkipVerify, cfgs.TLSInsecureSkipVerifySet)
	if !cfgs.RouteRandomly {
		cfgs.RouteRandomly = true
	}
	return cfgs
}

func Connect(addr string, password string) (redis.UniversalClient, error) {
	return ConnectWithConfig(RedisUniversalCfgs{Addr: addr, Password: password})
}

func ConnectWithConfig(cfgs RedisUniversalCfgs) (redis.UniversalClient, error) {
	cfgs = normalizeRedisConfig(cfgs)
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:       []string{cfgs.Addr},
		Password:    cfgs.Password,
		ReadTimeout: cfgs.ReadTimeout,
		MaxRetries:  cfgs.MaxRetries,
		DB:          cfgs.DB,
		TLSConfig:   cfgs.TLSConfig,
	})

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}

func New(addr string, password string) (redis.UniversalClient, error) {
	return Connect(addr, password)
}

// MustNew is a convenience wrapper that panics on error.
func MustNew(addr string, password string) redis.UniversalClient {
	rdb, err := New(addr, password)
	if err != nil {
		panic(err)
	}
	return rdb
}

func NewWithConfig(cfgs RedisUniversalCfgs) (redis.UniversalClient, error) {
	return ConnectWithConfig(cfgs)
}

// MustNewWithConfig is a convenience wrapper that panics on error.
func MustNewWithConfig(cfgs RedisUniversalCfgs) redis.UniversalClient {
	rdb, err := NewWithConfig(cfgs)
	if err != nil {
		panic(err)
	}
	return rdb
}

func ConnectClusterWithConfig(cfgs RedisClusterCfgs) (*redis.ClusterClient, error) {
	cfgs = normalizeRedisClusterConfig(cfgs)
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         cfgs.Addrs,
		ReadTimeout:   cfgs.ReadTimeout,
		MaxRetries:    cfgs.MaxRetries,
		RouteRandomly: cfgs.RouteRandomly,
		TLSConfig:     cfgs.TLSConfig,
	})

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()

	if err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	}); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}

func ConnectCluster(addrs []string) (*redis.ClusterClient, error) {
	return ConnectClusterWithConfig(RedisClusterCfgs{Addrs: addrs})
}

func NewCluster(addrs []string) (*redis.ClusterClient, error) {
	return ConnectCluster(addrs)
}

// MustNewCluster is a convenience wrapper that panics on error.
func MustNewCluster(addrs []string) *redis.ClusterClient {
	rdb, err := NewCluster(addrs)
	if err != nil {
		panic(err)
	}
	return rdb
}

func ConnectFailOverWithConfig(cfgs RedisFailOverCfgs) (*redis.Client, error) {
	cfgs = normalizeRedisFailOverConfig(cfgs)
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfgs.MasterName,
		SentinelAddrs: cfgs.Addrs,
		ReadTimeout:   cfgs.ReadTimeout,
		MaxRetries:    cfgs.MaxRetries,
		RouteRandomly: cfgs.RouteRandomly,
		TLSConfig:     cfgs.TLSConfig,
	})

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}

func ConnectFailOver(masterName string, addrs []string) (*redis.Client, error) {
	return ConnectFailOverWithConfig(RedisFailOverCfgs{MasterName: masterName, Addrs: addrs})
}

func NewFailOver(masterName string, addrs []string) (*redis.Client, error) { // Sentinel
	return ConnectFailOver(masterName, addrs)
}

// MustNewFailOver is a convenience wrapper that panics on error.
func MustNewFailOver(masterName string, addrs []string) *redis.Client { // Sentinel
	rdb, err := NewFailOver(masterName, addrs)
	if err != nil {
		panic(err)
	}
	return rdb
}
