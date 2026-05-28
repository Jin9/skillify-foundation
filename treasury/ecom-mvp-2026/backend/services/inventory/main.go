package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/app/inventory/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/router"

	_ "embed"
	_ "time/tzdata"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second
	handlerTimeout           = serverWriteTimeout - 100*time.Millisecond
)

// go build -ldflags "-X main.commit=abc123"
var commit string

//go:embed VERSION
var version string

func init() {
	if os.Getenv("GOMAXPROCS") != "" {
		runtime.GOMAXPROCS(0)
	} else {
		runtime.GOMAXPROCS(1)
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1)
	}
}

func main() {
	cfg := config.C(commonconfig.Env)
	_ = logger.New(
		logger.AWSKeyReplacer,
		logger.OpenSearchKeyReplacer,
		logger.CensorReplacer,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// ── Postgres pool (shared by service + sweeper + outbox relay + consumers) ──
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		slog.Error("inventory: failed to parse DB_URL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		slog.Error("inventory: failed to connect to postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// ── Router-layer deps (Kafka producer, config) ────────────────────────────
	d, cleanupDeps := router.NewDeps(ctx, cfg)
	defer cleanupDeps()

	// ── Service (wires all storage on top of the shared pool) ────────────────
	// Concrete access implementations are constructed here (main.go imports both
	// inventory and access; service.go must not import access — that caused the cycle).
	svc := inventory.NewService(inventory.ServiceConfig{
		Pool:            pool,
		StockStorage:    access.NewStockLevelStorage(pool),
		ResvStorage:     access.NewReservationStorage(pool),
		OutboxStorage:   access.NewOutboxStorage(pool),
		ConsumedStorage: access.NewConsumedEventStorage(pool),
		Cfg:             cfg,
	})

	// ── Kafka consumer ────────────────────────────────────────────────────────
	consumerDone, stopConsumer := startConsumer(ctx, svc, cfg)
	defer func() {
		cancel()
		stopConsumer()
		waitForConsumer(consumerDone, gracefulShutdownDuration)
	}()

	// If the consumer exits unexpectedly, trigger global shutdown.
	go func() {
		<-consumerDone
		cancel()
	}()

	// ── Sweeper: expires reservations every SWEEPER_CADENCE_SECONDS ──────────
	sweeper := inventory.NewSweeper(svc, cfg)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("sweeper goroutine panicked", slog.Any("panic", r))
			}
		}()
		sweeper.Run(ctx)
	}()

	// ── Outbox relay: publishes outbox rows every 5s ──────────────────────────
	relay := inventory.NewOutboxRelay(svc, d.Producer(), cfg)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("outbox relay goroutine panicked", slog.Any("panic", r))
			}
		}()
		relay.Run(ctx)
	}()

	// ── HTTP server ───────────────────────────────────────────────────────────
	r := router.New(d, svc, version, commit, handlerTimeout)
	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		<-ctx.Done()
		slog.Info("shutting down HTTP server", "timeout", gracefulShutdownDuration.String())
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownDuration)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Warn("HTTP server Shutdown", "error", err)
		}
	}()

	slog.Info("inventory service starting", "port", cfg.Server.Port, "version", version)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("HTTP server panicked", slog.Any("panic", r))
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("HTTP server ListenAndServe", "error", err)
		return
	}

	slog.Info("inventory service stopped")
}

// ─── Kafka consumer wiring ────────────────────────────────────────────────────

func startConsumer(ctx context.Context, svc *inventory.Service, cfg config.Config) (<-chan struct{}, func()) {
	done := make(chan struct{})

	if !cfg.Kafka.Enabled || strings.TrimSpace(cfg.Kafka.Brokers) == "" {
		slog.Info("kafka consumer disabled (KAFKA_ENABLED=false or no brokers)")
		close(done)
		return done, func() {}
	}

	eventHandlers := svc.EventHandlers()

	consumerCtx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				slog.Error("kafka consumer panicked", slog.Any("panic", r))
			}
		}()

		topics := []string{
			inventory.TopicPaymentEvents,
			inventory.TopicOrderEvents,
			inventory.TopicCatalogEvents,
		}

		group, err := kafka.NewConsumerGroup(kafka.ConsumerConfig{
			Brokers:                  splitCSV(cfg.Kafka.Brokers),
			GroupID:                  cfg.Kafka.GroupID,
			OffsetsInitial:           "latest",
			RebalanceGroupStrategies: "range",
			KafkaConf:                kafka.NewConsumerConfigAtLeastOnce(),
		})
		if err != nil {
			slog.Error("failed to create kafka consumer group", slog.String("error", err.Error()))
			return
		}
		defer func() { _ = group.Close() }()

		processor := kafka.NewEventRouter(eventHandlers)
		handler := kafka.NewConsumerGroupHandler(consumerCtx, processor)

		slog.Info("kafka consumer started", "topics", topics, "group_id", cfg.Kafka.GroupID)
		defer slog.Info("kafka consumer stopped")

		for {
			if err := group.Consume(consumerCtx, topics, handler); err != nil {
				slog.Error("kafka consumer group Consume", slog.String("error", err.Error()))
				return
			}
			if consumerCtx.Err() != nil {
				return
			}
		}
	}()

	return done, cancel
}

func waitForConsumer(done <-chan struct{}, timeout time.Duration) {
	select {
	case <-done:
	case <-time.After(timeout):
		slog.Warn("kafka consumer shutdown timeout", "timeout", timeout.String())
	}
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
