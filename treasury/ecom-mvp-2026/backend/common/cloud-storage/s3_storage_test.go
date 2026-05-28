package cloudstorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3Config_Validate(t *testing.T) {
	t.Run("valid with region only", func(t *testing.T) {
		assert.NoError(t, S3Config{Region: "us-east-1"}.Validate())
	})

	t.Run("valid with full static credentials", func(t *testing.T) {
		assert.NoError(t, S3Config{
			Region:          "us-east-1",
			AccessKeyID:     "AKID",
			SecretAccessKey: "secret",
		}.Validate())
	})

	t.Run("error when region is empty", func(t *testing.T) {
		err := S3Config{}.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Region is required")
	})

	t.Run("error when only access key provided", func(t *testing.T) {
		err := S3Config{Region: "us-east-1", AccessKeyID: "AKID"}.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be provided together")
	})

	t.Run("error when only secret key provided", func(t *testing.T) {
		err := S3Config{Region: "us-east-1", SecretAccessKey: "secret"}.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be provided together")
	})
}

func TestNewS3(t *testing.T) {
	t.Run("validation error propagated", func(t *testing.T) {
		_, err := NewS3(context.Background(), S3Config{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Region is required")
	})

	t.Run("aws config load error propagated", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{}, errors.New("aws load boom")
		}

		_, err := NewS3(context.Background(), S3Config{Region: "us-east-1"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "aws load boom")
	})

	t.Run("success returns client", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}

		client, err := NewS3(context.Background(), S3Config{Region: "us-east-1"})
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("success with static credentials", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}

		client, err := NewS3(context.Background(), S3Config{
			Region:          "us-east-1",
			AccessKeyID:     "AKID",
			SecretAccessKey: "secret",
			SessionToken:    "token",
		})
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("success with endpoint and path style", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}

		client, err := NewS3(context.Background(), S3Config{
			Region:       "us-east-1",
			Endpoint:     "http://localhost:4566",
			UsePathStyle: true,
		})
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("success with connect timeout", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}

		client, err := NewS3(context.Background(), S3Config{
			Region:         "us-east-1",
			ConnectTimeout: 5 * time.Second,
		})
		require.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestNewS3Client(t *testing.T) {
	t.Run("validation error propagated", func(t *testing.T) {
		_, err := NewS3Client(context.Background(), S3Config{})
		require.Error(t, err)
	})

	t.Run("success returns wrapper", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}

		client, err := NewS3Client(context.Background(), S3Config{Region: "us-east-1"})
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Inner())
		assert.NoError(t, client.Close())
	})
}

func TestMustNewS3Client(t *testing.T) {
	t.Run("panics on validation error", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewS3Client(context.Background(), S3Config{})
		})
	})

	t.Run("returns client on success", func(t *testing.T) {
		origLoader := loadAWSConfig
		t.Cleanup(func() { loadAWSConfig = origLoader })

		loadAWSConfig = func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}

		assert.NotPanics(t, func() {
			c := MustNewS3Client(context.Background(), S3Config{Region: "us-east-1"})
			assert.NotNil(t, c)
		})
	})
}

func TestS3Client_Close(t *testing.T) {
	c := &S3Client{inner: nil}
	assert.NoError(t, c.Close())
}

func TestS3Client_Inner(t *testing.T) {
	c := &S3Client{inner: nil}
	assert.Nil(t, c.Inner())
}
