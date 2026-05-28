package cloudstorage

import (
	"context"
	"errors"
	"os"
	"time"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Config configures a Cloud Storage client.
//
// If EmulatorHost is provided, this package will set STORAGE_EMULATOR_HOST
// (unless already set) and will include option.WithoutAuthentication().
//
// CredentialsFile / CredentialsJSON are mutually exclusive.
//
// ConnectTimeout applies only to New/NewClient creation.
// If zero, the passed context controls the lifetime.
type Config struct {
	ProjectID string

	// EmulatorHost example: "http://localhost:4443".
	// For the official GCS emulator, this is typically a full URL with scheme.
	EmulatorHost string

	CredentialsFile string
	CredentialsJSON []byte

	ConnectTimeout time.Duration
}

func (cfg Config) Validate() error {
	if cfg.ProjectID == "" {
		return errors.New("cloudstorage: ProjectID is required")
	}
	if cfg.CredentialsFile != "" && len(cfg.CredentialsJSON) > 0 {
		return errors.New("cloudstorage: CredentialsFile and CredentialsJSON are mutually exclusive")
	}
	return nil
}

// New creates a Cloud Storage client.
func New(ctx context.Context, cfg Config, clientOptions ...option.ClientOption) (*gcs.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	ctx2, cancel := withOptionalTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	opts := make([]option.ClientOption, 0, len(clientOptions)+2)

	if cfg.EmulatorHost != "" {
		// Cloud Storage client respects STORAGE_EMULATOR_HOST.
		// Use full URL (including scheme) when applicable.
		if os.Getenv("STORAGE_EMULATOR_HOST") == "" {
			_ = os.Setenv("STORAGE_EMULATOR_HOST", cfg.EmulatorHost)
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

	client, err := gcs.NewClient(ctx2, opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Client is a thin wrapper around the official Cloud Storage client.
type Client struct {
	inner *gcs.Client
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

func (c *Client) Inner() *gcs.Client {
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
