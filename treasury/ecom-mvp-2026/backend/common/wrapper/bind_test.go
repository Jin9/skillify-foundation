package wrapper_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sampleRequest struct {
	Name string `json:"name" binding:"required"`
}

func TestBindJSON_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(sampleRequest{Name: "test"})
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	req, ok := wrapper.BindJSON[sampleRequest](c)
	require.True(t, ok)
	assert.Equal(t, "test", req.Name)
}

func TestBindJSON_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":""}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	req, ok := wrapper.BindJSON[sampleRequest](c)
	assert.False(t, ok)
	assert.Nil(t, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBindJSON_WithAttrs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	req, ok := wrapper.BindJSON[sampleRequest](c, slog.String("handler", "TestHandler"))
	assert.False(t, ok)
	assert.Nil(t, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
