package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func generateTestKeys(t *testing.T) (pubPEM string, privPEM string) {
	t.Helper()

	// Generate a 2048-bit RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Encode Private Key to PKCS#8 format
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	privPEM = string(pem.EncodeToMemory(privBlock))

	// Encode Public Key to PKIX format
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	pubPEM = string(pem.EncodeToMemory(pubBlock))

	return pubPEM, privPEM
}

func TestRSACipher(t *testing.T) {
	pubPEM, privPEM := generateTestKeys(t)

	cfg := RSAConfig{
		PublicKeyPEM:  pubPEM,
		PrivateKeyPEM: privPEM,
	}

	cipher, err := NewRSA(cfg)
	if err != nil {
		t.Fatalf("failed to create RSA cipher: %v", err)
	}

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		plaintext := "OWASP Top 10 Confidential Data"

		ciphertext, err := cipher.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("encryption failed: %v", err)
		}

		if ciphertext == plaintext {
			t.Fatal("ciphertext is identical to plaintext")
		}

		decrypted, err := cipher.Decrypt(ciphertext)
		if err != nil {
			t.Fatalf("decryption failed: %v", err)
		}

		if decrypted != plaintext {
			t.Fatalf("expected %q but got %q", plaintext, decrypted)
		}
	})

	t.Run("Sign and Verify", func(t *testing.T) {
		payload := "Transaction: $5,000,000"

		signature, err := cipher.Sign(payload)
		if err != nil {
			t.Fatalf("signing failed: %v", err)
		}

		err = cipher.Verify(payload, signature)
		if err != nil {
			t.Fatalf("verification failed on valid signature: %v", err)
		}

		err = cipher.Verify("Transaction: $5", signature)
		if err == nil {
			t.Fatal("verification succeeded on tampered payload")
		}

		err = cipher.Verify(payload, "invalid_signature_base64_data")
		if err == nil {
			t.Fatal("verification succeeded on malformed signature")
		}
	})

	t.Run("Key Size Enforcement", func(t *testing.T) {
		// Generate a 1024-bit key (insecure)
		weakKey, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			t.Fatalf("failed to generate weak private key: %v", err)
		}
		weakPrivBytes, _ := x509.MarshalPKCS8PrivateKey(weakKey)
		weakPrivPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: weakPrivBytes}))

		_, err = NewRSA(RSAConfig{PrivateKeyPEM: weakPrivPEM})
		if err == nil || err.Error() != "RSA private key must be at least 2048 bits to comply with OWASP guidelines" {
			t.Fatalf("expected key size enforcement error, got: %v", err)
		}

		weakPubBytes, _ := x509.MarshalPKIXPublicKey(&weakKey.PublicKey)
		weakPubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: weakPubBytes}))
		_, err = NewRSA(RSAConfig{PublicKeyPEM: weakPubPEM})
		if err == nil || err.Error() != "RSA public key must be at least 2048 bits to comply with OWASP guidelines" {
			t.Fatalf("expected key size enforcement error for public key, got: %v", err)
		}
	})

	t.Run("Missing Keys", func(t *testing.T) {
		_, err := NewRSA(RSAConfig{})
		if err == nil {
			t.Fatal("expected error for missing keys")
		}
	})

	t.Run("Invalid PEM Blocks", func(t *testing.T) {
		_, err := NewRSA(RSAConfig{PublicKeyPEM: "invalid"})
		if err == nil {
			t.Fatal("expected error for invalid public key PEM")
		}

		_, err = NewRSA(RSAConfig{PrivateKeyPEM: "invalid"})
		if err == nil {
			t.Fatal("expected error for invalid private key PEM")
		}
	})

	t.Run("PKCS1 Private Key", func(t *testing.T) {
		key, _ := rsa.GenerateKey(rand.Reader, 2048)
		bytes := x509.MarshalPKCS1PrivateKey(key)
		pemStr := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: bytes}))

		_, err := NewRSA(RSAConfig{PrivateKeyPEM: pemStr})
		if err != nil {
			t.Fatalf("unexpected error parsing PKCS1 key: %v", err)
		}
	})

	t.Run("Invalid Key Type - Garbage", func(t *testing.T) {
		_, err := NewRSA(RSAConfig{PublicKeyPEM: "-----BEGIN PUBLIC KEY-----\nMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALoO8H/\n-----END PUBLIC KEY-----"})
		if err == nil {
			t.Fatal("expected error for malformed public key payload")
		}

		_, err = NewRSA(RSAConfig{PrivateKeyPEM: "-----BEGIN PRIVATE KEY-----\nMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALoO8H/\n-----END PRIVATE KEY-----"})
		if err == nil {
			t.Fatal("expected error for malformed private key payload")
		}
	})

	t.Run("Encrypt with Missing Public Key", func(t *testing.T) {
		c, _ := NewRSA(RSAConfig{PrivateKeyPEM: privPEM})
		_, err := c.Encrypt("test")
		if err == nil {
			t.Fatal("expected error due to missing public key")
		}

        err = c.Verify("payload", "sig")
        if err == nil {
			t.Fatal("expected error due to missing public key on Verify")
		}
	})

	t.Run("Decrypt and Sign with Missing Private Key", func(t *testing.T) {
		c, _ := NewRSA(RSAConfig{PublicKeyPEM: pubPEM})
		_, err := c.Decrypt("test")
		if err == nil {
			t.Fatal("expected error due to missing private key")
		}

		_, err = c.Sign("test")
		if err == nil {
			t.Fatal("expected error due to missing private key")
		}
	})

	t.Run("Message Too Long For Encryption", func(t *testing.T) {
		// generate 3000 byte string exceeding RSA 2048 OAEP limits limit
		longStr := make([]byte, 3000)
		_, err := cipher.Encrypt(string(longStr))
		if err == nil {
			t.Fatal("expected encryption error for overly long message")
		}
	})

	t.Run("MustNewRSA", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected no panic")
			}
		}()
		_ = MustNewRSA(cfg)
	})

	t.Run("MustNewRSA_Panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		_ = MustNewRSA(RSAConfig{})
	})
}
