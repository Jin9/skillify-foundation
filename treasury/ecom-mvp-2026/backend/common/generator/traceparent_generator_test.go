package generator

import (
	"errors"
	"math/rand/v2"
	"sync"
	"testing"
)

func TestNewTraceParent(t *testing.T) {
	tp := NewTraceParent()
	s := tp.String()

	parsed, err := Parse(s)
	if err != nil {
		t.Fatalf("valid traceparent %q, error %s", s, err)
	}
	if parsed.SpanID.String() == "" {
		t.Error("expect non-empty span-id")
	}
	if parsed.TraceID.String() == "" {
		t.Error("expect non-empty trace-id")
	}
	if parsed.TraceFlags != TraceFlagsSampled {
		t.Errorf("expected TraceFlagsSampled=%d but got %d", TraceFlagsSampled, parsed.TraceFlags)
	}
}

func TestIDGeneratorUniqueness(t *testing.T) {
	gen := getIDGenerator()

	spanID1 := gen.NewSpanID()
	traceID1 := gen.NewTraceID()

	if len(spanID1) != 8 {
		t.Errorf("span-id len expect 8 but actual %d", len(spanID1))
	}
	if len(traceID1) != 16 {
		t.Errorf("trace-id len expect 16 but actual %d", len(traceID1))
	}

	spanID2 := gen.NewSpanID()
	traceID2 := gen.NewTraceID()

	if spanID1 == spanID2 {
		t.Error("span should be random every time")
	}
	if traceID1 == traceID2 {
		t.Error("trace should be random every time")
	}
}

func TestParse(t *testing.T) {
	t.Run("valid traceparent", func(t *testing.T) {
		tp, err := Parse("00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tp.SpanID.String() != "b7ad6b7169203331" {
			t.Errorf("expect b7ad6b7169203331 but got %q", tp.SpanID.String())
		}
		if tp.TraceID.String() != "0af7651916cd43dd8448eb211c80319c" {
			t.Errorf("expect 0af7651916cd43dd8448eb211c80319c but got %q", tp.TraceID.String())
		}
	})

	t.Run("empty", func(t *testing.T) {
		_, err := Parse("")
		if !errors.Is(err, ErrEmptyTraceParent) {
			t.Fatalf("expected ErrEmptyTraceParent but got %v", err)
		}
	})

	t.Run("unsupported version", func(t *testing.T) {
		_, err := Parse("01-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		if !errors.Is(err, ErrWrongTraceParent) {
			t.Fatalf("expected ErrWrongTraceParent but got %v", err)
		}
	})

	t.Run("all-zero trace-id", func(t *testing.T) {
		_, err := Parse("00-00000000000000000000000000000000-b7ad6b7169203331-01")
		if !errors.Is(err, ErrInvalidTraceID) {
			t.Fatalf("expected ErrInvalidTraceID but got %v", err)
		}
	})

	t.Run("all-zero span-id", func(t *testing.T) {
		_, err := Parse("00-0af7651916cd43dd8448eb211c80319c-0000000000000000-01")
		if !errors.Is(err, ErrInvalidSpanID) {
			t.Fatalf("expected ErrInvalidSpanID but got %v", err)
		}
	})

	t.Run("wrong format - no dashes", func(t *testing.T) {
		_, err := Parse("wrongformat")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("wrong version len", func(t *testing.T) {
		_, err := Parse("0-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("version parse string error", func(t *testing.T) {
		_, err := Parse("xx-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("trace hex parse error", func(t *testing.T) {
		_, err := Parse("00-0af7651916cd43dd8448eb211c80319z-b7ad6b7169203331-01")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("trace len error", func(t *testing.T) {
		// 30 hex chars (15 bytes) instead of 32
		_, err := Parse("00-0af7651916cd43dd8448eb211c8031-b7ad6b7169203331-01")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("span hex parse error", func(t *testing.T) {
		_, err := Parse("00-0af7651916cd43dd8448eb211c80319c-b7ad6b716920333z-01")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("span len error", func(t *testing.T) {
		// 14 hex chars (7 bytes) instead of 16 chars
		_, err := Parse("00-0af7651916cd43dd8448eb211c80319c-b7ad6b71692033-01")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("flag hex parse error", func(t *testing.T) {
		_, err := Parse("00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-0z")
		if err == nil { t.Fatal("expected err") }
	})

	t.Run("flag len error", func(t *testing.T) {
		// double flag: 4 chars (2 bytes)
		_, err := Parse("00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-0001")
		if err == nil { t.Fatal("expected err") }
	})
}

func TestCreateChild(t *testing.T) {
	parent := NewTraceParent()
	child := parent.CreateChild()

	if child.TraceID != parent.TraceID {
		t.Error("child should inherit parent trace-id")
	}
	if child.SpanID == parent.SpanID {
		t.Error("child should have a new span-id")
	}
	if child.TraceFlags != parent.TraceFlags {
		t.Error("child should inherit parent trace-flags")
	}
}

func TestGeneratorFallbackErrors(t *testing.T) {
	origCryptoRead := cryptoRandRead
	origChaChaRead := chaChaRead

	defer func() {
		cryptoRandRead = origCryptoRead
		chaChaRead = origChaChaRead
		generatorOnce = sync.Once{}
	}()

	cryptoRandRead = func(b []byte) (int, error) {
		return 0, errors.New("mock crypto error")
	}
	chaChaRead = func(c *rand.ChaCha8, b []byte) (int, error) {
		return 0, errors.New("mock chacha error")
	}

	generatorOnce = sync.Once{} // Reset to trigger coverage of seed error
	gen := getIDGenerator()

	sid := gen.NewSpanID()
	if sid[7] != 1 {
		t.Error("expected span id fallback not applied")
	}

	tid := gen.NewTraceID()
	if tid[15] != 1 {
		t.Error("expected trace id fallback not applied")
	}
}
