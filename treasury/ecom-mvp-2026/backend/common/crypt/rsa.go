package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// AsymmetricCipher defines operations for OWASP-compliant RSA asymmetric cryptography.
type AsymmetricCipher interface {
	// Encrypt uses RSA-OAEP with SHA-256 to encrypt the payload.
	Encrypt(plaintext string) (string, error)

	// Decrypt uses RSA-OAEP with SHA-256 to decrypt the payload.
	Decrypt(ciphertext string) (string, error)

	// Sign uses RSA-PSS with SHA-256 to create a secure cryptographic signature.
	Sign(payload string) (string, error)

	// Verify validates an RSA-PSS signature.
	Verify(payload, signature string) error
}

// RSAConfig holds the PEM-encoded keys for the RSA cipher.
// Both are optional depending on the desired operation (e.g. Encrypt/Verify only needs PublicKeyPEM).
type RSAConfig struct {
	PublicKeyPEM  string
	PrivateKeyPEM string
}

type rsaService struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRSA creates a new AsymmetricCipher.
func NewRSA(cfg RSAConfig) (AsymmetricCipher, error) {
	service := &rsaService{}

	if cfg.PublicKeyPEM != "" {
		block, _ := pem.Decode([]byte(cfg.PublicKeyPEM))
		if block == nil {
			return nil, errors.New("failed to parse public key PEM block")
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKIX public key: %w", err)
		}

		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("provided public key is not an RSA key")
		}

		if rsaPub.Size() < 256 { // 2048 bits = 256 bytes
			return nil, errors.New("RSA public key must be at least 2048 bits to comply with OWASP guidelines")
		}

		service.publicKey = rsaPub
	}

	if cfg.PrivateKeyPEM != "" {
		block, _ := pem.Decode([]byte(cfg.PrivateKeyPEM))
		if block == nil {
			return nil, errors.New("failed to parse private key PEM block")
		}

		// Try PKCS#8 first, then fallback to PKCS#1
		var rsaPriv *rsa.PrivateKey
		priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			rsaPriv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key (neither PKCS#8 nor PKCS#1): %w", err)
			}
		} else {
			var ok bool
			rsaPriv, ok = priv.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("provided private key is not an RSA key")
			}
		}

		if rsaPriv.Size() < 256 { // 2048 bits = 256 bytes
			return nil, errors.New("RSA private key must be at least 2048 bits to comply with OWASP guidelines")
		}

		service.privateKey = rsaPriv
	}

	if service.publicKey == nil && service.privateKey == nil {
		return nil, errors.New("at least one of PublicKeyPEM or PrivateKeyPEM must be provided")
	}

	return service, nil
}

// Encrypt encrypts the given string utilizing RSA-OAEP and SHA-256.
func (r *rsaService) Encrypt(plaintext string) (string, error) {
	if r.publicKey == nil {
		return "", errors.New("public key is required for encryption")
	}

	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, r.publicKey, []byte(plaintext), nil)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given string utilizing RSA-OAEP and SHA-256.
func (r *rsaService) Decrypt(ciphertext string) (string, error) {
	if r.privateKey == nil {
		return "", errors.New("private key is required for decryption")
	}

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, r.privateKey, decoded, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// Sign signs the payload utilizing RSA-PSS and SHA-256.
func (r *rsaService) Sign(payload string) (string, error) {
	if r.privateKey == nil {
		return "", errors.New("private key is required for signing")
	}

	hashed := sha256.Sum256([]byte(payload))
	signature, err := rsa.SignPSS(rand.Reader, r.privateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify validates the cryptographic signature utilizing RSA-PSS and SHA-256.
func (r *rsaService) Verify(payload, signature string) error {
	if r.publicKey == nil {
		return errors.New("public key is required for verification")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode base64 signature: %w", err)
	}

	hashed := sha256.Sum256([]byte(payload))
	err = rsa.VerifyPSS(r.publicKey, crypto.SHA256, hashed[:], sigBytes, nil)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// MustNewRSA is a convenience wrapper that panics on error.
func MustNewRSA(cfg RSAConfig) AsymmetricCipher {
	s, err := NewRSA(cfg)
	if err != nil {
		panic(err)
	}
	return s
}
