package validator

import (
	"testing"
)

type DummyConfig struct {
	Name string `validate:"required"`
}

func TestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cfg := DummyConfig{Name: "test"}
		if err := Validate(&cfg); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		cfg := DummyConfig{}
		if err := Validate(&cfg); err == nil {
			t.Errorf("expected error for empty name")
		}
	})
}

func TestMustValid(t *testing.T) {
	orig := logFatal
	defer func() { logFatal = orig }()

	t.Run("valid", func(t *testing.T) {
		fatalCalled := false
		logFatal = func(v ...any) {
			fatalCalled = true
		}

		cfg := DummyConfig{Name: "test"}
		MustValid(&cfg)
		if fatalCalled {
			t.Error("logFatal should not be called for valid config")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		fatalCalled := false
		logFatal = func(v ...any) {
			fatalCalled = true
		}

		cfg := DummyConfig{}
		MustValid(&cfg)
		if !fatalCalled {
			t.Error("logFatal should be called for invalid config")
		}
	})
}
