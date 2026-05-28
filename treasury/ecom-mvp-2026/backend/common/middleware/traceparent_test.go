package middleware

import (
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func assertHexBytes(t *testing.T, label string, s string, wantBytes int) {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("%s should be hex, got %q (err=%v)", label, s, err)
	}
	if len(b) != wantBytes {
		t.Fatalf("%s should be %d bytes, got %d (value=%q)", label, wantBytes, len(b), s)
	}
}

func TestNewTraceparent(t *testing.T) {
	traceparent := NewTraceParent().String()

	tp, err := Parse(traceparent)
	if err != nil {
		t.Errorf("valid traceparent %q, error %s\n", traceparent, err)
	}

	if tp.SpanID.String() == "" {
		t.Error("a new traceparent expect not empty span-id\n")
	}

	if tp.TraceID.String() == "" {
		t.Error("a new traceparent expect not empty trace-id\n")
	}

	if tp.TraceFlags != TraceFlagsSampled {
		t.Errorf("expected TraceFlagsSampled=%d but got %d", TraceFlagsSampled, tp.TraceFlags)
	}
	if got := tp.String(); got[len(got)-3:] != "-01" {
		t.Errorf("expected traceparent to end with -01 but got %q", got)
	}
}

func TestParseTraceparent(t *testing.T) {
	t.Run("valid traceparent", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
		tp, err := Parse(traceparent)
		if err != nil {
			t.Errorf("valid traceparent %q, error %s\n", traceparent, err)
		}

		if tp.SpanID.String() != "b7ad6b7169203331" {
			t.Errorf("expect b7ad6b7169203331 but got %q\n", tp.SpanID.String())
		}
		if tp.TraceID.String() != "0af7651916cd43dd8448eb211c80319c" {
			t.Errorf("expect 0af7651916cd43dd8448eb211c80319c but got %q\n", tp.TraceID.String())
		}
	})
	t.Run("unsupported version", func(t *testing.T) {
		traceparent := "01-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
		_, err := Parse(traceparent)
		if err == nil {
			t.Fatalf("expected error for unsupported version")
		}
		if !errors.Is(err, ErrWrongTraceParent) {
			t.Fatalf("expected ErrWrongTraceParent but got %v", err)
		}
	})
	t.Run("invalid version format", func(t *testing.T) {
		traceparent := "zz-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
		_, err := Parse(traceparent)
		if err == nil {
			t.Fatalf("expected error for invalid version format")
		}
		if !errors.Is(err, ErrWrongTraceParent) {
			t.Fatalf("expected ErrWrongTraceParent but got %v", err)
		}
	})
	t.Run("invalid all-zero trace-id", func(t *testing.T) {
		traceparent := "00-00000000000000000000000000000000-b7ad6b7169203331-01"
		_, err := Parse(traceparent)
		if err == nil {
			t.Fatalf("expected error for all-zero trace-id")
		}
		if !errors.Is(err, ErrInvalidTraceID) {
			t.Fatalf("expected ErrInvalidTraceID but got %v", err)
		}
	})
	t.Run("invalid all-zero span-id", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c-0000000000000000-01"
		_, err := Parse(traceparent)
		if err == nil {
			t.Fatalf("expected error for all-zero span-id")
		}
		if !errors.Is(err, ErrInvalidSpanID) {
			t.Fatalf("expected ErrInvalidSpanID but got %v", err)
		}
	})
	t.Run("invalid pattern traceparent", func(t *testing.T) {
		traceparent := "5f2b8701-6c14-436f-8e69-92f1670f0ec2"
		tp, err := Parse(traceparent)
		if err == nil {
			t.Errorf("invalid pattern traceparent expect error but we've got trace-id %q and span-id %q\n\n", tp.TraceID, tp.SpanID)
		}
	})
	t.Run("empty traceparent", func(t *testing.T) {
		traceparent := ""
		tp, err := Parse(traceparent)
		if err == nil {
			t.Errorf("empty pattern traceparent expect error but we've got trace-id %q and span-id %q\n\n", tp.TraceID, tp.SpanID)
		}
	})
	t.Run("invalid trace-id length", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c00-b7ad6b7169203331-01"
		_, err := Parse(traceparent)
		if err == nil {
			t.Errorf("invalid trace-id length expect error")
		}
	})
	t.Run("invalid span-id length", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331aa-01"
		_, err := Parse(traceparent)
		if err == nil {
			t.Errorf("invalid span-id length expect error")
		}
	})
	t.Run("invalid trace-flags length", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-0101"
		_, err := Parse(traceparent)
		if err == nil {
			t.Errorf("invalid trace-flags length expect error")
		}
	})
}

func TestTraceContextTraceIDMiddleware(t *testing.T) {
	t.Run("ref-id key as traceparent", func(t *testing.T) {
		handler := func(c *gin.Context) {
			if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
				c.String(200, v)
				return
			}
			c.String(500, "not found")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)

		req.Header.Add(traceparentHeaderKey, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		c.Request = req

		e.Use(TraceContextTraceIDMiddleware(traceparentHeaderKey))

		e.POST("/api", handler)

		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if w.Body.String() != "0af7651916cd43dd8448eb211c80319c" {
			t.Errorf("expect ref-id to be trace-id but got %q\n", w.Body.String())
		}
	})
	t.Run("default header key", func(t *testing.T) {
		handler := func(c *gin.Context) {
			if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
				c.String(200, v)
				return
			}
			c.String(500, "not found")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Add(traceparentHeaderKey, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		c.Request = req

		e.Use(TraceContextTraceIDMiddleware(""))
		e.POST("/api", handler)
		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if w.Body.String() != "0af7651916cd43dd8448eb211c80319c" {
			t.Errorf("expect trace-id from traceparent header")
		}
	})
	t.Run("preserve existing ref-id", func(t *testing.T) {
		handler := func(c *gin.Context) {
			if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
				c.String(200, v)
				return
			}
			c.String(500, "not found")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		c.Request = req

		c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), "existing-ref"))

		e.Use(TraceContextTraceIDMiddleware(""))
		e.POST("/api", handler)
		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if w.Body.String() != "existing-ref" {
			t.Errorf("expect existing ref-id preserved")
		}
	})
	t.Run("ref-id key as traceparent but empty", func(t *testing.T) {
		handler := func(c *gin.Context) {
			if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
				c.String(200, v)
				return
			}
			c.String(500, "not found")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)

		c.Request = req

		e.Use(TraceContextTraceIDMiddleware(traceparentHeaderKey))

		e.POST("/api", handler)

		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if got := w.Body.String(); got == "" {
			t.Fatalf("expect ref-id not empty\n")
		} else {
			assertHexBytes(t, "ref-id", got, 16)
		}
	})

	t.Run("invalid traceparent header falls back", func(t *testing.T) {
		handler := func(c *gin.Context) {
			c.String(200, RefID(c))
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Add(traceparentHeaderKey, "01-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		c.Request = req

		e.Use(TraceContextTraceIDMiddleware(""))
		e.POST("/api", handler)
		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Fatalf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		assertHexBytes(t, "ref-id", w.Body.String(), 16)
	})
}
