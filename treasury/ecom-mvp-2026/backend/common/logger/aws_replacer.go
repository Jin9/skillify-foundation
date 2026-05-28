package logger

import "log/slog"

var awsKeys = map[string]string{
	"msg":  "message",
	"time": "timestamp",
}

// AWSKeyReplacer replaces slog default keys with names that are easier to query in CloudWatch.
func AWSKeyReplacer(groups []string, a slog.Attr) (slog.Attr, bool) {
	for k, v := range awsKeys {
		if a.Key == k {
			return slog.String(v, a.Value.String()), true
		}
	}

	return a, false
}
