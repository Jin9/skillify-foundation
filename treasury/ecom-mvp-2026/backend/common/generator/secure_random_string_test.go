package generator

import (
	"crypto/rand"
	"errors"
	"strings"
	"testing"
)

func TestGenerateSecureString(t *testing.T) {
	original := randReader
	defer func() { randReader = original }()

	t.Run("success", func(t *testing.T) {
		randReader = rand.Reader // default reader
		s, err := GenerateSecureString(16, "abc123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s) != 16 {
			t.Errorf("expected length 16, got %d", len(s))
		}
	})

	t.Run("empty charset", func(t *testing.T) {
		_, err := GenerateSecureString(10, "")
		if err == nil || !strings.Contains(err.Error(), "charset") {
			t.Errorf("expected error for empty charset")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		_, err := GenerateSecureString(0, "abc")
		if err == nil || !strings.Contains(err.Error(), "length") {
			t.Errorf("expected error for zero length")
		}
	})

	t.Run("rand read error", func(t *testing.T) {
		randReader = fakeFailReader{}
		_, err := GenerateSecureString(5, "abc")
		if err == nil || !strings.Contains(err.Error(), "forced error") {
			t.Errorf("expected read error, got %v", err)
		}
	})

	t.Run("maxByte boundary", func(t *testing.T) {
		// Provide 255 which is > 252 (maxByte for len 6)
		// it will be skipped, and we should still get valid characters
		randReader = fakeSeqReader{data: []byte{255, 0, 1, 2, 3}}
		s, err := GenerateSecureString(3, "ab") // len 2. maxByte = 255 - (255%2) = 255 - 1 = 254
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s) != 3 {
			t.Errorf("expected length 3, got %d", len(s))
		}
	})

}

type fakeSeqReader struct {
	data []byte
	pos  int
}

func (r fakeSeqReader) Read(p []byte) (int, error) {
	n := copy(p, r.data)
	return n, nil // simple mock
}

type fakeFailReader struct{}

func (fakeFailReader) Read(p []byte) (int, error) {
	return 0, errors.New("forced error")
}

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func BenchmarkGenerateSecureString(b *testing.B) {
	for b.Loop() {
		_, err := GenerateSecureString(32, charset)
		if err != nil {
			b.Fatal(err)
		}
	}
}
