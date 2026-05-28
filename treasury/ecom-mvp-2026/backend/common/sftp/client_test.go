package sftp

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_InvalidAuth(t *testing.T) {
	_, err := NewClient(Config{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no password or private key provided")
}

func TestNewClient_InvalidPrivateKey(t *testing.T) {
	_, err := NewClient(Config{
		PrivateKeyAuth: []byte("invalid-private-key"),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse private key")
}

func TestNewClient_DialFailure(t *testing.T) {
	cfg := Config{
		Host:     "invalid.host.local",
		Port:     "22",
		Username: "user",
		Password: "password",
		Timeout:  time.Millisecond,
	}

	_, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to dial ssh")
}

func TestClient_ContextCancellations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	c := &client{}

	// Upload
	err := c.Upload(ctx, "/path", strings.NewReader("data"))
	assert.ErrorIs(t, err, context.Canceled)

	// Download
	err = c.Download(ctx, "/path", &strings.Builder{})
	assert.ErrorIs(t, err, context.Canceled)

	// ListDirectory
	_, err = c.ListDirectory(ctx, "/path")
	assert.ErrorIs(t, err, context.Canceled)

	// Delete
	err = c.Delete(ctx, "/path")
	assert.ErrorIs(t, err, context.Canceled)
}
