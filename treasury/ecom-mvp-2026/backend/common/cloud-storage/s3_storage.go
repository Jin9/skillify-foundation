package cloudstorage

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type awsConfigLoader func(context.Context, ...func(*awsconfig.LoadOptions) error) (aws.Config, error)

var loadAWSConfig awsConfigLoader = awsconfig.LoadDefaultConfig

// S3Config configures an AWS S3 client.
//
// Static credentials are optional. If omitted, the AWS SDK default credential chain is used.
// Endpoint is optional and useful for local S3-compatible services such as LocalStack or MinIO.
// ConnectTimeout applies only to configuration loading and client creation.
type S3Config struct {
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	UsePathStyle    bool
	ConnectTimeout  time.Duration
}

func (cfg S3Config) Validate() error {
	if cfg.Region == "" {
		return errors.New("cloudstorage: Region is required")
	}

	hasAccessKey := cfg.AccessKeyID != ""
	hasSecretKey := cfg.SecretAccessKey != ""
	if hasAccessKey != hasSecretKey {
		return errors.New("cloudstorage: AccessKeyID and SecretAccessKey must be provided together")
	}

	return nil
}

// NewS3 creates an AWS S3 client.
func NewS3(ctx context.Context, cfg S3Config, clientOptions ...func(*s3.Options)) (*s3.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	ctx2, cancel := withOptionalTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	loadOptions := make([]func(*awsconfig.LoadOptions) error, 0, 2)
	loadOptions = append(loadOptions, awsconfig.WithRegion(cfg.Region))

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		loadOptions = append(loadOptions,
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)),
		)
	}

	awsCfg, err := loadAWSConfig(ctx2, loadOptions...)
	if err != nil {
		return nil, err
	}

	options := make([]func(*s3.Options), 0, len(clientOptions)+2)
	if cfg.Endpoint != "" {
		endpoint := cfg.Endpoint
		options = append(options, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	if cfg.UsePathStyle {
		options = append(options, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}
	options = append(options, clientOptions...)

	return s3.NewFromConfig(awsCfg, options...), nil
}

// S3Client is a thin wrapper around the official AWS S3 client.
type S3Client struct {
	inner *s3.Client
}

// NewS3Client creates a wrapper client.
func NewS3Client(ctx context.Context, cfg S3Config, clientOptions ...func(*s3.Options)) (*S3Client, error) {
	c, err := NewS3(ctx, cfg, clientOptions...)
	if err != nil {
		return nil, err
	}

	return &S3Client{inner: c}, nil
}

// MustNewS3Client is a convenience wrapper that panics on error.
func MustNewS3Client(ctx context.Context, cfg S3Config, clientOptions ...func(*s3.Options)) *S3Client {
	c, err := NewS3Client(ctx, cfg, clientOptions...)
	if err != nil {
		panic(err)
	}

	return c
}

func (c *S3Client) Inner() *s3.Client {
	return c.inner
}

func (c *S3Client) Close() error {
	return nil
}
