package checkout

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// UpstreamError is returned by all HTTP clients when the downstream service
// returns a non-2xx, non-retryable status.
type UpstreamError struct {
	Service    string
	HTTPStatus int
	Code       string
	Message    string
}

func (e *UpstreamError) Error() string {
	return fmt.Sprintf("%s: HTTP %d code=%s msg=%s", e.Service, e.HTTPStatus, e.Code, e.Message)
}

// IsUpstreamError unwraps an *UpstreamError if present.
func IsUpstreamError(err error, code string) bool {
	var u *UpstreamError
	if errors.As(err, &u) {
		return u.Code == code
	}
	return false
}

// isRetryableErr returns true for connection-level errors (not HTTP 4xx/5xx).
func isRetryableErr(err error) bool {
	if err == nil {
		return false
	}
	// context errors are not retryable
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	// net/http errors (connection refused, timeout dial) are retryable
	return true
}

// bearerAuthOption sets the Authorization header on an outgoing request.
func bearerAuthOption(token string) httpclient.OptionFunc {
	return func(req *http.Request, _ ...context.Context) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// internalSecretOption sets the X-Internal-Secret header for service-to-service calls.
func internalSecretOption(secret string) httpclient.OptionFunc {
	return func(req *http.Request, _ ...context.Context) {
		req.Header.Set("X-Internal-Secret", secret)
	}
}
