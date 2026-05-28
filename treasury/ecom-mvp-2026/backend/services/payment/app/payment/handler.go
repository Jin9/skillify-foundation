package payment

import (
	"encoding/json"
	"fmt"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// handler holds the dependencies shared across all HTTP handlers.
type handler struct {
	svc        Service
	hmacSecret []byte // CALLBACK_HMAC_SECRET raw bytes
}

// HandlerConfig is passed to NewHandler.
type HandlerConfig struct {
	Service    Service
	HMACSecret []byte
}

// NewHandler constructs the handler.
func NewHandler(cfg HandlerConfig) *handler {
	return &handler{
		svc:        cfg.Service,
		hmacSecret: cfg.HMACSecret,
	}
}

// codeOf / msgOf convert plain strings to the wrapper's typed Code / Message.
// They keep handler code concise without repeating wrapper.Code("...") everywhere.
func codeOf(s string) wrapper.Code   { return wrapper.Code(s) }
func msgOf(s string) wrapper.Message { return wrapper.Message(s) }

// bindRawJSON parses rawBody into dst and validates required fields.
// Used by handler_callback.go which must read the raw body for HMAC verification
// before gin's ShouldBindJSON can be called (body is already consumed).
func bindRawJSON(rawBody []byte, dst any) error {
	if err := json.Unmarshal(rawBody, dst); err != nil {
		return fmt.Errorf("json parse: %w", err)
	}
	// Basic non-nil validation is deferred to ProcessCallbackInput construction.
	return nil
}
