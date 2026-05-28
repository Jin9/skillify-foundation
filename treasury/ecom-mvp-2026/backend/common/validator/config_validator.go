// Package validator is a thin convenience wrapper around go-playground/validator.
// It exists so that callers only need to import one package for struct validation
// without caring about the underlying library version.
package validator

import (
	"log"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

var logFatal = log.Fatal

// Validate validates config and returns an error instead of exiting the process.
// Use this in libraries to avoid log.Fatal in constructors.
func Validate(i any) error {
	return validator.New().Struct(i)
}

// Deprecated: MustValid calls log.Fatal on validation failure.
// Use Validate instead and handle the error explicitly.
func MustValid(i any) {
	if err := Validate(i); err != nil {
		logFatal(err.Error(), slog.String("tag", "init configs"))
	}
}
