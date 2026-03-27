package plate

import (
	"strconv"
	"strings"
	"time"
)

const defaultSyncInterval = 6 * time.Hour

func parseSyncDuration(spec string) time.Duration {
	s := strings.TrimSpace(strings.ToLower(spec))
	if s == "" {
		return defaultSyncInterval
	}

	if d, err := time.ParseDuration(s); err == nil && d > 0 {
		return d
	}

	parts := strings.Fields(s)
	if len(parts) == 2 {
		if n, err := strconv.Atoi(parts[0]); err == nil && n > 0 {
			switch strings.TrimSpace(parts[1]) {
			case "h", "hour", "hours":
				return time.Duration(n) * time.Hour
			case "m", "minute", "minutes":
				return time.Duration(n) * time.Minute
			case "s", "second", "seconds":
				return time.Duration(n) * time.Second
			}
		}
	}

	return defaultSyncInterval
}
