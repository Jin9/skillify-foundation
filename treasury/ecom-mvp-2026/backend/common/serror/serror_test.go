package serror

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWrap(t *testing.T) {
	given := "testing"
	var pc uintptr
	var file string
	var line int

	var first = func() error {
		pc, file, line, _ = runtime.Caller(0)
		return Wrap(errors.New(given))
	}

	err := first()

	expect := fmt.Sprintf("((testing+%s:%d:%s))", filepath.Base(file), line+1, filepath.Base(runtime.FuncForPC(pc).Name()))

	if err.Error() != expect {
		t.Errorf("expect:\n%q\n\nactual:\n%q\n", expect, err.Error())
	}
}

func TestWrapSkipZeroStep_ErrorOccurredFromCaller(t *testing.T) {
	given := "testing"
	var pc uintptr
	var file string
	var line int

	var first = func() error {
		pc, file, line, _ = runtime.Caller(0)
		return WrapSkip(errors.New(given), 0)
	}

	err := first()

	expect := fmt.Sprintf("((testing+%s:%d:%s))", filepath.Base(file), line+1, filepath.Base(runtime.FuncForPC(pc).Name()))

	if err.Error() != expect {
		t.Errorf("expect:\n%q\n\nactual:\n%q\n", expect, err.Error())
	}
}

func TestWrapSkipOneStepBack_ErrorOccurredFromCallerofCaller(t *testing.T) {
	given := "testing"
	var pc uintptr
	var file string
	var line int

	var second = func() error {
		return WrapSkip(errors.New(given), 1)
	}

	var first = func() error {
		pc, file, line, _ = runtime.Caller(0)
		return second()
	}

	err := first()

	expect := fmt.Sprintf("((testing+%s:%d:%s))", filepath.Base(file), line+1, filepath.Base(runtime.FuncForPC(pc).Name()))

	if err.Error() != expect {
		t.Errorf("expect:\n%q\n\nactual:\n%q\n", expect, err.Error())
	}
}

func TestCallerNotRecoverSkip(t *testing.T) {
	const notRecover = 4
	s := caller(notRecover)
	if s != "" {
		t.Error("the origin not recover skip is 4")
	}
}

func TestDecodeMessageEmptyString(t *testing.T) {
	msg, atts := DecodeMessage("")

	if msg != "" {
		t.Errorf("given empty string to Decode expect empty string msg but actual %q\n", msg)
	}

	if len(atts) != 0 {
		t.Errorf("given empty string to DecodeMessage expect 0 lenght Attrs but actual %d\n", len(atts))
	}
}

func TestSErrorDecode(t *testing.T) {
	given := "test message"
	err := New(given)
	msg, attrs := DecodeMessage(err.Error())

	if msg != "test message" {
		t.Errorf("%q message is expected but actual %q", given, msg)
	}

	if len(attrs) != 3 {
		t.Errorf("attrs should have 3 elements but actual %d", len(attrs))
	}
}

func TestSErrorWithMessageDecode(t *testing.T) {
	given := "test message"
	err := New(given)

	message := fmt.Sprintf("first message: %s", err)

	msg, attrs := DecodeMessage(message)

	if msg != "first message: test message" {
		t.Errorf("message: %q\n", message)
		t.Errorf("%q message is expected but actual %q\n", "first message: "+given, msg)
	}

	if len(attrs) != 3 {
		t.Errorf("attrs should have 3 elements but actual %d", len(attrs))
	}
}

func TestDecodeMessagePlainError(t *testing.T) {
	given := "test message"
	err := errors.New(given)
	msg, attrs := DecodeMessage(err.Error())

	if msg != "test message" {
		t.Errorf("%q message is expected but actual %q", given, msg)
	}

	if len(attrs) != 0 {
		t.Errorf("attrs should have 3 elements but actual %d", len(attrs))
	}
}

func TestDecodeMessageWrongPattern(t *testing.T) {
	givenMsg := prefix + "message" + suffix
	msg, atts := DecodeMessage(givenMsg)

	if givenMsg != msg {
		t.Errorf("expect %q string msg but actual %q\n", givenMsg, msg)
	}

	if len(atts) != 0 {
		t.Errorf("given wrong pattern to DecodeMessage expect 0 lenght Attrs but actual %d\n", len(atts))
	}
}

func TestDecodeMessageWrongSourcePattern(t *testing.T) {
	givenMsg := prefix + "message" + separator + "filename.go" + suffix
	msg, atts := DecodeMessage(givenMsg)

	if givenMsg != msg {
		t.Errorf("expect %q string msg but actual %q\n", givenMsg, msg)
	}

	if len(atts) != 0 {
		t.Errorf("given wrong pattern to DecodeMessage expect 0 lenght Attrs but actual %d\n", len(atts))
	}
}

func TestAttrsEmptyByDefault(t *testing.T) {
	err := Wrap(errors.New("test"))
	attrs := err.Attrs()

	if len(attrs) != 0 {
		t.Errorf("expected 0 attrs by default but got %d", len(attrs))
	}
}

func TestWrapWith(t *testing.T) {
	err := Wrap(errors.New("test")).With(
		slog.String("member_id", "m-123"),
		slog.String("organization_id", "o-456"),
	)

	attrs := err.Attrs()
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attrs but got %d", len(attrs))
	}

	if attrs[0].Key != "member_id" || attrs[0].Value.String() != "m-123" {
		t.Errorf("first attr: got key=%q value=%q", attrs[0].Key, attrs[0].Value.String())
	}

	if attrs[1].Key != "organization_id" || attrs[1].Value.String() != "o-456" {
		t.Errorf("second attr: got key=%q value=%q", attrs[1].Key, attrs[1].Value.String())
	}
}

func TestWithChaining(t *testing.T) {
	err := Wrap(errors.New("test")).
		With(slog.String("a", "1")).
		With(slog.String("b", "2"))

	attrs := err.Attrs()
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attrs after chaining but got %d", len(attrs))
	}

	if attrs[0].Key != "a" {
		t.Errorf("expected key 'a' but got %q", attrs[0].Key)
	}
	if attrs[1].Key != "b" {
		t.Errorf("expected key 'b' but got %q", attrs[1].Key)
	}
}

func TestNewWith(t *testing.T) {
	err := New("test error").With(slog.String("handler", "TestHandler"))

	attrs := err.Attrs()
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attr but got %d", len(attrs))
	}

	if attrs[0].Key != "handler" || attrs[0].Value.String() != "TestHandler" {
		t.Errorf("attr: got key=%q value=%q", attrs[0].Key, attrs[0].Value.String())
	}
}

func TestWrapSkipWith(t *testing.T) {
	err := WrapSkip(errors.New("test"), 0).With(slog.String("key", "val"))

	attrs := err.Attrs()
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attr but got %d", len(attrs))
	}

	if attrs[0].Key != "key" {
		t.Errorf("expected key 'key' but got %q", attrs[0].Key)
	}
}

func TestWithDoesNotAffectErrorMessage(t *testing.T) {
	plain := Wrap(errors.New("original"))
	withAttrs := Wrap(errors.New("original")).With(slog.String("k", "v"))

	// Error() should contain the error message regardless of attrs
	if plain.Error() == "" {
		t.Error("plain error should not be empty")
	}

	// With() should not change the Error() format
	_, plainAttrs := DecodeMessage(plain.Error())
	_, withMsgAttrs := DecodeMessage(withAttrs.Error())

	if len(plainAttrs) != len(withMsgAttrs) {
		t.Error("With() should not affect DecodeMessage output")
	}
}

func TestUnwrap(t *testing.T) {
	orig := errors.New("original error")
	err := Wrap(orig)

	if !errors.Is(err, orig) {
		t.Error("Unwrap should allow errors.Is to match the original error")
	}
	if err.Unwrap() != orig {
		t.Error("Unwrap should return exact original error")
	}
}

func TestWrapSkipNegative(t *testing.T) {
	orig := errors.New("test")
	err := WrapSkip(orig, -5)

	if err == nil {
		t.Fatal("WrapSkip should return a valid error even with negative skip")
	}
	// The negative skip should be clamped to wrapSkipBase
	// We can check it has the same caller properties roughly as 0
	msg, attrs := DecodeMessage(err.Error())
	if msg != "test" {
		t.Errorf("expected 'test', got %q", msg)
	}
	if len(attrs) == 0 {
		t.Error("expected caller info to be attached")
	}
}
