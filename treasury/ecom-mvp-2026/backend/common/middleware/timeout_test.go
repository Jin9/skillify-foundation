package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	timeout := 100 * time.Millisecond

	r := gin.New()
	r.Use(Timeout(timeout))

	var hasDeadline bool
	var deadline time.Time

	r.GET("/test-timeout", func(c *gin.Context) {
		ctx := c.Request.Context()
		deadline, hasDeadline = ctx.Deadline()
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test-timeout", nil)
	w := httptest.NewRecorder()

	startTime := time.Now()
	r.ServeHTTP(w, req)

	if !hasDeadline {
		t.Fatal("expected request context to have a deadline")
	}

	expectedDeadline := startTime.Add(timeout)
	diff := deadline.Sub(expectedDeadline)
	if diff < -10*time.Millisecond || diff > 10*time.Millisecond {
		t.Errorf("expected deadline around %v, got %v (diff: %v)", expectedDeadline, deadline, diff)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
