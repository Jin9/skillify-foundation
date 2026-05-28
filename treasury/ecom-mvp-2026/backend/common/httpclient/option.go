package httpclient

import (
	"context"
	"fmt"
	"net/http"
)

type debugContextKey struct{}

// DebugOption enables or disables debug logging for the request lifecycle.
// When enabled, request and response bodies are read for logging.
func DebugOption(enabled bool) OptionFunc {
	return func(r *http.Request, _ ...context.Context) {
		ctx := context.WithValue(r.Context(), debugContextKey{}, enabled)
		*r = *r.WithContext(ctx)
	}
}

// RequestOption adapts a simple request-mutating function into an OptionFunc.
// Useful for wrapping existing option helpers that don't need context.
func RequestOption(fn func(r *http.Request)) OptionFunc {
	return func(r *http.Request, _ ...context.Context) {
		fn(r)
	}
}

// HeaderOption sets a request header.
func HeaderOption(key, value string) OptionFunc {
	return func(r *http.Request, _ ...context.Context) {
		r.Header.Set(key, value)
	}
}

// BearerTokenOption sets Authorization header as a Bearer token.
func BearerTokenOption(token string) OptionFunc {
	return HeaderOption("Authorization", fmt.Sprintf("Bearer %s", token))
}

func jsonOption(r *http.Request, ctx ...context.Context) {
	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json")
	}
}
