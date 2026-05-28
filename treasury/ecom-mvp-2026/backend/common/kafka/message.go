package kafka

import (
	"time"

	"github.com/google/uuid"
)

// Message is the standard Kafka event envelope.
//
// Every event produced or consumed through Kafka should use this structure.
// This aligns with the DDD event-driven pattern and ensures:
//   - Consistent event naming convention: <DOMAIN>_<ACTION>_<STAGE>
//   - Idempotency via EventID (unique per event occurrence)
//   - Correlation via AggregateID (links events to the same entity)
//
// Example:
//
//	msg := kafka.NewMessage("MEMBER_REGISTER_SUCCESS", memberID, payload)
type Message[T any] struct {
	EventName   string    `json:"event_name"`
	AggregateID string    `json:"aggregate_id"`
	EventID     string    `json:"event_id"`
	Timestamp   time.Time `json:"timestamp"`
	Payload     T         `json:"payload"`
}

// NewMessage creates a Message with an auto-generated EventID and current timestamp.
func NewMessage[T any](eventName, aggregateID string, payload T) Message[T] {
	return Message[T]{
		EventName:   eventName,
		AggregateID: aggregateID,
		EventID:     uuid.New().String(),
		Timestamp:   time.Now().UTC(),
		Payload:     payload,
	}
}

// NewMessageWithID creates a Message with a caller-supplied EventID.
// Useful for replaying or forwarding events with the same idempotency key.
func NewMessageWithID[T any](eventName, aggregateID, eventID string, payload T) Message[T] {
	return Message[T]{
		EventName:   eventName,
		AggregateID: aggregateID,
		EventID:     eventID,
		Timestamp:   time.Now().UTC(),
		Payload:     payload,
	}
}
