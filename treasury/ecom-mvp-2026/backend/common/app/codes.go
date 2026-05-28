package app

import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

// Code and Message are type aliases so that all packages can continue to use
// app.Code / app.Message while wrapper stays free of domain dependencies.
type Code = wrapper.Code
type Message = wrapper.Message
type Response[T any] = wrapper.Response[T]

// Standard response code re-exports for ergonomic use via app.CodeXxx / app.MessageXxx.
// The canonical definitions live in the wrapper package.
const (
	CodeSuccess           = wrapper.CodeSuccess
	MessageSuccess        = wrapper.MessageSuccess
	CodeBadRequest        = wrapper.CodeBadRequest
	MessageBadRequest     = wrapper.MessageBadRequest
	CodeUnauthorized      = wrapper.CodeUnauthorized
	MessageUnauthorized   = wrapper.MessageUnauthorized
	CodeForbidden         = wrapper.CodeForbidden
	MessageForbidden      = wrapper.MessageForbidden
	CodeNotFound          = wrapper.CodeNotFound
	MessageNotFound       = wrapper.MessageNotFound
	CodeInternalError     = wrapper.CodeInternalError
	MessageInternalError  = wrapper.MessageInternalError
	CodeServiceUnavail    = wrapper.CodeServiceUnavail
	MessageServiceUnavail = wrapper.MessageServiceUnavail
	CodeTimeout           = wrapper.CodeTimeout
	MessageTimeout        = wrapper.MessageTimeout
)

// Domain-specific codes — add aggregate codes below.
const (
	// Product aggregate codes
	CodeProductCreated    Code    = "PD201"
	MessageProductCreated Message = "Product Created"

	CodeProductNotFound    Code    = "PD404"
	MessageProductNotFound Message = "Product Not Found"

	CodeProductDeleted    Code    = "PD200"
	MessageProductDeleted Message = "Product Deleted"
)
