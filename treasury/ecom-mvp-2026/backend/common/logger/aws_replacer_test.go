package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func TestAWSKeyReplacer(t *testing.T) {
	oldLogLevel := LogLevel
	LogLevel = slog.LevelInfo
	defer func() { LogLevel = oldLogLevel }()

	t.Run("AWS key replacer", func(t *testing.T) {
		os.Unsetenv("ENV")
		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		l := New(AWSKeyReplacer)

		l.Info("message")

		var m map[string]any
		json.NewDecoder(buf).Decode(&m)

		if _, ok := m["level"]; !ok {
			t.Errorf("replacer should keep key %q for CloudWatch queries: actual %v\n", "level", m)
		}

		if _, ok := m["message"]; !ok {
			t.Errorf("replacer should replace key %q to %q: actual %v\n", "msg", "message", m)
		}

		if _, ok := m["timestamp"]; !ok {
			t.Errorf("replacer should replace key %q to %q: actual %v\n", "time", "timestamp", m)
		}

		if _, ok := m["msg"]; ok {
			t.Errorf("replacer should not keep key %q after replacement: actual %v\n", "msg", m)
		}

		if _, ok := m["time"]; ok {
			t.Errorf("replacer should not keep key %q after replacement: actual %v\n", "time", m)
		}
	})

	t.Run("not found any key", func(t *testing.T) {
		_, ok := AWSKeyReplacer([]string{}, slog.Attr{})
		if ok {
			t.Errorf("not any matched key expect false")
		}
	})
}
