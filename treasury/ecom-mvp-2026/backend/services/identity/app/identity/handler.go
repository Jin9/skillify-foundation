package identity

import "errors"

// Sentinel errors returned by Service and mapped to HTTP codes in handlers.
var (
	ErrEmailAlreadyRegistered = errors.New("EMAIL_ALREADY_REGISTERED")
	ErrAuthInvalid            = errors.New("AUTH_INVALID")
	ErrAuthSuspended          = errors.New("AUTH_SUSPENDED")
	ErrAuthRevoked            = errors.New("AUTH_REVOKED")
	ErrUserNotFound           = errors.New("USER_NOT_FOUND")
	ErrAddressNotFound        = errors.New("ADDRESS_NOT_FOUND")
)

// HandlerConfig wires service dependencies into the handler set.
type HandlerConfig struct {
	Service Service
}

type handler struct {
	svc Service
}

func NewHandler(cfg HandlerConfig) *handler {
	return &handler{svc: cfg.Service}
}
