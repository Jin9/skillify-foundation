package firestore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Config.Validate
// ---------------------------------------------------------------------------

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

	t.Run("valid with emulator host", func(t *testing.T) {
		assert.NoError(t, Config{ProjectID: "p", EmulatorHost: "localhost:8080"}.Validate())
	})

	t.Run("valid with database ID", func(t *testing.T) {
		assert.NoError(t, Config{ProjectID: "p", DatabaseID: "my-db"}.Validate())
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

// ---------------------------------------------------------------------------
// shouldUseNamedDatabase
// ---------------------------------------------------------------------------

func TestShouldUseNamedDatabase(t *testing.T) {
	t.Run("empty string returns false", func(t *testing.T) {
		db, ok := shouldUseNamedDatabase("")
		assert.False(t, ok)
		assert.Equal(t, "", db)
	})

	t.Run("whitespace only returns false", func(t *testing.T) {
		db, ok := shouldUseNamedDatabase("   ")
		assert.False(t, ok)
		assert.Equal(t, "", db)
	})

	t.Run("default database ID returns false", func(t *testing.T) {
		db, ok := shouldUseNamedDatabase("(default)")
		assert.False(t, ok)
		assert.Equal(t, "", db)
	})

	t.Run("custom database returns true", func(t *testing.T) {
		db, ok := shouldUseNamedDatabase("my-database")
		assert.True(t, ok)
		assert.Equal(t, "my-database", db)
	})

	t.Run("custom database with whitespace is trimmed", func(t *testing.T) {
		db, ok := shouldUseNamedDatabase("  my-database  ")
		assert.True(t, ok)
		assert.Equal(t, "my-database", db)
	})
}

// ---------------------------------------------------------------------------
// withOptionalTimeout
// ---------------------------------------------------------------------------

func TestWithOptionalTimeout(t *testing.T) {
	t.Run("zero timeout returns context without deadline", func(t *testing.T) {
		ctx := context.Background()
		ctx2, cancel := withOptionalTimeout(ctx, 0)
		defer cancel()
		_, hasDeadline := ctx2.Deadline()
		assert.False(t, hasDeadline)
	})

	t.Run("negative timeout returns context without deadline", func(t *testing.T) {
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

// ---------------------------------------------------------------------------
// Client.Close (nil-guard paths)
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Client.Inner
// ---------------------------------------------------------------------------

func TestClient_Inner(t *testing.T) {
	c := &Client{inner: nil}
	assert.Nil(t, c.Inner())
}

// ---------------------------------------------------------------------------
// MustNewClient — panic path only
// (Successful creation requires Firestore emulator — integration test.)
// ---------------------------------------------------------------------------

func TestMustNewClient_PanicsOnValidationError(t *testing.T) {
	assert.Panics(t, func() {
		MustNewClient(context.Background(), Config{})
	})
}

// ---------------------------------------------------------------------------
// New — validation error path only
// (Successful creation requires Firestore emulator — integration test.)
// ---------------------------------------------------------------------------

func TestNew_ValidationError(t *testing.T) {
	_, err := New(context.Background(), Config{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ProjectID is required")
}

// ---------------------------------------------------------------------------
// NewClient — validation error path only
// ---------------------------------------------------------------------------

func TestNewClient_ValidationError(t *testing.T) {
	_, err := NewClient(context.Background(), Config{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ProjectID is required")
}
