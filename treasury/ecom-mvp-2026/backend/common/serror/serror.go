// Package serror provides structured error wrapping with source-location tracking.
package serror

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

// SError represents error structure with
// line of code that program had been executing when an error occurred
type SError struct {
	err   error
	at    string
	attrs []slog.Attr
}

// With attaches structured context fields to the error for investigation.
// Fields flow through the observability middleware and are auto-censored
// by the registered CensorReplacer (e.g. password, api_key are masked).
func (err *SError) With(attrs ...slog.Attr) *SError {
	err.attrs = append(err.attrs, attrs...)
	return err
}

// Attrs returns the structured context fields attached to this error.
func (err *SError) Attrs() []slog.Attr {
	return err.attrs
}

const (
	prefix          = "(("
	suffix          = "))"
	separator       = "+"
	sourceSeparator = ":"
	prefixSize      = len(prefix)
	suffixSize      = len(suffix)

	wrapSkipBase = 2
)

// Error returns error message
func (err *SError) Error() string {
	return fmt.Sprintf("((%s+%s))", err.err, err.at)
}

// Unwrap enables errors.Is / errors.As
func (err *SError) Unwrap() error {
	return err.err
}

// New returns [SError]
func New(s string) *SError {
	return &SError{
		err: errors.New(s),
		at:  caller(wrapSkipBase),
	}
}

// WrapSkip wraps the previous error and includes the line of code the occurred error
// with skip caller given step
//
// Example: if function A caught an error
//
//	func A() {
//	  err := occurred()
//	  err = serror.WrapSkip(err, 0)
//	}
//
// Example: if function A call function B and B caught an error
// but you need to keep line of code in function A
//
//	func A() {
//	  err := B()
//	}
//
//	func B() {
//	  err := occurred()
//	  err = serror.WrapSkip(err, 1)
//	}
func WrapSkip(err error, skip int) *SError {
	skip += wrapSkipBase
	if skip < wrapSkipBase {
		skip = wrapSkipBase
	}
	return &SError{
		err: err,
		at:  caller(skip),
	}
}

func Wrap(err error) *SError {
	return &SError{
		err: err,
		at:  caller(wrapSkipBase),
	}
}

func caller(skip int) string {
	pc, file, no, ok := runtime.Caller(skip)
	if ok {
		b := filepath.Base(file)
		f := filepath.Base(runtime.FuncForPC(pc).Name())
		return fmt.Sprintf("%s:%d:%s", b, no, f)
	}
	return ""
}

// DecodeMessage decodes error message that was generated from serror (New,Wrap or WrapSkip)
// to plain message and []slog.Attr
//
// Example: DecodeMessage("!!original error message:@handler.go:23:handler.Serve")
//
//	 message = "original error message"
//		slog.Attr[0] = slog.String("file", "handler.go")
//		slog.Attr[1] = slog.String("line", "23")
//		slog.Attr[2] = slog.String("func", "handler.Serve")
func DecodeMessage(s string) (msg string, attrs []slog.Attr) {
	if s == "" {
		return "", []slog.Attr{}
	}

	serrorFrom := strings.Index(s, prefix)
	serrorTo := strings.LastIndex(s, suffix)

	if serrorFrom == -1 || serrorTo == -1 || serrorFrom >= serrorTo {
		return s, []slog.Attr{}
	}

	serrorMsg := s[serrorFrom+prefixSize : serrorTo]

	elem := strings.Split(serrorMsg, separator)
	if len(elem) != 2 {
		return s, []slog.Attr{}
	}

	sources := strings.Split(elem[1], sourceSeparator)
	if len(sources) != 3 {
		return s, []slog.Attr{}
	}

	s = strings.Replace(s, s[serrorFrom:serrorTo+suffixSize], elem[0], 1)

	return s, []slog.Attr{
		slog.String("file", sources[0]),
		slog.String("line", sources[1]),
		slog.String("func", sources[2]),
	}
}
