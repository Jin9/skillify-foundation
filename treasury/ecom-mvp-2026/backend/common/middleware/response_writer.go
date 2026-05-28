package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"

	"github.com/gin-gonic/gin"
)

// responseWriter wraps gin.ResponseWriter to intercept non-OK responses for logging.
type responseWriter struct {
	gin.ResponseWriter
	ctx         context.Context
	statusCode  int
	meta        map[string]string
	successCode app.Code
}

func newResponseWriter(w gin.ResponseWriter, ctx context.Context, meta map[string]string, successCode app.Code) *responseWriter {
	if ctx == nil {
		ctx = context.Background()
	}
	return &responseWriter{ResponseWriter: w, ctx: ctx, statusCode: http.StatusOK, meta: meta, successCode: successCode}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	ctx := w.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	if w.statusCode != http.StatusOK {
		var resp app.Response[any]
		if err := json.Unmarshal(b, &resp); err != nil {
			slog.DebugContext(ctx, "response json decode failed", slog.Any("error", err))
		}

		if (resp.Code == "" || resp.Code == w.successCode) && resp.Message == "" && resp.Data == nil {
			slog.WarnContext(ctx, "response not standard")
		}

		origMessage := string(resp.Message)
		msg, attrs := serror.DecodeMessage(origMessage)

		attrs = append(attrs, slog.Int("status-code", w.statusCode))
		for k, v := range w.meta {
			attrs = append(attrs, slog.String(k, v))
		}
		if resp.Code != "" && resp.Code != w.successCode {
			attrs = append(attrs, slog.String("code", string(resp.Code)))
		}

		level := logLevelFromStatus(w.statusCode)
		slog.LogAttrs(ctx, level, msg, attrs...)

		if origMessage != msg {
			var out bytes.Buffer
			json.HTMLEscape(&out, []byte(origMessage))
			b = bytes.Replace(b, out.Bytes(), []byte(msg), 1)
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func logLevelFromStatus(statusCode int) slog.Level {
	if statusCode >= http.StatusInternalServerError {
		return slog.LevelError
	}
	if statusCode >= http.StatusBadRequest {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func AutoLoggingMiddleware(successCode app.Code) gin.HandlerFunc {
	return func(c *gin.Context) {
		meta := map[string]string{}
		if ref, ok := c.Request.Context().Value(refIDKey).(string); ok {
			meta[string(refIDKey)] = ref
		}

		c.Writer = newResponseWriter(c.Writer, c.Request.Context(), meta, successCode)
		c.Next()
	}
}
