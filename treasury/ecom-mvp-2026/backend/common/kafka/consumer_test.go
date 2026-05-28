package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestConsumerConfigConstructors(t *testing.T) {
	t.Run("no guarantee", func(t *testing.T) {
		config := NewConsumerConfigNoGuarantee()

		if config.Version != sarama.DefaultVersion {
			t.Fatalf("expected version %v, got %v", sarama.DefaultVersion, config.Version)
		}
	})

	t.Run("at most once", func(t *testing.T) {
		config := NewConsumerConfigAtMostOnce()

		if !config.Consumer.Offsets.AutoCommit.Enable {
			t.Fatalf("expected auto commit enabled")
		}

		if config.Consumer.Offsets.AutoCommit.Interval != time.Second {
			t.Fatalf("expected auto commit interval %v, got %v", time.Second, config.Consumer.Offsets.AutoCommit.Interval)
		}
	})

	t.Run("at least once", func(t *testing.T) {
		config := NewConsumerConfigAtLeastOnce()

		if config.Consumer.Offsets.AutoCommit.Enable {
			t.Fatalf("expected auto commit disabled")
		}
	})
}

func TestConsumerGroupRequiresConfig(t *testing.T) {
	_, err := NewConsumerGroup(ConsumerConfig{})
	if err == nil {
		t.Fatalf("expected error for nil kafka config")
	}
}

func TestMustNewConsumerGroupPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic from MustNewConsumerGroup")
		}
	}()

	_ = MustNewConsumerGroup(ConsumerConfig{})
}

func TestOffsetInitialFrom(t *testing.T) {
	if got := offsetInitialFrom("earliest"); got != sarama.OffsetOldest {
		t.Fatalf("expected oldest offset, got %d", got)
	}

	if got := offsetInitialFrom("latest"); got != sarama.OffsetNewest {
		t.Fatalf("expected newest offset, got %d", got)
	}

	if got := offsetInitialFrom("unknown"); got != sarama.OffsetNewest {
		t.Fatalf("expected default newest offset, got %d", got)
	}
}

func TestBalanceStrategyFrom(t *testing.T) {
	if got := fmt.Sprintf("%T", balanceStrategyFrom("range")); got == "<nil>" {
		t.Fatalf("expected range strategy")
	}

	if got, want := fmt.Sprintf("%T", balanceStrategyFrom("sticky")), fmt.Sprintf("%T", sarama.NewBalanceStrategySticky()); got != want {
		t.Fatalf("expected sticky strategy %s, got %s", want, got)
	}

	if got, want := fmt.Sprintf("%T", balanceStrategyFrom("roundrobin")), fmt.Sprintf("%T", sarama.NewBalanceStrategyRoundRobin()); got != want {
		t.Fatalf("expected round robin strategy %s, got %s", want, got)
	}

	if got, want := fmt.Sprintf("%T", balanceStrategyFrom("unknown")), fmt.Sprintf("%T", sarama.NewBalanceStrategyRange()); got != want {
		t.Fatalf("expected default range strategy %s, got %s", want, got)
	}
}

func TestNewEventRouter_DispatchesMatchingEvent(t *testing.T) {
	called := false
	handlers := map[string]KafkaHandler{
		"TEST_EVENT": func(_ context.Context, msg Message[json.RawMessage]) error {
			called = true
			if msg.EventName != "TEST_EVENT" {
				t.Fatalf("expected event TEST_EVENT, got %s", msg.EventName)
			}
			if msg.AggregateID != "agg-123" {
				t.Fatalf("expected aggregate_id agg-123, got %s", msg.AggregateID)
			}
			return nil
		},
	}

	processor := NewEventRouter(handlers)

	envelope := Message[json.RawMessage]{
		EventName:   "TEST_EVENT",
		AggregateID: "agg-123",
		EventID:     "evt-1",
		Payload:     json.RawMessage(`{"key":"value"}`),
	}
	body, _ := json.Marshal(envelope)

	msg := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    42,
		Value:     body,
	}

	err := processor(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("handler was not called")
	}
}

func TestNewEventRouter_UnknownEventSkipped(t *testing.T) {
	handlers := map[string]KafkaHandler{
		"REGISTERED_EVENT": func(_ context.Context, _ Message[json.RawMessage]) error {
			t.Fatal("handler should not be called for unknown event")
			return nil
		},
	}

	processor := NewEventRouter(handlers)

	envelope := Message[json.RawMessage]{
		EventName: "UNKNOWN_EVENT",
		EventID:   "evt-2",
	}
	body, _ := json.Marshal(envelope)

	msg := &sarama.ConsumerMessage{Value: body}
	if err := processor(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error for unknown event: %v", err)
	}
}

func TestNewEventRouter_MalformedMessage(t *testing.T) {
	processor := NewEventRouter(nil)

	msg := &sarama.ConsumerMessage{Value: []byte(`not json`)}
	err := processor(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error for malformed message")
	}
}

func TestNewEventRouter_HandlerError(t *testing.T) {
	handlerErr := errors.New("processing failed")
	handlers := map[string]KafkaHandler{
		"FAIL_EVENT": func(_ context.Context, _ Message[json.RawMessage]) error {
			return handlerErr
		},
	}

	processor := NewEventRouter(handlers)

	envelope := Message[json.RawMessage]{EventName: "FAIL_EVENT", EventID: "evt-3"}
	body, _ := json.Marshal(envelope)

	msg := &sarama.ConsumerMessage{Value: body}
	err := processor(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error from handler")
	}
	if !errors.Is(err, handlerErr) {
		t.Fatalf("expected wrapped handler error, got: %v", err)
	}
}

func TestNewEventRouter_PayloadPreserved(t *testing.T) {
	type testPayload struct {
		Name string `json:"name"`
	}

	handlers := map[string]KafkaHandler{
		"PAYLOAD_EVENT": func(_ context.Context, msg Message[json.RawMessage]) error {
			var p testPayload
			if err := json.Unmarshal(msg.Payload, &p); err != nil {
				t.Fatalf("failed to unmarshal payload: %v", err)
			}
			if p.Name != "test-name" {
				t.Fatalf("expected payload name test-name, got %s", p.Name)
			}
			return nil
		},
	}

	processor := NewEventRouter(handlers)

	payload, _ := json.Marshal(testPayload{Name: "test-name"})
	envelope := Message[json.RawMessage]{
		EventName: "PAYLOAD_EVENT",
		EventID:   "evt-4",
		Payload:   payload,
	}
	body, _ := json.Marshal(envelope)

	msg := &sarama.ConsumerMessage{Value: body}
	if err := processor(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDefaultConsumerGroup(t *testing.T) {
	config := NewDefaultConsumerGroup()
	assert.True(t, config.Consumer.Offsets.AutoCommit.Enable)
	assert.Equal(t, time.Second, config.Consumer.Offsets.AutoCommit.Interval)
}

func TestNewConsumerGroup_Error(t *testing.T) {
	cfg := ConsumerConfig{
		KafkaConf: nil,
	}
	_, err := NewConsumerGroup(cfg)
	assert.Error(t, err)

	cfg.KafkaConf = NewConsumerConfigNoGuarantee()
	cfg.EnableSSL = true
	cfg.CACertPEM = "invalid_pem"
	_, err = NewConsumerGroup(cfg)
	assert.Error(t, err) // Should error on TLS parsing
}

func TestNewConsumerGroupHandler(t *testing.T) {
	h := NewConsumerGroupHandler(context.Background(), nil)
	assert.NotNil(t, h)
	assert.NotNil(t, h.processor) // Should fallback to DefaultProcessor

	// Setup and Cleanup return nil
	assert.NoError(t, h.Setup(nil))
	assert.NoError(t, h.Cleanup(nil))
}

type mockSession struct {
	ctx context.Context
}

func (m *mockSession) Claims() map[string][]int32 { return nil }
func (m *mockSession) MemberID() string { return "" }
func (m *mockSession) GenerationID() int32 { return 0 }
func (m *mockSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {}
func (m *mockSession) Commit() {}
func (m *mockSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {}
func (m *mockSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {}
func (m *mockSession) Context() context.Context { return m.ctx }

type mockClaim struct {
	messages chan *sarama.ConsumerMessage
}

func (m *mockClaim) Topic() string { return "test-topic" }
func (m *mockClaim) Partition() int32 { return 0 }
func (m *mockClaim) InitialOffset() int64 { return 0 }
func (m *mockClaim) HighWaterMarkOffset() int64 { return 0 }
func (m *mockClaim) Messages() <-chan *sarama.ConsumerMessage { return m.messages }

func TestConsumeClaim(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sess := &mockSession{ctx: context.Background()}
	claim := &mockClaim{messages: make(chan *sarama.ConsumerMessage, 2)}

	h := NewConsumerGroupHandler(ctx, nil)

	// Test normal message processing
	claim.messages <- &sarama.ConsumerMessage{Topic: "test-topic", Partition: 0, Offset: 1}
	close(claim.messages) // Will cause claim.Messages() to return !ok

	err := h.ConsumeClaim(sess, claim)
	assert.NoError(t, err)

	// Test cancellation
	claim2 := &mockClaim{messages: make(chan *sarama.ConsumerMessage)}
	cancel() // Cancel the signal context
	err = h.ConsumeClaim(sess, claim2)
	assert.NoError(t, err)
}

func TestDefaultProcessor(t *testing.T) {
	err := DefaultProcessor(context.Background(), &sarama.ConsumerMessage{})
	assert.NoError(t, err)
}
