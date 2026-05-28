package generator

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAppUUID_GenerateUUID(t *testing.T) {
	u := NewAppUUID()
	id := u.GenerateUUID()
	if id == "" {
		t.Fatal("expected non-empty UUID string")
	}
	if _, err := uuid.Parse(id); err != nil {
		t.Fatalf("expected valid UUID, got %q: %v", id, err)
	}
}

func TestAppUUID_GenerateUUIDTypeStrict(t *testing.T) {
	u := NewAppUUID()
	id := u.GenerateUUIDTypeStrict()
	if id == uuid.Nil {
		t.Fatal("expected non-nil UUID")
	}
}

func TestAppUUID_GenerateUUIDV7(t *testing.T) {
	u := NewAppUUID()
	id := u.GenerateUUIDV7()
	if id == uuid.Nil {
		t.Fatal("expected non-nil UUID v7")
	}
	if id.Version() != 7 {
		t.Fatalf("expected UUID v7, got v%d", id.Version())
	}
}

func TestGeneratePwd(t *testing.T) {
	key := "my-key"
	date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	pwd := GeneratePwd(key, date, 16)

	if len(pwd) != 16 {
		t.Fatalf("expected length 16, got %d", len(pwd))
	}

	// Deterministic — same input should produce same output
	pwd2 := GeneratePwd(key, date, 16)
	if pwd != pwd2 {
		t.Fatalf("expected deterministic output, got %q and %q", pwd, pwd2)
	}

	// Different key produces different output
	pwd3 := GeneratePwd("different-key", date, 16)
	if pwd == pwd3 {
		t.Fatal("expected different output for different key")
	}

	// Different date produces different output
	pwd4 := GeneratePwd(key, date.AddDate(0, 0, 1), 16)
	if pwd == pwd4 {
		t.Fatal("expected different output for different date")
	}
}

func TestGeneratePwd_LongLength(t *testing.T) {
	key := "test"
	date := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	// SHA-512 hex is 128 chars. Request more than 128
	pwd := GeneratePwd(key, date, 200)
	if len(pwd) != 128 {
		t.Fatalf("expected full hex string length 128, got %d", len(pwd))
	}
}

func TestToHexString(t *testing.T) {
	result := toHexString([]byte{0, 15, 255})
	if result != "000fff" {
		t.Fatalf("expected 000fff, got %q", result)
	}
}
