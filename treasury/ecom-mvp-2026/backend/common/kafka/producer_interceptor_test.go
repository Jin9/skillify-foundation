package kafka

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

type mockProducer struct {
	sendErr    error
	sentTopic  string
	sentMsg    any
	sentOption SendMessageOption
	closed     bool
}

func (m *mockProducer) Close() error { m.closed = true; return nil }
func (m *mockProducer) SendMessage(ctx context.Context, topic string, message any) error {
	return m.SendMessageWithOption(ctx, topic, message, SendMessageOption{})
}
func (m *mockProducer) SendMessageWithOption(ctx context.Context, topic string, message any, option SendMessageOption) error {
	m.sentTopic = topic
	m.sentMsg = message
	m.sentOption = option
	return m.sendErr
}

func captureLog(t *testing.T, level slog.Level) *bytes.Buffer {
	t.Helper()
	buf := new(bytes.Buffer)
	prev := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prev) })
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: level})))
	return buf
}

func TestLoggingProducer_DebugMode(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{}

	msg := NewMessage("TEST_EVENT", "agg-1", map[string]string{"key": "value"})
	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeDebug})

	err := p.SendMessage(context.Background(), "my-topic", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log := buf.String()
	if !strings.Contains(log, "kafka produce ok") {
		t.Errorf("expected 'kafka produce ok', got:\n%s", log)
	}
	if !strings.Contains(log, "my-topic") {
		t.Errorf("expected topic in log, got:\n%s", log)
	}
	if !strings.Contains(log, "TEST_EVENT") {
		t.Errorf("expected event_name in log, got:\n%s", log)
	}
	if !strings.Contains(log, "agg-1") {
		t.Errorf("expected aggregate_id in log, got:\n%s", log)
	}
	if !strings.Contains(log, `key`) || !strings.Contains(log, `value`) {
		t.Errorf("expected full payload in debug log, got:\n%s", log)
	}
}

func TestLoggingProducer_MetaMode_NoPayload(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{}

	msg := NewMessage("USER_REGISTER", "u-1", map[string]string{
		"password": "s3cret!",
		"email":    "test@example.com",
	})
	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeMeta})

	if err := p.SendMessage(context.Background(), "user-topic", msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log := buf.String()

	// Envelope metadata should be visible
	if !strings.Contains(log, "USER_REGISTER") {
		t.Errorf("expected event_name visible, got:\n%s", log)
	}
	if !strings.Contains(log, "user-topic") {
		t.Errorf("expected topic visible, got:\n%s", log)
	}

	// Payload should NOT be logged at all
	if strings.Contains(log, "s3cret!") {
		t.Errorf("leaked password in meta mode, got:\n%s", log)
	}
	if strings.Contains(log, "test@example.com") {
		t.Errorf("leaked email in meta mode, got:\n%s", log)
	}
}

func TestLoggingProducer_MetaMode_WithLogAttrs(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{}

	msg := NewMessage("ORDER_CREATED", "o-1", map[string]string{
		"card_number": "4111111111111111",
	})
	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeMeta})

	err := p.SendMessageWithOption(context.Background(), "order-topic", msg, SendMessageOption{
		LogAttrs: []slog.Attr{
			slog.String("organization_id", "org-456"),
			slog.Int("item_count", 3),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log := buf.String()

	// Custom attrs should be visible
	if !strings.Contains(log, "org-456") {
		t.Errorf("expected organization_id in log, got:\n%s", log)
	}
	if !strings.Contains(log, "3") {
		t.Errorf("expected item_count in log, got:\n%s", log)
	}

	// Payload should NOT be logged
	if strings.Contains(log, "4111111111111111") {
		t.Errorf("leaked card_number in meta mode, got:\n%s", log)
	}
}

func TestLoggingProducer_DebugMode_WithLogAttrs(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{}

	msg := NewMessage("TEST", "agg-1", "data")
	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeDebug})

	err := p.SendMessageWithOption(context.Background(), "topic", msg, SendMessageOption{
		LogAttrs: []slog.Attr{
			slog.String("trace_id", "abc-123"),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log := buf.String()

	// Both custom attrs AND payload should be visible in debug mode
	if !strings.Contains(log, "abc-123") {
		t.Errorf("expected trace_id in log, got:\n%s", log)
	}
	if !strings.Contains(log, "payload") {
		t.Errorf("expected payload in debug log, got:\n%s", log)
	}
}

func TestLoggingProducer_SilentMode(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{}

	msg := NewMessage("BATCH_JOB", "batch-1", "data")
	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeSilent})

	err := p.SendMessage(context.Background(), "batch-topic", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() > 0 {
		t.Errorf("expected no log output in silent mode, got:\n%s", buf.String())
	}
	if inner.sentTopic != "batch-topic" {
		t.Errorf("expected message to still be sent, topic: %s", inner.sentTopic)
	}
}

func TestLoggingProducer_ErrorLogsAtErrorLevel(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{sendErr: errors.New("broker unavailable")}

	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeDebug})

	err := p.SendMessage(context.Background(), "fail-topic", "data")
	if err == nil {
		t.Fatal("expected error")
	}

	log := buf.String()
	if !strings.Contains(log, "kafka produce failed") {
		t.Errorf("expected 'kafka produce failed', got:\n%s", log)
	}
	if !strings.Contains(log, "broker unavailable") {
		t.Errorf("expected error message in log, got:\n%s", log)
	}
}

func TestLoggingProducer_Close(t *testing.T) {
	inner := &mockProducer{}
	p := WithLogging(inner, LogInterceptorOption{})

	if err := p.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !inner.closed {
		t.Error("expected inner producer to be closed")
	}
}

func TestLoggingProducer_WithKey(t *testing.T) {
	buf := captureLog(t, slog.LevelDebug)
	inner := &mockProducer{}

	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeDebug})
	err := p.SendMessageWithOption(context.Background(), "topic", "msg", SendMessageOption{Key: "partition-key"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log := buf.String()
	if !strings.Contains(log, "partition-key") {
		t.Errorf("expected key in log, got:\n%s", log)
	}
}

func TestLoggingProducer_LevelGate_SkipsWork(t *testing.T) {
	buf := captureLog(t, slog.LevelInfo)
	inner := &mockProducer{}

	msg := NewMessage("TEST_EVENT", "agg-1", map[string]string{"key": "value"})
	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeDebug})

	err := p.SendMessage(context.Background(), "topic", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Success logs at DEBUG, but slog level is INFO → no output
	if buf.Len() > 0 {
		t.Errorf("expected no log when level is below threshold, got:\n%s", buf.String())
	}
}

func TestLoggingProducer_LevelGate_ErrorStillLogs(t *testing.T) {
	buf := captureLog(t, slog.LevelInfo)
	inner := &mockProducer{sendErr: errors.New("timeout")}

	p := WithLogging(inner, LogInterceptorOption{Mode: LogModeDebug})
	_ = p.SendMessage(context.Background(), "topic", "data")

	log := buf.String()
	if !strings.Contains(log, "kafka produce failed") {
		t.Errorf("expected error to still be logged at INFO level, got:\n%s", log)
	}
}

func TestWithLoggingFromEnv(t *testing.T) {
	tests := []struct {
		env      string
		wantMode LogMode
	}{
		{"LOCAL", LogModeDebug},
		{"DEV", LogModeDebug},
		{"local", LogModeDebug},
		{"UAT", LogModeMeta},
		{"PROD", LogModeMeta},
		{"prod", LogModeMeta},
		{"", LogModeDebug},
	}

	for _, tt := range tests {
		t.Run("env="+tt.env, func(t *testing.T) {
			inner := &mockProducer{}
			p := WithLoggingFromEnv(inner, tt.env)

			lp, ok := p.(*loggingProducer)
			if !ok {
				t.Fatal("expected *loggingProducer")
			}
			if lp.option.Mode != tt.wantMode {
				t.Errorf("env %q: expected mode %d, got %d", tt.env, tt.wantMode, lp.option.Mode)
			}
		})
	}
}
