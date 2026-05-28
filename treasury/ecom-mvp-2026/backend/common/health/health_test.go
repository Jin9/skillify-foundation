package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLiveness(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler := Liveness("v1.0.0", "commit-abc")
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "hostname")
		assert.Contains(t, w.Body.String(), "v1.0.0")
		assert.Contains(t, w.Body.String(), "commit-abc")
	})

	t.Run("hostname error", func(t *testing.T) {
		orig := osHostname
		defer func() { osHostname = orig }()
		osHostname = func() (string, error) {
			return "", errors.New("mock error")
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler := Liveness("v1.0.0", "commit-abc")
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mock error")
	})
}

func TestReadiness(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := Readiness()
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := Metrics()
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "memory")
	assert.Contains(t, body, "alloc")
	assert.Contains(t, body, "heapInuse")
}
