package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/cenkalti/backoff/v4"
)

type Producer interface {
	Close() error
	SendMessage(ctx context.Context, topic string, message any) error
	SendMessageWithOption(ctx context.Context, topic string, message any, option SendMessageOption) error
}

type ProducerConfig struct {
	KafkaConf     *sarama.Config
	Brokers       []string
	EnableSSL     bool
	CACertPEM     string
	ClientCertPEM string
	ClientKeyPEM  string
	EnableSASL    bool
	SaslMechanism sarama.SASLMechanism
	SaslUsername  string
	SaslPassword  string
}

type producer struct {
	client sarama.SyncProducer
}

var _ Producer = (*producer)(nil)

func NewSyncProducerGuarantee() *sarama.Config {
	config := newConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = true

	return config
}

func NewSyncProducerFireAndForget() *sarama.Config {
	config := newConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = true

	return config
}

func NewProducer(cfg ProducerConfig) (Producer, error) {

	kafkaConf := cfg.KafkaConf
	if kafkaConf == nil {
		return nil, fmt.Errorf("kafka config must not be nil")
	}
	if err := applyTLSConfig(kafkaConf, cfg.EnableSSL, cfg.CACertPEM, cfg.ClientCertPEM, cfg.ClientKeyPEM); err != nil {
		return nil, err
	}
	if err := applySASLConfig(kafkaConf, cfg.EnableSASL, cfg.SaslMechanism, cfg.SaslUsername, cfg.SaslPassword); err != nil {
		return nil, err
	}

	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, kafkaConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka sync producer: %w", err)
	}

	producer := &producer{
		client: syncProducer,
	}

	return producer, nil
}

// Deprecated: MustNewProducer panics on error. Use NewProducer instead.
func MustNewProducer(cfg ProducerConfig) Producer {
	producer, err := NewProducer(cfg)
	if err != nil {
		log.Panicf("new producer: %v", err)
	}
	return producer
}

func (p *producer) Close() error {
	return p.client.Close()
}

type SendMessageOption struct {
	WithRetry   bool
	FailedTopic string
	MaxRetries  int
	Header      []sarama.RecordHeader
	Key         string

	// LogAttrs are caller-provided structured attributes that the logging
	// interceptor includes in every log line (all modes except Silent).
	// Use this to attach investigation-safe values (IDs, counts, status codes)
	// without exposing the full payload.
	LogAttrs []slog.Attr
}

func (p *producer) SendMessage(ctx context.Context, topic string, message any) error {
	return p.SendMessageWithOption(ctx, topic, message, SendMessageOption{})
}

func (p *producer) SendMessageWithOption(ctx context.Context, topic string, message any, option SendMessageOption) error {
	msgBytes, err := encodeMessage(message)
	if err != nil {
		return err
	}

	if option.WithRetry {
		if option.FailedTopic == "" {
			return fmt.Errorf("failed topic is required")
		}
		return p.sendMessageWithRetry(topic, msgBytes, option)
	}
	return p.sendMessage(topic, msgBytes, option)
}

func (p *producer) sendMessage(topic string, message []byte, option SendMessageOption) error {
	msg := buildProducerMessage(topic, message, option)

	if _, _, err := p.client.SendMessage(msg); err != nil {
		return fmt.Errorf("failed to send sync message: %w", err)
	}

	return nil
}

func (p *producer) sendMessageWithRetry(topic string, message []byte, option SendMessageOption) error {
	msg := buildProducerMessage(topic, message, option)

	// Retries Step
	err := backoff.RetryNotify(
		func() error {
			_, _, err := p.client.SendMessage(msg)
			if err != nil {
				return fmt.Errorf("failed to send sync message: %w", err)
			}
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(option.MaxRetries)),
		func(err error, duration time.Duration) {
			log.Printf("retrying to send message in %s: %v", duration, err)
		},
	)

	// Failed step, When reached max retries
	if err != nil {
		msg.Topic = option.FailedTopic
		_, _, retryErr := p.client.SendMessage(msg)
		if retryErr != nil {
			return fmt.Errorf("failed to send message to failed topic: %w (original error: %v)", retryErr, err)
		}

		return fmt.Errorf("message sent to failed topic successfully (original error: %w)", err)
	}

	return nil
}

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func encodeMessage(message any) ([]byte, error) {
	switch v := message.(type) {
	case []byte:
		return v, nil
	default:
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)

		if err := json.NewEncoder(buf).Encode(message); err != nil {
			return nil, fmt.Errorf("failed to encode message to JSON: %w", err)
		}
		// Copy the buffer bytes because Put will recycle it
		dst := make([]byte, buf.Len())
		copy(dst, buf.Bytes())
		return dst, nil
	}
}

func buildProducerMessage(topic string, message []byte, option SendMessageOption) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.ByteEncoder(message),
		Headers: option.Header,
	}

	if option.Key != "" {
		msg.Key = sarama.StringEncoder(option.Key)
	}

	return msg
}
