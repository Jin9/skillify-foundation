package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/generator"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

const traceparentHeaderKey = "traceparent"

// Type aliases for backwards compatibility.
type TraceID = generator.TraceID
type SpanID = generator.SpanID
type TraceFlags = generator.TraceFlags
type TraceParent = generator.TraceParent

// Re-export constants and sentinel errors from generator.
const TraceFlagsSampled = generator.TraceFlagsSampled

var (
	ErrEmptyTraceParent = generator.ErrEmptyTraceParent
	ErrWrongTraceParent = generator.ErrWrongTraceParent
	ErrInvalidTraceID   = generator.ErrInvalidTraceID
	ErrInvalidSpanID    = generator.ErrInvalidSpanID
)

// NewTraceParent delegates to generator.NewTraceParent.
func NewTraceParent() TraceParent {
	return generator.NewTraceParent()
}

// Parse delegates to generator.Parse.
func Parse(parent string) (TraceParent, error) {
	return generator.Parse(parent)
}

func TraceContextTraceIDMiddleware(headerKey string) gin.HandlerFunc {
	hk := headerKey
	if hk == "" {
		hk = traceparentHeaderKey
	}

	return func(c *gin.Context) {
		if refID, ok := c.Request.Context().Value(refIDKey).(string); ok && refID != "" {
			c.Next()
			return
		}

		parent := c.Request.Header.Get(hk)
		if parent == "" {
			slog.DebugContext(c.Request.Context(), "missing traceparent", slog.String("header-key", hk))
			refID := NewTraceParent().TraceID.String()
			c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), refID))
			c.Set(wrapper.CtxTraceID, refID)
			c.Next()
			return
		}

		tp, err := Parse(parent)
		if err != nil {
			slog.DebugContext(c.Request.Context(), "invalid traceparent", slog.String("header-key", hk), slog.Any("error", err))
			refID := NewTraceParent().TraceID.String()
			c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), refID))
			c.Set(wrapper.CtxTraceID, refID)
			c.Next()
			return
		}

		refID := tp.TraceID.String()
		c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), refID))
		c.Set(wrapper.CtxTraceID, refID)
		c.Next()
	}
}
