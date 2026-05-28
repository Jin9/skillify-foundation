package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/generator"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/validator"
)

const (
	ivCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
	ivLength  = 12
)

type Cipher interface {
	GenerateIV() ([]byte, error)
	Encrypt(plaintext string, iv []byte) (string, error)
	Decrypt(encodedPayload string) (string, error)
}

var _ Cipher = (*service)(nil)

type Config struct {
	Key string `validate:"required,len=32"`
}

type service struct {
	key []byte
}

func New(cfg Config) (Cipher, error) {
	if err := validator.Validate(cfg); err != nil {
		return nil, err
	}

	return &service{
		key: []byte(cfg.Key),
	}, nil
}

// MustNew is a convenience wrapper that panics on error.
func MustNew(cfg Config) Cipher {
	a, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return a
}

// GenerateIV generates a random IV of ivLength bytes.
func (a *service) GenerateIV() ([]byte, error) {
	b, err := generator.GenerateSecureString(ivLength, ivCharset)
	if err != nil {
		return nil, fmt.Errorf("generate secure string: %w", err)
	}
	return []byte(b), nil
}

// Encrypt encrypts the plaintext using AES-GCM with the provided IV
func (a *service) Encrypt(plaintext string, iv []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext) + "|" + base64.StdEncoding.EncodeToString(iv), nil
}

// Decrypt decrypts the ciphertext using AES-GCM with the provided IV
func (a *service) Decrypt(encodedPayload string) (string, error) {
	parts := strings.Split(string(encodedPayload), "|")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid payload format")
	}

	decoded := make([][]byte, len(parts))
	labels := []string{"ciphertext", "iv"}

	for i, part := range parts {
		data, err := base64.StdEncoding.DecodeString(part)
		if err != nil {
			return "", fmt.Errorf("decode base64 %s: %w", labels[i], err)
		}
		decoded[i] = data
	}

	ciphertext, iv := decoded[0], decoded[1]

	block, err := aes.NewCipher(a.key) // requires 32 bytes for AES-256
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block) // GCM mode
	if err != nil {
		return "", fmt.Errorf("new GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}
