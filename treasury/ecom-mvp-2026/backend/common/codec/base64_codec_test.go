package codec

import (
	"testing"
)

func TestBase64Coder_EncodeBase64(t *testing.T) {
	coder := NewBase64Coder()
	raw := "hello world"
	expected := "aGVsbG8gd29ybGQ="

	encoded := coder.EncodeBase64(raw)
	if encoded != expected {
		t.Errorf("expected %q, got %q", expected, encoded)
	}
}

func TestBase64Coder_DecodeBase64(t *testing.T) {
	coder := NewBase64Coder()

	t.Run("valid base64", func(t *testing.T) {
		encoded := "aGVsbG8gd29ybGQ="
		expected := "hello world"

		decoded, err := coder.DecodeBase64(encoded)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if decoded != expected {
			t.Errorf("expected %q, got %q", expected, decoded)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := coder.DecodeBase64("")
		if err == nil {
			t.Error("expected error for empty string")
		}
	})

	t.Run("invalid base64", func(t *testing.T) {
		_, err := coder.DecodeBase64("invalid bas64!@#")
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})
}
