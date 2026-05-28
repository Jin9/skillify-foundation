package middleware

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

var logAttrPool = sync.Pool{
	New: func() any {
		attrs := make([]any, 0, 16)
		return &attrs
	},
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		latency := time.Since(start)
		latencyMs := latency.Milliseconds()

		method := c.Request.Method
		path := c.FullPath()
		status := c.Writer.Status()

		meta, ok := c.Get(wrapper.CtxResponseMeta)
		if !ok {
			// fallback (should be rare)
			slog.WarnContext(c.Request.Context(), "api.request.completed_no_meta",
				slog.String("event", "api.request.completed_no_meta"),
				slog.String("component", "access_log"),
				slog.String("http_method", method),
				slog.String("http_request_path", path),
				slog.Int("http_status", status),
				slog.Int64("latency_ms", latencyMs),
			)
			return
		}

		responseMeta := meta.(wrapper.ResponseMeta)

		// Borrow a pre-allocated slice from the pool
		attrsPtr := logAttrPool.Get().(*[]any)
		logAttrs := (*attrsPtr)[:0] // Reset length to 0, retaining capacity

		logAttrs = append(logAttrs,
			slog.String("event", "api.request.completed"),
			slog.String("component", "access_log"),
			slog.String("http_method", method),
			slog.String("http_request_path", path),
			slog.Int("http_status", status),
			slog.String("code", string(responseMeta.Code)),
			slog.String("trace_id", responseMeta.TraceID),
			slog.Int64("latency_ms", latencyMs),
		)

		if responseMeta.Err == nil {
			slog.InfoContext(c.Request.Context(), "api.request.completed", logAttrs...)
			logAttrPool.Put(attrsPtr)
			return
		}

		logAttrs = appendErrorAttrs(logAttrs, responseMeta.Err)
		slog.ErrorContext(c.Request.Context(), "api.request.completed", logAttrs...)
		logAttrPool.Put(attrsPtr)
	}
}

// appendErrorAttrs enriches log attributes with error details extracted from
// the serror package: the decoded message, source location, and any
// structured context attached via SError.With().
func appendErrorAttrs(attrs []any, err error) []any {
	msg, sourceAttrs := serror.DecodeMessage(err.Error())
	attrs = append(attrs, slog.String("error", msg))

	if len(sourceAttrs) > 0 {
		attrs = append(attrs, attrsToGroup("error_source", sourceAttrs))
	}

	var sErr *serror.SError
	if errors.As(err, &sErr) {
		if ctxAttrs := sErr.Attrs(); len(ctxAttrs) > 0 {
			attrs = append(attrs, attrsToGroup("error_context", ctxAttrs))
		}
	}

	return attrs
}

// attrsToGroup converts []slog.Attr into a single slog.Group attribute.
func attrsToGroup(name string, attrs []slog.Attr) slog.Attr {
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return slog.Group(name, args...)
}
