package logger_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"
)

func init() {
	// Force non-LOCAL env so logger.New produces JSON (deterministic output for examples).
	os.Unsetenv("ENV")
	logger.LogLevel = slog.LevelInfo
}

// ExampleCensorReplacer_basic demonstrates masking a single sensitive field.
//
// Any slog attribute whose key matches a registered censor key will have
// its value replaced with the configured mask string.
func ExampleCensorReplacer_basic() {
	logger.ResetCensors()
	buf := new(bytes.Buffer)
	logger.SetOutput(buf)
	defer logger.ResetOutput()

	logger.AddCensor("password", "***REDACTED***")

	l := logger.New(logger.CensorReplacer)

	l.Info("user login",
		slog.String("username", "john.doe"),
		slog.String("password", "s3cret!"),
	)

	var m map[string]any
	json.NewDecoder(buf).Decode(&m)

	fmt.Println("username:", m["username"])
	fmt.Println("password:", m["password"])
	// Output:
	// username: john.doe
	// password: ***REDACTED***
}

// ExampleCensorReplacer_multipleFields demonstrates masking multiple fields
// at once. Only registered keys are masked; all other attributes pass through untouched.
func ExampleCensorReplacer_multipleFields() {
	logger.ResetCensors()
	buf := new(bytes.Buffer)
	logger.SetOutput(buf)
	defer logger.ResetOutput()

	logger.LoadCensorsFromMap(map[string]string{
		"email":       "***EMAIL***",
		"card_number": "****-****-****-****",
		"cvv":         "***",
	})

	l := logger.New(logger.CensorReplacer)

	l.Info("payment processed",
		slog.String("order_id", "ORD-12345"),
		slog.String("email", "john@example.com"),
		slog.String("card_number", "4111111111111111"),
		slog.String("cvv", "123"),
		slog.Float64("amount", 1500.00),
	)

	var m map[string]any
	json.NewDecoder(buf).Decode(&m)

	fmt.Println("order_id:", m["order_id"])
	fmt.Println("email:", m["email"])
	fmt.Println("card_number:", m["card_number"])
	fmt.Println("cvv:", m["cvv"])
	fmt.Println("amount:", m["amount"])
	// Output:
	// order_id: ORD-12345
	// email: ***EMAIL***
	// card_number: ****-****-****-****
	// cvv: ***
	// amount: 1500
}

// ExampleCensorReplacer_withAWSKeyReplacer demonstrates composing CensorReplacer
// with other replacers. Replacers are evaluated in order — AWSKeyReplacer runs first,
// then CensorReplacer. This is how it is wired in main.go:
//
//	_ = logger.New(logger.AWSKeyReplacer, logger.CensorReplacer)
func ExampleCensorReplacer_withAWSKeyReplacer() {
	logger.ResetCensors()
	buf := new(bytes.Buffer)
	logger.SetOutput(buf)
	defer logger.ResetOutput()

	logger.AddCensor("api_key", "***API_KEY***")

	l := logger.New(logger.AWSKeyReplacer, logger.CensorReplacer)

	l.Info("calling external service",
		slog.String("service", "payment-gateway"),
		slog.String("api_key", "sk-live-abc123xyz"),
	)

	var m map[string]any
	json.NewDecoder(buf).Decode(&m)

	fmt.Println("service:", m["service"])
	fmt.Println("api_key:", m["api_key"])
	// Output:
	// service: payment-gateway
	// api_key: ***API_KEY***
}
