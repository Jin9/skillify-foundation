package clock

import (
	"testing"
	"time"
)

func TestClock_Now(t *testing.T) {
	c := New()
	now := c.Now()

	if now.Location() != time.UTC {
		t.Errorf("Now() should return UTC, got %v", now.Location())
	}
}

func TestClock_NowIn(t *testing.T) {
	c := New()
	bkk, err := c.NowIn(Bangkok)
	if err != nil {
		t.Fatalf("NowIn(%q) error: %v", Bangkok, err)
	}
	if bkk.Location().String() != Bangkok {
		t.Errorf("expected location %q, got %q", Bangkok, bkk.Location())
	}
}

func TestClock_NowIn_InvalidLocation(t *testing.T) {
	c := New()
	_, err := c.NowIn("Invalid/Location")
	if err == nil {
		t.Error("expected error for invalid location")
	}
}

func TestClock_Parse(t *testing.T) {
	c := New()
	parsed, err := c.Parse(DefaultLayout, "2025-06-15T10:30:00")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if parsed.Year() != 2025 || parsed.Month() != 6 || parsed.Day() != 15 {
		t.Errorf("unexpected parsed time: %v", parsed)
	}
}

func TestClock_ParseIn(t *testing.T) {
	c := New()
	parsed, err := c.ParseIn(DefaultLayout, "2025-06-15T10:30:00", Bangkok)
	if err != nil {
		t.Fatalf("ParseIn() error: %v", err)
	}
	if parsed.Location().String() != Bangkok {
		t.Errorf("expected location %q, got %q", Bangkok, parsed.Location())
	}
}

func TestNowUTC(t *testing.T) {
	now := NowUTC()
	if now.Location() != time.UTC {
		t.Errorf("NowUTC() should return UTC, got %v", now.Location())
	}
}

func TestNowBangkok(t *testing.T) {
	bkk, err := NowBangkok()
	if err != nil {
		t.Fatalf("NowBangkok() error: %v", err)
	}
	if bkk.Location().String() != Bangkok {
		t.Errorf("expected location %q, got %q", Bangkok, bkk.Location())
	}
}

func TestParseUTC(t *testing.T) {
	parsed, err := ParseUTC("2025-06-15T10:30:00")
	if err != nil {
		t.Fatalf("ParseUTC() error: %v", err)
	}
	if parsed.Location() != time.UTC {
		t.Errorf("ParseUTC() should return UTC, got %v", parsed.Location())
	}
}

func TestParseUTC_Invalid(t *testing.T) {
	_, err := ParseUTC("not-a-time")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestParseBangkok(t *testing.T) {
	parsed, err := ParseBangkok("2025-06-15T10:30:00")
	if err != nil {
		t.Fatalf("ParseBangkok() error: %v", err)
	}
	if parsed.Location().String() != Bangkok {
		t.Errorf("expected location %q, got %q", Bangkok, parsed.Location())
	}
}

func TestToBangkok(t *testing.T) {
	utcTime := time.Date(2025, 6, 15, 3, 0, 0, 0, time.UTC)
	bkk, err := ToBangkok(utcTime)
	if err != nil {
		t.Fatalf("ToBangkok() error: %v", err)
	}
	if bkk.Hour() != 10 {
		t.Errorf("expected hour 10 (UTC+7), got %d", bkk.Hour())
	}
}
