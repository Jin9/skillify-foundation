package wrapper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupGinTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return c, w
}

func TestRespond_WithData(t *testing.T) {
	c, w := setupGinTestContext()

	type TestData struct {
		ID string `json:"id"`
	}

	data := TestData{ID: "123"}

	Respond(c, ResponseOption[TestData]{
		HTTPStatus: http.StatusOK,
		Code:       Code("MA0000"),
		Message:    Message("success"),
		TraceID:    "trace-001",
		Data:       &data,
	})

	// ---- HTTP assertions ----
	require.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	require.Equal(t, "MA0000", body["code"])
	require.Equal(t, "success", body["message"])
	require.NotNil(t, body["data"])

	// ---- Context meta assertions ----
	meta, ok := c.Get(CtxResponseMeta)
	require.True(t, ok)

	rm := meta.(ResponseMeta)
	require.Equal(t, http.StatusOK, rm.HTTPStatus)
	require.Equal(t, Code("MA0000"), rm.Code)
	require.Equal(t, "trace-001", rm.TraceID)
}
