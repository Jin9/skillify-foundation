package logger

import (
	"log/slog"
	"maps"
	"sync"
)

var (
	censorsMu sync.RWMutex
	censors   = map[string]string{}
)

// DefaultCensors defines the fields that are masked by default in all log output.
// Add entries here when a new sensitive field is introduced across the codebase.
var DefaultCensors = map[string]string{
	"password":     "***REDACTED***",
	"access_token": "***TOKEN***",
	"api_key":      "***API_KEY***",
	"secret":       "***SECRET***",
}

func init() {
	LoadCensorsFromMap(DefaultCensors)
}

// AddCensor adds a key-value pair to censor sensitive information in logs.
// key: the log attribute key to censor (e.g., "password", "api_key")
// maskedValue: the value to display instead (e.g., "***REDACTED***")
func AddCensor(key, maskedValue string) {
	censorsMu.Lock()
	defer censorsMu.Unlock()
	censors[key] = maskedValue
}

// LoadCensorsFromMap loads multiple censors from a map[string]string.
// This is useful when loading censor configuration from JSON or other config files.
func LoadCensorsFromMap(censorMap map[string]string) {
	censorsMu.Lock()
	defer censorsMu.Unlock()
	maps.Copy(censors, censorMap)
}

// CensorReplacer replaces or masks log values for sensitive information.
func CensorReplacer(groups []string, a slog.Attr) (slog.Attr, bool) {
	censorsMu.RLock()
	defer censorsMu.RUnlock()
	for k, v := range censors {
		if a.Key == k {
			return slog.String(k, v), true
		}
	}
	return a, false
}

// ResetCensors clears all registered censors. Intended for testing only.
func ResetCensors() {
	censorsMu.Lock()
	defer censorsMu.Unlock()
	clear(censors)
}
