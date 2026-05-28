package hash

import (
	"errors"
	"testing"
)

func TestHashManager_NewHashManager(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := HashManagerCfgs{Pepper: "1234567812345678"} // exactly 16 bytes
		hm, err := NewHashManager(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hm == nil {
			t.Fatal("expected non-nil hash manager")
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		cfg := HashManagerCfgs{Pepper: "short"}
		_, err := NewHashManager(cfg)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestHashManager_MustNewHashManager(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unexpected panic: %v", r)
			}
		}()
		cfg := HashManagerCfgs{Pepper: "1234567812345678"}
		hm := MustNewHashManager(cfg)
		if hm == nil {
			t.Fatal("expected non-nil hash manager")
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		cfg := HashManagerCfgs{Pepper: "short"}
		_ = MustNewHashManager(cfg)
	})
}

func TestHashManager_GenerateSalt(t *testing.T) {
	cfg := HashManagerCfgs{Pepper: "1234567812345678"}
	hm := MustNewHashManager(cfg)

	t.Run("success", func(t *testing.T) {
		salt, err := hm.GenerateSalt()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(salt) != 16 {
			t.Errorf("expected length 16, got %d", len(salt))
		}
	})

	t.Run("error", func(t *testing.T) {
		orig := generateSecureString
		defer func() { generateSecureString = orig }()
		generateSecureString = func(length int, charset string) (string, error) {
			return "", errors.New("mock generator error")
		}

		_, err := hm.GenerateSalt()
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestHashManager_HashMethods(t *testing.T) {
	cfg := HashManagerCfgs{Pepper: "1234567812345678"}
	hm := MustNewHashManager(cfg)

	plainText := "my-secret-password"
	salt := "my-salt"

	hash1 := hm.HashSha256Encode(plainText)
	if hash1 == "" {
		t.Error("expected non-empty hash")
	}

	hash2 := hm.HashSha256EncodeSalt(plainText, salt)
	if hash2 == "" || hash2 == hash1 {
		t.Error("expected non-empty, distinct hash for salt")
	}

	hash3 := hm.HashSha256EncodePepper(plainText)
	if hash3 == "" || hash3 == hash1 || hash3 == hash2 {
		t.Error("expected non-empty, distinct hash for pepper")
	}
}
