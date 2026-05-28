package cloudstorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	t.Run("valid with project ID only", func(t *testing.T) {
		assert.NoError(t, Config{ProjectID: "my-project"}.Validate())
	})

	t.Run("valid with credentials file", func(t *testing.T) {
		assert.NoError(t, Config{ProjectID: "p", CredentialsFile: "/key.json"}.Validate())
	})

	t.Run("valid with credentials JSON", func(t *testing.T) {
		assert.NoError(t, Config{ProjectID: "p", CredentialsJSON: []byte(`{}`)}.Validate())
	})

	t.Run("error when project ID is empty", func(t *testing.T) {
		err := Config{}.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ProjectID is required")
	})

	t.Run("error when both credentials are set", func(t *testing.T) {
		err := Config{
			ProjectID:       "p",
			CredentialsFile: "/key.json",
			CredentialsJSON: []byte(`{}`),
		}.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mutually exclusive")
	})
}

func TestWithOptionalTimeout(t *testing.T) {
	t.Run("no timeout returns same context", func(t *testing.T) {
		ctx := context.Background()
		ctx2, cancel := withOptionalTimeout(ctx, 0)
		defer cancel()
		_, hasDeadline := ctx2.Deadline()
		assert.False(t, hasDeadline)
	})

	t.Run("negative timeout returns same context", func(t *testing.T) {
		ctx := context.Background()
		ctx2, cancel := withOptionalTimeout(ctx, -1*time.Second)
		defer cancel()
		_, hasDeadline := ctx2.Deadline()
		assert.False(t, hasDeadline)
	})

	t.Run("positive timeout returns context with deadline", func(t *testing.T) {
		ctx := context.Background()
		ctx2, cancel := withOptionalTimeout(ctx, 5*time.Second)
		defer cancel()
		_, hasDeadline := ctx2.Deadline()
		assert.True(t, hasDeadline)
	})
}

func TestClient_Close(t *testing.T) {
	t.Run("nil client", func(t *testing.T) {
		var c *Client
		assert.NoError(t, c.Close())
	})

	t.Run("nil inner", func(t *testing.T) {
		c := &Client{inner: nil}
		assert.NoError(t, c.Close())
	})
}

func TestClient_Inner(t *testing.T) {
	c := &Client{inner: nil}
	assert.Nil(t, c.Inner())
}

func TestMustNewClient_PanicsOnValidationError(t *testing.T) {
	assert.Panics(t, func() {
		MustNewClient(context.Background(), Config{})
	})
}

func TestNew_ValidationError(t *testing.T) {
	_, err := New(context.Background(), Config{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ProjectID is required")
}

func TestNewClient_ValidationError(t *testing.T) {
	_, err := NewClient(context.Background(), Config{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ProjectID is required")
}
