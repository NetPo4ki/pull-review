package log

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(appEnv, level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:       lvl,
		AddSource:   appEnv != "prod" && appEnv != "production",
		ReplaceAttr: nil,
	}

	var handler slog.Handler
	if appEnv == "prod" || appEnv == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
