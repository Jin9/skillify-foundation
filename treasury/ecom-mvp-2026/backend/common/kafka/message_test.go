package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage("TEST_EVENT", "agg-1", map[string]string{"key": "value"})

	assert.Equal(t, "TEST_EVENT", msg.EventName)
	assert.Equal(t, "agg-1", msg.AggregateID)
	assert.NotEmpty(t, msg.EventID)
	assert.False(t, msg.Timestamp.IsZero())
	assert.Equal(t, "value", msg.Payload["key"])
}

func TestNewMessageWithID(t *testing.T) {
	msg := NewMessageWithID("TEST_EVENT", "agg-1", "evt-fixed", "payload")

	assert.Equal(t, "TEST_EVENT", msg.EventName)
	assert.Equal(t, "agg-1", msg.AggregateID)
	assert.Equal(t, "evt-fixed", msg.EventID)
	assert.False(t, msg.Timestamp.IsZero())
	assert.Equal(t, "payload", msg.Payload)
}
