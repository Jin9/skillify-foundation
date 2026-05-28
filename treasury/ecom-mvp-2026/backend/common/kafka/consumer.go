package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
)

func newConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	return config
}

// ConsumerConfig holds the settings needed to create a Kafka consumer group.
type ConsumerConfig struct {
	Brokers                  []string
	GroupID                  string
	Topics                   []string
	OffsetsInitial           string
	RebalanceGroupStrategies string
	EnableSSL                bool
	CACertPEM                string
	ClientCertPEM            string
	ClientKeyPEM             string
	EnableSASL               bool
	SaslMechanism            sarama.SASLMechanism
	SaslUsername             string
	SaslPassword             string
	KafkaConf                *sarama.Config
}

var balanceStrategy = map[string]sarama.BalanceStrategy{
	"":           sarama.NewBalanceStrategyRange(),
	"range":      sarama.NewBalanceStrategyRange(),
	"sticky":     sarama.NewBalanceStrategySticky(),
	"roundrobin": sarama.NewBalanceStrategyRoundRobin(),
}

var offsetsInitial = map[string]int64{
	"":         sarama.OffsetNewest,
	"earliest": sarama.OffsetOldest,
	"latest":   sarama.OffsetNewest,
	"oldest":   sarama.OffsetOldest,
	"newest":   sarama.OffsetNewest,
}

func offsetInitialFrom(value string) int64 {
	offset, ok := offsetsInitial[value]
	if !ok {
		return sarama.OffsetNewest
	}
	return offset
}

func balanceStrategyFrom(value string) sarama.BalanceStrategy {
	strategy, ok := balanceStrategy[value]
	if !ok || strategy == nil {
		return sarama.NewBalanceStrategyRange()
	}
	return strategy
}

func NewConsumerConfigNoGuarantee() *sarama.Config {
	return newConfig()
}

func NewConsumerConfigAtMostOnce() *sarama.Config {
	config := newConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second

	return config
}

func NewConsumerConfigAtLeastOnce() *sarama.Config {
	config := newConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false
	return config
}

// Deprecated: NewDefaultConsumerGroup is identical to NewConsumerConfigAtMostOnce. Use that instead.
func NewDefaultConsumerGroup() *sarama.Config {
	config := newConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second

	return config
}

// NewConsumerGroup creates a Kafka consumer group from the given config.
// Use NewConsumerConfigAtLeastOnce or NewConsumerConfigAtMostOnce to pre-configure delivery semantics.
func NewConsumerGroup(cfg ConsumerConfig) (sarama.ConsumerGroup, error) {
	kafkaConf := cfg.KafkaConf
	if kafkaConf == nil {
		return nil, fmt.Errorf("kafka config must not be nil")
	}
	kafkaConf.Consumer.Offsets.Initial = offsetInitialFrom(cfg.OffsetsInitial)
	kafkaConf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{ // default is [sarama.BalanceStrategyRange, sarama.BalanceStrategyRoundRobin]
		balanceStrategyFrom(cfg.RebalanceGroupStrategies),
	}
	if err := applyTLSConfig(kafkaConf, cfg.EnableSSL, cfg.CACertPEM, cfg.ClientCertPEM, cfg.ClientKeyPEM); err != nil {
		return nil, err
	}
	if err := applySASLConfig(kafkaConf, cfg.EnableSASL, cfg.SaslMechanism, cfg.SaslUsername, cfg.SaslPassword); err != nil {
		return nil, err
	}

	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, kafkaConf)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// MustNewConsumerGroup is a convenience wrapper that panics on error.
// Deprecated: MustNewConsumerGroup panics on error. Use NewConsumerGroup instead.
func MustNewConsumerGroup(cfg ConsumerConfig) sarama.ConsumerGroup {
	group, err := NewConsumerGroup(cfg)
	if err != nil {
		log.Panicf("new consumer group: %v", err)
	}

	return group
}

// Processor is the function signature for handling a single Kafka message.
type Processor func(ctx context.Context, msg *sarama.ConsumerMessage) error

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler.
// It delegates message processing to the provided Processor and handles
// graceful shutdown via signalCtx.
type ConsumerGroupHandler struct {
	signalCtx context.Context
	processor Processor
}

// NewConsumerGroupHandler creates a ConsumerGroupHandler with the given processor.
// If processor is nil, DefaultProcessor is used.
func NewConsumerGroupHandler(signalCtx context.Context, processor Processor) *ConsumerGroupHandler {
	if processor == nil {
		processor = DefaultProcessor
	}
	return &ConsumerGroupHandler{
		signalCtx: signalCtx,
		processor: processor,
	}
}

func (*ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// ConsumeClaim is already called in a goroutine.
	for {
		// Strict check for graceful shutdown before pulling messages
		select {
		case <-h.signalCtx.Done():
			return nil
		default:
		}

		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			if err := h.processor(sess.Context(), msg); err != nil {
				slog.ErrorContext(sess.Context(), "failed to process message",
					slog.String("topic", msg.Topic),
					slog.Int("partition", int(msg.Partition)),
					slog.Int64("offset", msg.Offset),
					slog.String("error", err.Error()),
				)
			}
			sess.Commit()
		case <-sess.Context().Done():
			slog.Debug("consumer rebalance")
			return nil
		case <-h.signalCtx.Done():
			return nil
		}
	}
}

// DefaultProcessor logs the received message and returns nil.
func DefaultProcessor(_ context.Context, msg *sarama.ConsumerMessage) error {
	slog.Info("message received",
		slog.String("topic", msg.Topic),
		slog.Int("partition", int(msg.Partition)),
		slog.Int64("offset", msg.Offset),
	)
	return nil
}

// KafkaHandler is the function signature for handling a domain event.
// The message payload is kept as json.RawMessage so each handler can
// unmarshal it into its own domain-specific struct.
type KafkaHandler func(ctx context.Context, msg Message[json.RawMessage]) error

// NewEventRouter returns a Processor that routes incoming Kafka messages
// to the appropriate KafkaHandler based on the event name.
//
// It acts as an access-log interceptor for the async path — each handler
// invocation produces a single structured "event.completed" log with latency,
// metadata, and error details. Individual handlers should not log
// received/completed events themselves.
//
// Unknown events are logged at Debug level and silently skipped.
// Malformed messages that cannot be unmarshalled return an error.
func NewEventRouter(handlers map[string]KafkaHandler) Processor {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		start := time.Now()

		var envelope Message[json.RawMessage]
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			return fmt.Errorf("failed to unmarshal message envelope: %w", err)
		}

		handler, ok := handlers[envelope.EventName]
		if !ok {
			slog.DebugContext(ctx, "event.skipped",
				slog.String("event", "event.skipped"),
				slog.String("component", "event_router"),
				slog.String("event_name", envelope.EventName),
				slog.String("topic", msg.Topic),
			)
			return nil
		}

		err := handler(ctx, envelope)
		latencyMs := time.Since(start).Milliseconds()

		attrs := []any{
			slog.String("event", "event.completed"),
			slog.String("component", "event_router"),
			slog.String("event_name", envelope.EventName),
			slog.String("event_id", envelope.EventID),
			slog.String("aggregate_id", envelope.AggregateID),
			slog.String("topic", msg.Topic),
			slog.Int("partition", int(msg.Partition)),
			slog.Int64("offset", msg.Offset),
			slog.Int64("latency_ms", latencyMs),
		}

		if err != nil {
			msg, sourceAttrs := serror.DecodeMessage(err.Error())
			attrs = append(attrs, slog.String("error", msg))
			if len(sourceAttrs) > 0 {
				groupArgs := make([]any, len(sourceAttrs))
				for i, a := range sourceAttrs {
					groupArgs[i] = a
				}
				attrs = append(attrs, slog.Group("error_source", groupArgs...))
			}

			var sErr *serror.SError
			if errors.As(err, &sErr) {
				if ctxAttrs := sErr.Attrs(); len(ctxAttrs) > 0 {
					ctxArgs := make([]any, len(ctxAttrs))
					for i, a := range ctxAttrs {
						ctxArgs[i] = a
					}
					attrs = append(attrs, slog.Group("error_context", ctxArgs...))
				}
			}

			slog.ErrorContext(ctx, "event.completed", attrs...)
			return fmt.Errorf("handler failed for event %s: %w", envelope.EventName, err)
		}

		slog.InfoContext(ctx, "event.completed", attrs...)
		return nil
	}
}
