package inventory

import (
	"context"
	"log/slog"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/config"
)

const (
	outboxRelayInterval = 5 * time.Second
	outboxRelayBatch    = 100
	outboxTopic         = "ecom.inventory.events"
)

// OutboxRelay polls outbox_events WHERE published_at IS NULL and publishes to Kafka.
// If KAFKA_ENABLED=false, it marks rows published_at=NOW() without publishing
// (logs once per batch at Info level).
//
// Delivery guarantee: at-least-once. The relayer:
//  1. Begins a tx; SELECT FOR UPDATE SKIP LOCKED LIMIT 100.
//  2. For each row: publish to Kafka (or skip if disabled).
//  3. UPDATE published_at=NOW() inside the same tx.
//  4. Commit.
type OutboxRelay struct {
	svc      *Service
	producer kafka.Producer
	cfg      config.Config
}

// NewOutboxRelay constructs an OutboxRelay.
// producer may be nil when KAFKA_ENABLED=false.
func NewOutboxRelay(svc *Service, producer kafka.Producer, cfg config.Config) *OutboxRelay {
	return &OutboxRelay{svc: svc, producer: producer, cfg: cfg}
}

// Run starts the relay loop. Blocks until ctx is cancelled.
// Start as a goroutine from main.
func (r *OutboxRelay) Run(ctx context.Context) {
	ticker := time.NewTicker(outboxRelayInterval)
	defer ticker.Stop()

	slog.InfoContext(ctx, "outbox relay started",
		slog.Bool("kafka_enabled", r.cfg.Kafka.Enabled),
		slog.Duration("interval", outboxRelayInterval),
		slog.Int("batch_size", outboxRelayBatch),
	)

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "outbox relay stopped")
			return
		case <-ticker.C:
			r.relay(ctx)
		}
	}
}

// relay processes one batch of unpublished outbox events.
func (r *OutboxRelay) relay(ctx context.Context) {
	tx, err := r.svc.Pool.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "outbox relay begin tx", slog.String("error", err.Error()))
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	events, err := r.svc.OutboxStorage.FetchUnpublished(ctx, tx, outboxRelayBatch)
	if err != nil {
		slog.ErrorContext(ctx, "outbox relay fetch", slog.String("error", err.Error()))
		return
	}

	if len(events) == 0 {
		_ = tx.Rollback(ctx)
		return
	}

	publishedCount := 0
	skippedCount := 0

	for _, evt := range events {
		if r.producer != nil && r.cfg.Kafka.Enabled {
			// Publish to Kafka with aggregate_id as the partition key.
			if err := r.producer.SendMessageWithOption(ctx, outboxTopic, evt.PayloadJSON, kafka.SendMessageOption{
				Key: evt.AggregateID, // orderId → partition key per contracts.json
			}); err != nil {
				slog.ErrorContext(ctx, "outbox relay publish",
					slog.String("outbox_id", evt.ID.String()),
					slog.String("event_type", evt.EventType),
					slog.String("error", err.Error()),
				)
				// Abort this batch on publish failure — row stays unpublished for next tick.
				return
			}
		} else {
			// KAFKA_ENABLED=false: mark published without actually publishing.
			skippedCount++
		}

		if err := r.svc.OutboxStorage.MarkPublished(ctx, tx, evt.ID); err != nil {
			slog.ErrorContext(ctx, "outbox relay mark published",
				slog.String("outbox_id", evt.ID.String()),
				slog.String("error", err.Error()),
			)
			return
		}
		publishedCount++
	}

	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "outbox relay commit", slog.String("error", err.Error()))
		return
	}

	if skippedCount > 0 {
		slog.InfoContext(ctx, "outbox relay: kafka disabled, marked events published without sending",
			slog.Int("marked_published", skippedCount),
		)
	}
	if publishedCount > 0 {
		slog.InfoContext(ctx, "outbox relay batch complete",
			slog.Int("published", publishedCount),
		)
	}
}
