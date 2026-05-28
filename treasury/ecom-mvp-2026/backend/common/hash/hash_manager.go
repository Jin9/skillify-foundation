package hash

import (
	"crypto/sha256"
	"encoding/base64"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/generator"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/validator"
)

const (
	saltCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
	saltLength  = 16
)

type HashManager interface {
	GenerateSalt() (string, error)
	HashSha256Encode(plainText string) string
	HashSha256EncodeSalt(plainText string, salt string) string
	HashSha256EncodePepper(plainText string) string
}

var _ HashManager = (*hashManager)(nil)

var generateSecureString = generator.GenerateSecureString

type HashManagerCfgs struct {
	Pepper string `validate:"required,len=16"`
}

type hashManager struct {
	pepper string
}

func NewHashManager(cfg HashManagerCfgs) (HashManager, error) {
	if err := validator.Validate(cfg); err != nil {
		return nil, err
	}

	return &hashManager{
		pepper: cfg.Pepper,
	}, nil
}

// MustNewHashManager is a convenience wrapper that panics on error.
func MustNewHashManager(cfg HashManagerCfgs) HashManager {
	hm, err := NewHashManager(cfg)
	if err != nil {
		panic(err)
	}
	return hm
}

func (s *hashManager) GenerateSalt() (string, error) {
	salt, err := generateSecureString(saltLength, saltCharset)
	if err != nil {
		return "", err
	}

	return salt, nil
}

func (s *hashManager) HashSha256Encode(plainText string) string {
	hash := sha256.Sum256([]byte(plainText))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (s *hashManager) HashSha256EncodeSalt(plainText, salt string) string {
	hash := sha256.Sum256([]byte(plainText + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (s *hashManager) HashSha256EncodePepper(plainText string) string {
	hash := sha256.Sum256([]byte(plainText + s.pepper))
	return base64.StdEncoding.EncodeToString(hash[:])
}
