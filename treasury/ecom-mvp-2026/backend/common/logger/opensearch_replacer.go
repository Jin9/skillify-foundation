package logger

import "log/slog"

var opensearchKeys = map[string]string{
	"msg":   "message",
	"time":  "@timestamp",
	"level": "log.level",
}

// OpenSearchKeyReplacer replaces slog default keys with names that conform
// to the OpenSearch / Elasticsearch Common Schema (ECS) conventions.
//
//   - "msg"   → "message"     — standard ECS message field
//   - "time"  → "@timestamp"  — required by OpenSearch Dashboards for time-based queries
//   - "level" → "log.level"   — ECS-compliant log level field
func OpenSearchKeyReplacer(groups []string, a slog.Attr) (slog.Attr, bool) {
	for k, v := range opensearchKeys {
		if a.Key == k {
			return slog.String(v, a.Value.String()), true
		}
	}

	return a, false
}
