package wrapper

import (
	"github.com/gin-gonic/gin"
)

const CtxResponseMeta = "response_meta"

// CtxTraceID is the gin-context key where the trace/ref id is stored by the
// observability middleware. Respond falls back to this value when
// ResponseOption.TraceID is empty so handlers don't need to plumb the trace id
// through every call site.
const CtxTraceID = "trace_id"

// Code is a short application-level result code included in every response.
type Code string

// Message is a human-readable description included in every response.
type Message string

// Response is the standard JSON envelope written to the wire.
type Response[T any] struct {
	Code    Code    `json:"code"`
	Message Message `json:"message"`
	Data    *T      `json:"data,omitempty"`
	TraceID string  `json:"traceId,omitempty"`
}

type ResponseMeta struct {
	HTTPStatus int
	Code       Code
	TraceID    string
	Err        error
}

type ResponseOption[T any] struct {
	HTTPStatus int
	Code       Code
	Message    Message
	Data       *T

	TraceID string
	Err     error
}

func Respond[T any](c *gin.Context, opt ResponseOption[T]) {
	traceID := opt.TraceID
	if traceID == "" {
		traceID = c.GetString(CtxTraceID)
	}

	c.Set(CtxResponseMeta, ResponseMeta{
		HTTPStatus: opt.HTTPStatus,
		Code:       opt.Code,
		TraceID:    traceID,
		Err:        opt.Err,
	})

	c.JSON(opt.HTTPStatus, Response[T]{
		Code:    opt.Code,
		Message: opt.Message,
		Data:    opt.Data, // nil → omitted
		TraceID: traceID,
	})
}
