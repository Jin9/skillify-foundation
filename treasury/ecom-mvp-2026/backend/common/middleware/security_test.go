package middleware

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	SecurityHeaders()(ctx)

	headers := recorder.Header()
	assertHeader(t, headers, "X-Frame-Options", "DENY")
	assertHeader(t, headers, "Cross-Origin-Opener-Policy", "same-origin")
	assertHeader(t, headers, "Cross-Origin-Resource-Policy", "same-origin")
	assertHeader(t, headers, "X-Permitted-Cross-Domain-Policies", "none")
	assertHeader(t, headers, "Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	assertHeader(t, headers, "Referrer-Policy", "strict-origin")
	assertHeader(t, headers, "X-Content-Type-Options", "nosniff")
}

func TestAccessControl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("sets CORS response headers", func(t *testing.T) {
		r := gin.New()
		r.Use(AccessControl("https://example.com", []string{"Authorization", "Content-Type"}))
		nextCalled := false
		r.GET("/", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(recorder, req)

		headers := recorder.Header()
		assertHeader(t, headers, "Access-Control-Allow-Origin", "https://example.com")
		assertHeader(t, headers, "Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		assertHeader(t, headers, "Access-Control-Allow-Headers", "Authorization,Content-Type")

		varyValues := headers.Values("Vary")
		assertContains(t, varyValues, "Origin")
		assertContains(t, varyValues, "Access-Control-Request-Method")
		assertContains(t, varyValues, "Access-Control-Request-Headers")

		if !nextCalled {
			t.Fatalf("expected next handler to be called")
		}
	})

	t.Run("options request aborts with no content", func(t *testing.T) {
		r := gin.New()
		r.Use(AccessControl("https://example.com", []string{"Authorization"}))
		nextCalled := false
		r.OPTIONS("/", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
		}

		if nextCalled {
			t.Fatalf("expected next handler not to be called for OPTIONS request")
		}
	})
}

func assertHeader(t *testing.T, headers http.Header, key, want string) {
	t.Helper()

	got := headers.Get(key)
	if got != want {
		t.Fatalf("expected header %q to be %q, got %q", key, want, got)
	}
}

func assertContains(t *testing.T, values []string, want string) {
	t.Helper()

	if slices.Contains(values, want) {
		return
	}

	t.Fatalf("expected value %q in %v", want, values)
}
