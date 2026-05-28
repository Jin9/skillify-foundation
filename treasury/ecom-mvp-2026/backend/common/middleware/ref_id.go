package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

type ctxKey string

const refIDKey ctxKey = "ref-id"

const DefaultRefIDHeaderKey = "X-REF-ID"

func RefIDMiddleware(headerKey string) gin.HandlerFunc {
	hk := headerKey
	if hk == "" {
		hk = DefaultRefIDHeaderKey
	}

	return func(c *gin.Context) {
		refID := c.Request.Header.Get(hk)
		if refID == "" {
			slog.DebugContext(c.Request.Context(), "missing ref-id", slog.String("header-key", hk))
			refID = uuid.NewString()
		}

		c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), refID))
		c.Set(wrapper.CtxTraceID, refID)
		c.Next()
	}
}

func ForwardRefIDOption(r *http.Request, ctxs ...context.Context) {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	if refID, ok := ctx.Value(refIDKey).(string); ok {
		if r.Header.Get(DefaultRefIDHeaderKey) == "" {
			r.Header.Set(DefaultRefIDHeaderKey, refID)
		}
		if r.Header.Get(string(refIDKey)) == "" {
			r.Header.Set(string(refIDKey), refID)
		}
	}
}

func newRefIDContext(ctx context.Context, refID string) context.Context {
	return context.WithValue(ctx, refIDKey, refID)
}

func SetRefID(c *gin.Context, refID string) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), refIDKey, refID))
}

func RefID(c *gin.Context) string {
	if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
		return v
	}
	return ""
}
