package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

func TestAccessLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		hasMeta      bool
		err          error
		expectedCode int
	}{
		{
			name:         "Success log with meta",
			hasMeta:      true,
			err:          nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Error log with plain error",
			hasMeta:      true,
			err:          errors.New("plain error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error log with serror.Wrap (includes error_source)",
			hasMeta:      true,
			err:          serror.Wrap(errors.New("wrapped error")),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error log with serror.New (includes error_source)",
			hasMeta:      true,
			err:          serror.New("serror new error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Fallback log without meta",
			hasMeta:      false,
			err:          nil,
			expectedCode: http.StatusOK,
		},
		{
			name:    "Error log with serror.Wrap + With context attrs (includes error_context)",
			hasMeta: true,
			err: serror.Wrap(errors.New("db timeout")).With(
				slog.String("query", "SELECT * FROM users"),
				slog.String("db_host", "db-replica-01"),
			),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error log with serror.New without context attrs (no error_context group)",
			hasMeta:      true,
			err:          serror.New("validation failed"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(AccessLog())

			r.GET("/test", func(c *gin.Context) {
				if tt.hasMeta {
					meta := wrapper.ResponseMeta{
						Code:    "0",
						TraceID: "test-trace-id",
					}
					if tt.err != nil {
						meta.Err = tt.err
						meta.HTTPStatus = tt.expectedCode
					}
					c.Set(wrapper.CtxResponseMeta, meta)
				}
				c.Status(tt.expectedCode)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}
