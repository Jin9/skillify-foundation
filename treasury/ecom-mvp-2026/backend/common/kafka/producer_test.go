package kafka

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	saramaMocks "github.com/IBM/sarama/mocks"
)

func TestProducerConfigConstructors(t *testing.T) {
	t.Run("guarantee config", func(t *testing.T) {
		config := NewSyncProducerGuarantee()

		if config.Version != sarama.DefaultVersion {
			t.Fatalf("expected version %v, got %v", sarama.DefaultVersion, config.Version)
		}

		if config.Producer.RequiredAcks != sarama.WaitForAll {
			t.Fatalf("expected required acks %v, got %v", sarama.WaitForAll, config.Producer.RequiredAcks)
		}
	})

	t.Run("fire and forget config", func(t *testing.T) {
		config := NewSyncProducerFireAndForget()

		if config.Version != sarama.DefaultVersion {
			t.Fatalf("expected version %v, got %v", sarama.DefaultVersion, config.Version)
		}

		if config.Producer.RequiredAcks != sarama.NoResponse {
			t.Fatalf("expected required acks %v, got %v", sarama.NoResponse, config.Producer.RequiredAcks)
		}
	})
}

func TestProducerRequiresConfig(t *testing.T) {
	_, err := NewProducer(ProducerConfig{})
	if err == nil {
		t.Fatalf("expected error for nil kafka config")
	}
}

func TestMustNewProducerPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic from MustNewProducer")
		}
	}()

	_ = MustNewProducer(ProducerConfig{})
}

func TestProducerSendMessage(t *testing.T) {
	t.Run("send message success", func(t *testing.T) {
		mockProducer := saramaMocks.NewSyncProducer(t, NewSyncProducerGuarantee())
		defer func() { _ = mockProducer.Close() }()

		mockProducer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
			if msg.Topic != "topic-a" {
				return fmt.Errorf("unexpected topic %q", msg.Topic)
			}

			encoded, err := msg.Value.Encode()
			if err != nil {
				return err
			}

			if string(encoded) != "payload" {
				return fmt.Errorf("unexpected payload %q", string(encoded))
			}

			return nil
		})

		p := &producer{client: mockProducer}
		if err := p.SendMessage(context.Background(), "topic-a", []byte("payload")); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("send message encode error", func(t *testing.T) {
		p := &producer{}

		err := p.SendMessage(context.Background(), "topic-a", func() {})
		if err == nil {
			t.Fatalf("expected encode error")
		}
	})

	t.Run("send message producer error", func(t *testing.T) {
		mockProducer := saramaMocks.NewSyncProducer(t, NewSyncProducerGuarantee())
		defer func() { _ = mockProducer.Close() }()

		mockProducer.ExpectSendMessageAndFail(errors.New("send failed"))

		p := &producer{client: mockProducer}
		err := p.SendMessage(context.Background(), "topic-a", []byte("payload"))
		if err == nil {
			t.Fatalf("expected send error")
		}
	})
}

func TestProducerSendMessageWithOption(t *testing.T) {
	t.Run("with retry requires failed topic", func(t *testing.T) {
		p := &producer{}

		err := p.SendMessageWithOption(context.Background(), "topic-a", []byte("payload"), SendMessageOption{WithRetry: true})
		if err == nil {
			t.Fatalf("expected error when failed topic is missing")
		}
	})

	t.Run("with retry success on first attempt", func(t *testing.T) {
		mockProducer := saramaMocks.NewSyncProducer(t, NewSyncProducerGuarantee())
		defer func() { _ = mockProducer.Close() }()

		mockProducer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
			if msg.Topic != "topic-a" {
				return fmt.Errorf("unexpected topic %q", msg.Topic)
			}
			return nil
		})

		p := &producer{client: mockProducer}
		err := p.SendMessageWithOption(context.Background(), "topic-a", []byte("payload"), SendMessageOption{
			WithRetry:   true,
			FailedTopic: "topic-failed",
			MaxRetries:  0,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("with retry falls back to failed topic", func(t *testing.T) {
		mockProducer := saramaMocks.NewSyncProducer(t, NewSyncProducerGuarantee())
		defer func() { _ = mockProducer.Close() }()

		mockProducer.ExpectSendMessageWithMessageCheckerFunctionAndFail(func(msg *sarama.ProducerMessage) error {
			if msg.Topic != "topic-a" {
				return fmt.Errorf("unexpected original topic %q", msg.Topic)
			}
			return nil
		}, errors.New("send failed"))

		mockProducer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
			if msg.Topic != "topic-failed" {
				return fmt.Errorf("unexpected failed topic %q", msg.Topic)
			}
			return nil
		})

		p := &producer{client: mockProducer}
		err := p.SendMessageWithOption(context.Background(), "topic-a", []byte("payload"), SendMessageOption{
			WithRetry:   true,
			FailedTopic: "topic-failed",
			MaxRetries:  0,
		})
		if err == nil {
			t.Fatalf("expected wrapped fallback error")
		}
	})

	t.Run("with retry failed topic send error", func(t *testing.T) {
		mockProducer := saramaMocks.NewSyncProducer(t, NewSyncProducerGuarantee())
		defer func() { _ = mockProducer.Close() }()

		mockProducer.ExpectSendMessageAndFail(errors.New("send failed"))
		mockProducer.ExpectSendMessageAndFail(errors.New("failed topic send failed"))

		p := &producer{client: mockProducer}
		err := p.SendMessageWithOption(context.Background(), "topic-a", []byte("payload"), SendMessageOption{
			WithRetry:   true,
			FailedTopic: "topic-failed",
			MaxRetries:  0,
		})
		if err == nil {
			t.Fatalf("expected failed topic send error")
		}
	})
}

func TestBuildProducerMessage(t *testing.T) {
	headers := []sarama.RecordHeader{{Key: []byte("x-id"), Value: []byte("123")}}
	msg := buildProducerMessage("topic-a", []byte("payload"), SendMessageOption{
		Header: headers,
		Key:    "key-1",
	})

	if msg.Topic != "topic-a" {
		t.Fatalf("expected topic topic-a, got %q", msg.Topic)
	}

	if len(msg.Headers) != 1 {
		t.Fatalf("expected one header, got %d", len(msg.Headers))
	}

	key, err := msg.Key.Encode()
	if err != nil {
		t.Fatalf("encode key: %v", err)
	}
	if string(key) != "key-1" {
		t.Fatalf("expected key key-1, got %q", string(key))
	}

	value, err := msg.Value.Encode()
	if err != nil {
		t.Fatalf("encode value: %v", err)
	}
	if string(value) != "payload" {
		t.Fatalf("expected payload, got %q", string(value))
	}
}

func TestEncodeMessage(t *testing.T) {
	t.Run("bytes", func(t *testing.T) {
		msg, err := encodeMessage([]byte("payload"))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(msg) != "payload" {
			t.Fatalf("expected payload, got %q", string(msg))
		}
	})

	t.Run("json", func(t *testing.T) {
		msg, err := encodeMessage(map[string]string{"hello": "world"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(msg) != "{\"hello\":\"world\"}\n" {
			t.Fatalf("unexpected json payload %q", string(msg))
		}
	})

	t.Run("json error", func(t *testing.T) {
		_, err := encodeMessage(make(chan int))
		if err == nil {
			t.Fatalf("expected json encode error")
		}
	})
}

func TestProducerClose(t *testing.T) {
	mockProducer := saramaMocks.NewSyncProducer(t, NewSyncProducerGuarantee())
	p := &producer{client: mockProducer}

	if err := p.Close(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewProducer_Error(t *testing.T) {
	cfg := ProducerConfig{KafkaConf: nil}
	_, err := NewProducer(cfg)
	if err == nil {
		t.Fatal("expected error on nil kafka config")
	}

	cfg.KafkaConf = NewSyncProducerGuarantee()
	cfg.EnableSSL = true
	cfg.CACertPEM = "invalid_pem"
	_, err = NewProducer(cfg)
	if err == nil {
		t.Fatal("expected error on tls cert parsing")
	}
}
