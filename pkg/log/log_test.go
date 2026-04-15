package log

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("NewLogger() returned nil")
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	logger := NewLogger()
	ctxWithLogger := SetLoggerToContext(ctx, logger)

	retrievedLogger := FromContext(ctxWithLogger)
	if retrievedLogger == nil {
		t.Error("FromContext() returned nil")
	}
}

func TestFromContextWithoutLogger(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx)
	if logger == nil {
		t.Error("FromContext() should return fallback logger, got nil")
	}
}

func TestSetLoggerToContext(t *testing.T) {
	ctx := context.Background()
	logger := NewLogger()

	ctxWithLogger := SetLoggerToContext(ctx, logger)
	if ctxWithLogger == nil {
		t.Error("SetLoggerToContext() returned nil")
	}

	retrievedLogger := FromContext(ctxWithLogger)
	if retrievedLogger == nil {
		t.Error("Logger not found in context")
	}
}

func TestLoggerMethods(t *testing.T) {
	logger := NewLogger()

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
}

func TestGetLevel(t *testing.T) {
	level := getLevel()
	if level != slog.LevelInfo {
		t.Errorf("Default level = %d, want %d", level, slog.LevelInfo)
	}

	t.Setenv("ZERO_LOG_LEVEL", "ERROR")
	level = getLevel()
	if level != slog.LevelError {
		t.Errorf("ERROR level = %d, want %d", level, slog.LevelError)
	}
}

func TestGetLevelInvalid(t *testing.T) {
	t.Setenv("ZERO_LOG_LEVEL", "INVALID")
	level := getLevel()
	if level != slog.LevelInfo {
		t.Errorf("Invalid level should default to INFO, got %d", level)
	}
}

func TestLoggerWithCustomLevel(t *testing.T) {
	t.Setenv("ZERO_LOG_LEVEL", "ERROR")

	logger := NewLogger()

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
}

func TestLoggerLevelMethodsWithCustomLevels(t *testing.T) {
	testCases := []struct {
		level string
	}{
		{"DEBUG"},
		{"INFO"},
		{"WARNING"},
		{"ERROR"},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			t.Setenv("ZERO_LOG_LEVEL", tc.level)

			logger := NewLogger()

			logger.Debug("Debug message")
			logger.Info("Info message")
			logger.Warn("Warning message")
			logger.Error("Error message")
		})
	}
}

func BenchmarkNewLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewLogger()
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	logger := NewLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", "i", i)
	}
}
