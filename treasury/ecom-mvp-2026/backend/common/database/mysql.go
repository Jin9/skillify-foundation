package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (cfg MySQLConfig) DSN() (string, error) {
	if cfg.Host == "" {
		return "", errors.New("database: mysql host is required")
	}
	if cfg.User == "" {
		return "", errors.New("database: mysql user is required")
	}
	if cfg.DBName == "" {
		return "", errors.New("database: mysql dbname is required")
	}

	port := cfg.Port
	if port == "" {
		port = "3306"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return "", errors.New("database: mysql port must be numeric")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True&charset=utf8&loc=Asia%%2FBangkok",
		cfg.User,
		cfg.Password,
		cfg.Host,
		port,
		cfg.DBName,
	)
	return dsn, nil
}

// ConnectMySQLDB creates a connection pool and pings the database.
// It returns an error instead of panicking.
func ConnectMySQLDB(cfg MySQLConfig) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
	return ConnectMySQLDBWithContext(ctx, cfg)
}

// ConnectMySQLDBWithContext creates a connection pool and pings the database using ctx.
func ConnectMySQLDBWithContext(ctx context.Context, cfg MySQLConfig) (*sql.DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// Deprecated: NewMySQL panics on error. Use ConnectMySQLDBWithContext instead.
func NewMySQL(dbUrl string) *sql.DB {
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		log.Panic("could not ping database", pingErr)
	}
	return db
}

// MustNewMySQLWithConfig panics on error. Use ConnectMySQLDB instead for graceful error handling.
func MustNewMySQLWithConfig(cfg MySQLConfig) *sql.DB {
	conn, err := ConnectMySQLDB(cfg)
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}
	return conn
}

func IsMysqlReady() bool {
	// TODO: implement check if mongo is ready
	return true
}
