package firestore

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	gcpfirestore "cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// Config configures a Firestore client.
//
// If EmulatorHost is provided, this package will set FIRESTORE_EMULATOR_HOST
// (unless already set) and will include option.WithoutAuthentication().
//
// CredentialsFile / CredentialsJSON are mutually exclusive.
//
// ConnectTimeout applies only to New/NewClient creation.
// If zero, the passed context controls the lifetime.
type Config struct {
	ProjectID string

	// EmulatorHost example: "localhost:8080".
	EmulatorHost string

	CredentialsFile string
	CredentialsJSON []byte

	// DatabaseID selects a named Firestore database.
	// Leave empty (or set to "(default)") to use the default database.
	DatabaseID string

	ConnectTimeout time.Duration
}

func (cfg Config) Validate() error {
	if cfg.ProjectID == "" {
		return errors.New("firestore: ProjectID is required")
	}
	if cfg.CredentialsFile != "" && len(cfg.CredentialsJSON) > 0 {
		return errors.New("firestore: CredentialsFile and CredentialsJSON are mutually exclusive")
	}
	return nil
}

func shouldUseNamedDatabase(databaseID string) (string, bool) {
	db := strings.TrimSpace(databaseID)
	if db == "" || db == gcpfirestore.DefaultDatabaseID {
		return "", false
	}
	return db, true
}

// New creates a Firestore client.
func New(ctx context.Context, cfg Config, clientOptions ...option.ClientOption) (*gcpfirestore.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	ctx2, cancel := withOptionalTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	opts := make([]option.ClientOption, 0, len(clientOptions)+2)

	if cfg.EmulatorHost != "" {
		if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
			_ = os.Setenv("FIRESTORE_EMULATOR_HOST", cfg.EmulatorHost)
		}
		// For emulator usage, ADC isn't required.
		opts = append(opts, option.WithoutAuthentication())
	}

	if len(cfg.CredentialsJSON) > 0 {
		opts = append(opts, option.WithCredentialsJSON(cfg.CredentialsJSON))
	} else if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}

	opts = append(opts, clientOptions...)

	var (
		client *gcpfirestore.Client
		err    error
	)
	if db, ok := shouldUseNamedDatabase(cfg.DatabaseID); ok {
		client, err = gcpfirestore.NewClientWithDatabase(ctx2, cfg.ProjectID, db, opts...)
	} else {
		client, err = gcpfirestore.NewClient(ctx2, cfg.ProjectID, opts...)
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Client is a thin wrapper around the official Firestore client.
type Client struct {
	inner *gcpfirestore.Client
}

// NewClient creates a wrapper client.
func NewClient(ctx context.Context, cfg Config, clientOptions ...option.ClientOption) (*Client, error) {
	c, err := New(ctx, cfg, clientOptions...)
	if err != nil {
		return nil, err
	}
	return &Client{inner: c}, nil
}

// MustNewClient is a convenience wrapper that panics on error.
func MustNewClient(ctx context.Context, cfg Config, clientOptions ...option.ClientOption) *Client {
	c, err := NewClient(ctx, cfg, clientOptions...)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Client) Inner() *gcpfirestore.Client {
	return c.inner
}

func (c *Client) Close() error {
	if c == nil || c.inner == nil {
		return nil
	}
	return c.inner.Close()
}

func withOptionalTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}
