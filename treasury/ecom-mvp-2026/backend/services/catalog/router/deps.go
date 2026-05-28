package router

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/database"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"

	"github.com/example/shoppilot/catalog/config"
)

// Deps holds all shared infrastructure clients for the catalog service.
type Deps struct {
	Cfg      config.Config
	DB       *pgxpool.Pool
	Producer kafka.Producer // nil when KAFKA_ENABLED=false
}

// NewDeps constructs all infrastructure clients from config.
// The returned cleanup func must be deferred by the caller.
// No Must* constructors — errors are returned and the caller calls log.Fatal.
func NewDeps(ctx context.Context, cfg config.Config) (Deps, func(), error) {
	// Connect to PostgreSQL via pgx pool (non-Must)
	db, err := database.ConnectPostgresDBWithContext(ctx, database.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
	})
	if err != nil {
		return Deps{}, nil, err
	}

	// Kafka producer: only initialized when KAFKA_ENABLED=true AND brokers are set.
	var producer kafka.Producer
	if cfg.Kafka.Enabled && strings.TrimSpace(cfg.Kafka.Brokers) != "" {
		producer, err = kafka.NewProducer(kafka.ProducerConfig{
			KafkaConf: kafka.NewSyncProducerGuarantee(),
			Brokers:   splitCSV(cfg.Kafka.Brokers),
		})
		if err != nil {
			db.Close()
			return Deps{}, nil, err
		}
	}

	d := Deps{
		Cfg:      cfg,
		DB:       db,
		Producer: producer,
	}

	cleanup := func() {
		db.Close()
		if producer != nil {
			_ = producer.Close()
		}
	}

	return d, cleanup, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
