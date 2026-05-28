package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"
)

// LogMode controls how the interceptor logs Kafka messages.
type LogMode int

const (
	// LogModeDebug logs full payload + envelope metadata + custom attrs.
	// Suitable for LOCAL/DEV/SIT.
	LogModeDebug LogMode = iota
	// LogModeMeta logs envelope metadata + custom attrs only (no payload).
	// Suitable for UAT/PROD where payloads may contain PII.
	LogModeMeta
	// LogModeSilent suppresses all produce logs. Useful for high-frequency handlers.
	LogModeSilent
)

// LogInterceptorOption configures the logging interceptor.
type LogInterceptorOption struct {
	// Mode controls the verbosity of Kafka produce logs.
	Mode LogMode
}

// loggingProducer decorates a Producer with structured logging.
// It implements the Producer interface (Decorator / Interceptor pattern).
type loggingProducer struct {
	inner  Producer
	option LogInterceptorOption
}

var _ Producer = (*loggingProducer)(nil)

// WithLogging wraps a Producer with a logging interceptor.
//
// Usage:
//
//	producer := kafka.MustNewProducer(cfg)
//	logged   := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeDebug})
//
// In PROD, use LogModeMeta and attach investigation attrs via SendMessageOption.LogAttrs:
//
//	prod := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeMeta})
//	prod.SendMessageWithOption("topic", msg, kafka.SendMessageOption{
//	    LogAttrs: []slog.Attr{slog.String("organization_id", orgID)},
//	})
func WithLogging(p Producer, opt LogInterceptorOption) Producer {
	return &loggingProducer{inner: p, option: opt}
}

// WithLoggingFromEnv wraps a Producer and auto-selects the log mode from the
// ENV environment variable:
//
//	LOCAL, DEV → LogModeDebug
//	UAT, PROD  → LogModeMeta
func WithLoggingFromEnv(p Producer, env string) Producer {
	mode := logModeFromEnv(env)
	return WithLogging(p, LogInterceptorOption{Mode: mode})
}

func (lp *loggingProducer) Close() error {
	return lp.inner.Close()
}

func (lp *loggingProducer) SendMessage(ctx context.Context, topic string, message any) error {
	return lp.sendWithLog(ctx, topic, message, SendMessageOption{})
}

func (lp *loggingProducer) SendMessageWithOption(ctx context.Context, topic string, message any, option SendMessageOption) error {
	return lp.sendWithLog(ctx, topic, message, option)
}

func (lp *loggingProducer) sendWithLog(ctx context.Context, topic string, message any, option SendMessageOption) error {
	if lp.option.Mode == LogModeSilent {
		return lp.inner.SendMessageWithOption(ctx, topic, message, option)
	}

	start := time.Now()
	err := lp.inner.SendMessageWithOption(ctx, topic, message, option)
	elapsed := time.Since(start)

	// Determine log level before doing any serialization work.
	level := slog.LevelDebug
	if err != nil {
		level = slog.LevelError
	}

	// Level gate: skip all work when the log won't be emitted.
	if !slog.Default().Enabled(ctx, level) {
		return err
	}

	attrs := []slog.Attr{
		slog.String("topic", topic),
		slog.Duration("duration", elapsed),
	}

	if option.Key != "" {
		attrs = append(attrs, slog.String("key", option.Key))
	}

	// Attempt to extract event envelope metadata for structured logging.
	if meta := extractEnvelopeMeta(message); meta != nil {
		attrs = append(attrs,
			slog.String("event_name", meta.eventName),
			slog.String("aggregate_id", meta.aggregateID),
			slog.String("event_id", meta.eventID),
		)
	}

	// Caller-provided investigation attrs are always logged (all modes).
	attrs = append(attrs, option.LogAttrs...)

	// Only debug mode includes the full payload.
	if lp.option.Mode == LogModeDebug {
		payloadJSON, _ := json.Marshal(message)
		attrs = append(attrs, slog.String("payload", string(payloadJSON)))
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "kafka produce failed", attrs...)
		return err
	}

	slog.LogAttrs(ctx, slog.LevelDebug, "kafka produce ok", attrs...)
	return nil
}

// envelopeMeta holds the standard fields from kafka.Message.
type envelopeMeta struct {
	eventName   string
	aggregateID string
	eventID     string
}

// extractEnvelopeMeta tries to read event_name/aggregate_id/event_id from
// the message. Works for kafka.Message[T] or any struct with those JSON keys.
func extractEnvelopeMeta(msg any) *envelopeMeta {
	raw, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	var m struct {
		EventName   string `json:"event_name"`
		AggregateID string `json:"aggregate_id"`
		EventID     string `json:"event_id"`
	}
	if err := json.Unmarshal(raw, &m); err != nil || m.EventName == "" {
		return nil
	}
	return &envelopeMeta{
		eventName:   m.EventName,
		aggregateID: m.AggregateID,
		eventID:     m.EventID,
	}
}

func logModeFromEnv(env string) LogMode {
	switch strings.ToUpper(env) {
	case "LOCAL", "DEV":
		return LogModeDebug
	case "UAT", "PROD":
		return LogModeMeta
	default:
		return LogModeDebug
	}
}
