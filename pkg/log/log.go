package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type contextKey struct{}

var loggerContextKey = contextKey{}

// NewLogger creates a new logger with the default settings
// The logger is configured to log to stdout in JSON format
// The log level is set to the value of the ZERO_LOG_LEVEL environment variable
// If the ZERO_LOG_LEVEL environment variable is not set, the log level is set to INFO
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLevel(),
	}))
}

// FromContext gets the logger from the context
// If the logger is not found in the context, a new logger is created and set to the context
func FromContext(rctx context.Context) *slog.Logger {
	if l, ok := rctx.Value(loggerContextKey).(*slog.Logger); ok {
		return l
	}
	return New()
}

// SetToContext sets the logger to the context
func SetToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// WithArgs adds the given arguments to the logger in the context
func WithArgs(ctx context.Context, args ...any) context.Context {
	if len(args) > 0 {
		return SetToContext(ctx, FromContext(ctx).With(args...))
	}
	return ctx
}

// getLevel gets the log level from the environment variable
func getLevel() slog.Level {
	if level := os.Getenv("ZERO_LOG_LEVEL"); level != "" {
		switch strings.ToUpper(level) {
		case "DEBUG":
			return slog.LevelDebug
		case "INFO":
			return slog.LevelInfo
		case "WARNING":
			return slog.LevelWarn
		case "ERROR":
			return slog.LevelError
		}
	}
	return slog.LevelInfo
}
