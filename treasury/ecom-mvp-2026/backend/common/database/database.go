package database

import "time"

const (
	connMaxLifetime       = time.Hour
	connMaxLifetimeJitter = time.Minute * 5
	connMaxIdleTime       = time.Minute * 30
	maxIdleConns          = 10
	maxOpenConns          = 50
	healthCheckPeriod     = time.Minute

	mongoMaxPoolSize        = 100
	mongoConnectTimeout     = 5 * time.Second
	mongoSocketTimeout      = time.Second
	mongoTransactionTimeout = time.Second
)
