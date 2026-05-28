package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func TestAddCensor(t *testing.T) {
	// Clean up censors before test
	censors = map[string]string{}

	t.Run("add single censor", func(t *testing.T) {
		AddCensor("password", "***REDACTED***")

		if censors["password"] != "***REDACTED***" {
			t.Errorf("expected password to be masked with '***REDACTED***', got '%s'", censors["password"])
		}
	})

	t.Run("add multiple censors", func(t *testing.T) {
		censors = map[string]string{} // Reset

		AddCensor("api_key", "***API***")
		AddCensor("token", "***TOKEN***")
		AddCensor("secret", "***SECRET***")

		if len(censors) != 3 {
			t.Errorf("expected 3 censors, got %d", len(censors))
		}

		if censors["api_key"] != "***API***" {
			t.Errorf("api_key not set correctly")
		}
		if censors["token"] != "***TOKEN***" {
			t.Errorf("token not set correctly")
		}
		if censors["secret"] != "***SECRET***" {
			t.Errorf("secret not set correctly")
		}
	})

	t.Run("overwrite existing censor", func(t *testing.T) {
		censors = map[string]string{} // Reset

		AddCensor("password", "***OLD***")
		AddCensor("password", "***NEW***")

		if censors["password"] != "***NEW***" {
			t.Errorf("expected password to be overwritten with '***NEW***', got '%s'", censors["password"])
		}
	})
}

func TestLoadCensorsFromMap(t *testing.T) {
	t.Run("load censors from map", func(t *testing.T) {
		censors = map[string]string{} // Reset

		censorMap := map[string]string{
			"password":    "***REDACTED***",
			"card_number": "****-****-****-****",
			"cvv":         "***",
		}

		LoadCensorsFromMap(censorMap)

		if len(censors) != 3 {
			t.Errorf("expected 3 censors, got %d", len(censors))
		}

		for key, expectedValue := range censorMap {
			if censors[key] != expectedValue {
				t.Errorf("expected %s to be '%s', got '%s'", key, expectedValue, censors[key])
			}
		}
	})

	t.Run("load empty map", func(t *testing.T) {
		censors = map[string]string{} // Reset

		emptyMap := map[string]string{}
		LoadCensorsFromMap(emptyMap)

		if len(censors) != 0 {
			t.Errorf("expected 0 censors, got %d", len(censors))
		}
	})

	t.Run("merge with existing censors", func(t *testing.T) {
		censors = map[string]string{} // Reset

		// Add initial censors
		AddCensor("password", "***PASS***")
		AddCensor("api_key", "***API***")

		// Load additional censors
		newCensors := map[string]string{
			"token":  "***TOKEN***",
			"secret": "***SECRET***",
		}
		LoadCensorsFromMap(newCensors)

		if len(censors) != 4 {
			t.Errorf("expected 4 censors, got %d", len(censors))
		}

		// Check all censors exist
		if censors["password"] != "***PASS***" {
			t.Errorf("existing password censor was lost")
		}
		if censors["token"] != "***TOKEN***" {
			t.Errorf("new token censor was not added")
		}
	})
}

func TestCensorReplacer(t *testing.T) {
	t.Run("censor replacer with match", func(t *testing.T) {
		os.Unsetenv("ENV")

		LogLevel = slog.LevelInfo

		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		censors = map[string]string{} // Reset
		censors["cid"] = "xxxxxxxxxxxxx"

		l := New(CensorReplacer)

		l.Info("message", slog.String("cid", "123456789012"))

		var m map[string]string
		json.NewDecoder(buf).Decode(&m)
		v, ok := m["cid"]
		if !ok {
			t.Errorf("not found cid key\n")
		}

		if v != "xxxxxxxxxxxxx" {
			t.Errorf("%v\n", m)
			t.Errorf("replacer replace cid to %q: actual %q\n", "xxxxxxxxxxxxx", v)
		}
	})

	t.Run("censor replacer with no match", func(t *testing.T) {
		os.Unsetenv("ENV")

		LogLevel = slog.LevelInfo

		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		censors = map[string]string{} // Reset
		censors["password"] = "***REDACTED***"

		l := New(CensorReplacer)

		// Log a field that is NOT in censors map
		l.Info("message", slog.String("username", "john.doe"))

		var m map[string]string
		json.NewDecoder(buf).Decode(&m)
		username, ok := m["username"]
		if !ok {
			t.Errorf("not found username key\n")
		}

		// Should NOT be censored
		if username != "john.doe" {
			t.Errorf("username should not be censored, expected 'john.doe', got '%s'", username)
		}
	})

	t.Run("censor multiple fields", func(t *testing.T) {
		os.Unsetenv("ENV")

		LogLevel = slog.LevelInfo

		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		censors = map[string]string{} // Reset
		censors["password"] = "***PASS***"
		censors["api_key"] = "***API***"
		censors["token"] = "***TOKEN***"

		l := New(CensorReplacer)

		l.Info("sensitive data",
			slog.String("username", "john.doe"),
			slog.String("password", "secret123"),
			slog.String("api_key", "sk-12345"),
			slog.String("token", "Bearer xyz"),
		)

		var m map[string]string
		json.NewDecoder(buf).Decode(&m)

		// Check uncensored field
		if m["username"] != "john.doe" {
			t.Errorf("username should not be censored")
		}

		// Check censored fields
		if m["password"] != "***PASS***" {
			t.Errorf("password should be censored, got '%s'", m["password"])
		}
		if m["api_key"] != "***API***" {
			t.Errorf("api_key should be censored, got '%s'", m["api_key"])
		}
		if m["token"] != "***TOKEN***" {
			t.Errorf("token should be censored, got '%s'", m["token"])
		}
	})

	t.Run("censor replacer with empty censors map", func(t *testing.T) {
		os.Unsetenv("ENV")

		LogLevel = slog.LevelInfo

		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		censors = map[string]string{} // Reset - empty map

		l := New(CensorReplacer)

		l.Info("message", slog.String("password", "secret123"))

		var m map[string]string
		json.NewDecoder(buf).Decode(&m)

		// With empty censors map, nothing should be censored
		if m["password"] != "secret123" {
			t.Errorf("password should not be censored when censors map is empty, got '%s'", m["password"])
		}
	})

	t.Run("CensorReplacer function directly", func(t *testing.T) {
		censors = map[string]string{} // Reset
		censors["secret"] = "***MASKED***"

		// Test when key matches
		attr := slog.String("secret", "my-secret-value")
		result, replaced := CensorReplacer(nil, attr)

		if !replaced {
			t.Errorf("expected replaced to be true")
		}
		if result.Value.String() != "***MASKED***" {
			t.Errorf("expected masked value, got '%s'", result.Value.String())
		}

		// Test when key doesn't match
		attr2 := slog.String("public", "public-value")
		result2, replaced2 := CensorReplacer(nil, attr2)

		if replaced2 {
			t.Errorf("expected replaced to be false for non-matching key")
		}
		if result2.Value.String() != "public-value" {
			t.Errorf("expected original value, got '%s'", result2.Value.String())
		}
	})
}
