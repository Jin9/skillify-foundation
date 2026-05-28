package wrapper

// Shared response codes used across all aggregates.
const (
	CodeSuccess    Code    = "0000"
	MessageSuccess Message = "Success"

	CodeBadRequest    Code    = "HM400"
	MessageBadRequest Message = "Bad Request"

	CodeUnauthorized    Code    = "HM401"
	MessageUnauthorized Message = "Unauthorized"

	CodeForbidden    Code    = "HM403"
	MessageForbidden Message = "Forbidden"

	CodeNotFound    Code    = "HM404"
	MessageNotFound Message = "Not Found"

	CodeInternalError    Code    = "HM500"
	MessageInternalError Message = "Internal Server Error"

	CodeServiceUnavail    Code    = "HM503"
	MessageServiceUnavail Message = "Service Unavailable"

	CodeTimeout    Code    = "HM504"
	MessageTimeout Message = "Gateway Timeout"
)
