package constant

import (
	"log/slog"
)

var (
	Severitydefault = slog.Level(slog.LevelDebug)
	SeverityInfo    = slog.Level(slog.LevelInfo)
	SeverityWarn    = slog.Level(slog.LevelWarn)
	SeverityError   = slog.Level(slog.LevelError)
)

const SERVICE_NAME = "web-ya-hime"
