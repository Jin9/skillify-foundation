package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/google/uuid"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"

	"github.com/example/shoppilot/catalog/app/catalog/access"
	"github.com/example/shoppilot/catalog/config"
	"github.com/example/shoppilot/catalog/router"

	_ "time/tzdata"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second
	handlerTimeout           = serverWriteTimeout - (time.Millisecond * 100)

	outboxBatchSize = 100
	kafkaTopic      = "ecom.catalog.events"
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

	d, cleanupDeps, err := router.NewDeps(ctx, cfg)
	if err != nil {
		slog.Error("failed to initialise dependencies", "error", err)
		os.Exit(1)
	}
	defer cleanupDeps()

	// Warn at boot when running in no-kafka mode so operators know relay is a no-op.
	if !cfg.Kafka.Enabled {
		slog.Warn("KAFKA_ENABLED=false — outbox relay will mark events processed without publishing to Kafka")
	}

	// Start outbox relay goroutine. Respects context cancellation.
	outboxStorage := access.NewOutboxStorage(d.DB)
	relayInterval := time.Duration(cfg.Kafka.RelayIntervalSec) * time.Second
	go startOutboxRelay(ctx, outboxStorage, d.Producer, relayInterval, cfg.Kafka.Enabled)

	r := router.New(d, version, commit, handlerTimeout)

	srv := newServer(cfg, r)
	go gracefulShutdown(ctx, srv, gracefulShutdownDuration)

	slog.Info("catalog service running", "port", cfg.Server.Port, "kafka_enabled", cfg.Kafka.Enabled)

	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("HTTP server panicked", "panic", rec)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("HTTP server error", "error", err)
	}

	slog.Info("bye")
}

// startOutboxRelay polls outbox_events WHERE published_at IS NULL on each tick.
//
// When kafkaEnabled=true: sends each event to Kafka and marks published_at after ACK.
// Failed Kafka sends are NOT marked published — they are retried on the next tick.
//
// When kafkaEnabled=false (KAFKA_ENABLED=false): marks rows processed without
// publishing. A warning is logged once at boot. This is the MVP no-kafka mode.
//
// The goroutine respects ctx cancellation and exits cleanly after the current batch.
func startOutboxRelay(
	ctx context.Context,
	store *access.OutboxStorage,
	producer kafka.Producer,
	interval time.Duration,
	kafkaEnabled bool,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("outbox relay: stopping")
			return
		case <-ticker.C:
			if err := relayTick(ctx, store, producer, kafkaEnabled); err != nil {
				slog.Error("outbox relay: tick error", "error", err)
			}
		}
	}
}

// relayTick processes one batch of up to outboxBatchSize pending outbox events.
func relayTick(
	ctx context.Context,
	store *access.OutboxStorage,
	producer kafka.Producer,
	kafkaEnabled bool,
) error {
	rows, err := store.FetchPending(ctx, outboxBatchSize)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	published := make([]uuid.UUID, 0, len(rows))

	for _, row := range rows {
		if kafkaEnabled && producer != nil {
			// Extract partition key (sku) from payload_json
			partitionKey := extractSKU(row.PayloadJSON)

			sendErr := producer.SendMessageWithOption(ctx, kafkaTopic, row.PayloadJSON, kafka.SendMessageOption{
				Key: partitionKey,
				LogAttrs: []slog.Attr{
					slog.String("event_type", row.EventType),
					slog.String("outbox_id", row.ID.String()),
				},
			})
			if sendErr != nil {
				slog.Error("outbox relay: kafka send failed",
					"outbox_id", row.ID,
					"event_type", row.EventType,
					"error", sendErr,
				)
				// Do NOT mark as published — will be retried on next tick.
				continue
			}
		}

		// kafka disabled (no-op) OR kafka ACK received — mark for bulk update
		published = append(published, row.ID)
	}

	if len(published) == 0 {
		return nil
	}

	if err := store.MarkPublished(ctx, published); err != nil {
		return err
	}

	slog.Info("outbox relay: marked published", "count", len(published))
	return nil
}

// extractSKU pulls the sku field from a payload_json byte slice.
// Returns empty string if absent or malformed — Kafka will use round-robin partitioning.
func extractSKU(payloadJSON []byte) string {
	var payload map[string]any
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return ""
	}
	sku, _ := payload["sku"].(string)
	return sku
}

func newServer(cfg config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           handler,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}
}

func gracefulShutdown(ctx context.Context, srv *http.Server, timeout time.Duration) {
	<-ctx.Done()
	slog.Info("shutting down", "timeout", timeout.String())
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}
}
