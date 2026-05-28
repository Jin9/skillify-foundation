package crypt

import (
	"strings"
	"testing"
)

func TestAESCipher(t *testing.T) {
	key := "01234567890123456789012345678901" // 32 bytes

	c, err := New(Config{Key: key})
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}

	t.Run("GenerateIV", func(t *testing.T) {
		iv, err := c.GenerateIV()
		if err != nil {
			t.Fatalf("generate iv: %v", err)
		}
		if len(iv) != 12 {
			t.Fatalf("expected IV length 12, got %d", len(iv))
		}
	})

	t.Run("Encrypt and Decrypt roundtrip", func(t *testing.T) {
		iv, _ := c.GenerateIV()

		plaintext := "hello-world-secret"
		cipher, err := c.Encrypt(plaintext, iv)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}

		if cipher == plaintext {
			t.Fatal("ciphertext should not equal plaintext")
		}

		if !strings.Contains(cipher, "|") {
			t.Fatal("expected ciphertext|iv format")
		}

		decrypted, err := c.Decrypt(cipher)
		if err != nil {
			t.Fatalf("decrypt: %v", err)
		}

		if decrypted != plaintext {
			t.Fatalf("expected %q, got %q", plaintext, decrypted)
		}
	})

	t.Run("Decrypt with invalid format", func(t *testing.T) {
		_, err := c.Decrypt("no-pipe-separator")
		if err == nil {
			t.Fatal("expected error for invalid format")
		}
	})

	t.Run("Decrypt with invalid base64 ciphertext", func(t *testing.T) {
		_, err := c.Decrypt("!!!invalid!!!|dGVzdA==")
		if err == nil {
			t.Fatal("expected error for invalid base64 ciphertext")
		}
	})

	t.Run("Decrypt with invalid base64 iv", func(t *testing.T) {
		_, err := c.Decrypt("dGVzdA==|!!!invalid!!!")
		if err == nil {
			t.Fatal("expected error for invalid base64 iv")
		}
	})

	t.Run("Decrypt with wrong key produces error", func(t *testing.T) {
		iv, _ := c.GenerateIV()
		cipher, _ := c.Encrypt("test", iv)

		otherKey := "99999999999999999999999999999999"
		otherC, _ := New(Config{Key: otherKey})
		_, err := otherC.Decrypt(cipher)
		if err == nil {
			t.Fatal("expected error decrypting with wrong key")
		}
	})
}

func TestAESNew_InvalidConfig(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		_, err := New(Config{Key: ""})
		if err == nil {
			t.Fatal("expected error for empty key")
		}
	})

	t.Run("wrong key length", func(t *testing.T) {
		_, err := New(Config{Key: "too-short"})
		if err == nil {
			t.Fatal("expected error for wrong key length")
		}
	})
}

func TestMustNew_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustNew(Config{})
}

func TestMustNew_Success(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("expected no panic")
		}
	}()
	c := MustNew(Config{Key: "01234567890123456789012345678901"})
	if c == nil {
		t.Fatal("expected non-nil cipher")
	}
}

func TestEncrypt_CipherFailure(t *testing.T) {
	// Bypass validation to force aes.NewCipher failure
	invalidSvc := &service{key: []byte("short")}
	_, err := invalidSvc.Encrypt("test", []byte("iv"))
	if err == nil {
		t.Fatal("expected encryption error with bad key")
	}
}

func TestDecrypt_CipherFailure(t *testing.T) {
	// Bypass validation to force aes.NewCipher failure
	invalidSvc := &service{key: []byte("short")}
	payload := "dGVzdA==|dGVzdA==" // Valid base64 strings mimicking ciphertext|iv
	_, err := invalidSvc.Decrypt(payload)
	if err == nil {
		t.Fatal("expected decryption error with bad key")
	}
}
