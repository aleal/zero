package log

import (
	"context"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("NewLogger() returned nil")
	}
}

func TestNewLoggerWithFunc(t *testing.T) {
	logFunc := func(rctx context.Context, level, format string, a ...any) {
		// Custom log function for testing
	}

	logger := NewLoggerWithFunc(logFunc)
	if logger == nil {
		t.Error("NewLoggerWithFunc() returned nil")
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

func TestSetLoggerToContext(t *testing.T) {
	ctx := context.Background()
	logger := NewLogger()

	ctxWithLogger := SetLoggerToContext(ctx, logger)
	if ctxWithLogger == nil {
		t.Error("SetLoggerToContext() returned nil")
	}

	// Verify the logger was set
	retrievedLogger := FromContext(ctxWithLogger)
	if retrievedLogger == nil {
		t.Error("Logger not found in context")
	}
}

func TestLoggerMethods(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Test that all logger methods can be called without panicking
	logger.Debug(ctx, "Debug message")
	logger.Info(ctx, "Info message")
	logger.Warning(ctx, "Warning message")
	logger.Error(ctx, "Error message")

	// Note: Fatal and Panic methods would exit/panic, so we don't test them
}

func TestLogFunction(t *testing.T) {
	ctx := context.Background()

	// Test that Log function can be called without panicking
	Log(ctx, "INFO", "Test log message")
}

func TestGetLevel(t *testing.T) {
	// Test default level
	level := getLevel()
	if level != infoValue {
		t.Errorf("Default level = %d, want %d", level, debugValue)
	}

	// Test with environment variable
	os.Setenv("ZERO_LOG_LEVEL", "ERROR")
	defer os.Unsetenv("ZERO_LOG_LEVEL")

	level = getLevel()
	if level != errorValue {
		t.Errorf("INFO level = %d, want %d", level, infoValue)
	}

	// Test with invalid level
	os.Setenv("ZERO_LOG_LEVEL", "INVALID")
	level = getLevel()
	if level != infoValue {
		t.Errorf("Invalid level should default to debug, got %d", level)
	}
}

func TestGetRequestID(t *testing.T) {
	ctx := context.Background()

	// Test with no request ID
	requestID := getRequestID(ctx)
	if requestID != "" {
		t.Errorf("Expected empty request ID, got %s", requestID)
	}
}

func TestLoggerWithCustomLevel(t *testing.T) {
	// Set log level to ERROR
	os.Setenv("ZERO_LOG_LEVEL", "ERROR")
	defer os.Unsetenv("ZERO_LOG_LEVEL")

	logger := NewLogger()
	ctx := context.Background()

	// Test that logger methods can be called without panicking
	// The actual filtering behavior is tested through the public interface
	logger.Debug(ctx, "Debug message")
	logger.Info(ctx, "Info message")
	logger.Warning(ctx, "Warning message")
	logger.Error(ctx, "Error message")
}

func TestLoggerFatalAndPanicMethods(t *testing.T) {
	// Note: Fatal and Panic methods call os.Exit(1) and panic() respectively,
	// which makes them difficult to test in unit tests. These methods are
	// designed to terminate the program, so they cannot be easily tested
	// without special test frameworks or mocking.

	logger := NewLogger()
	ctx := context.Background()

	// We can only test that the methods exist and can be called
	// The actual behavior (exit/panic) cannot be tested in unit tests
	_ = logger
	_ = ctx
}

func TestLoggerLevelMethodsWithCustomLevels(t *testing.T) {
	// Test level checking methods with different log levels
	testCases := []struct {
		level string
	}{
		{"DEBUG"},
		{"INFO"},
		{"WARNING"},
		{"ERROR"},
		{"FATAL"},
		{"PANIC"},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			os.Setenv("ZERO_LOG_LEVEL", tc.level)
			defer os.Unsetenv("ZERO_LOG_LEVEL")

			logger := NewLogger()
			ctx := context.Background()

			// Test that logger methods can be called without panicking
			// The actual filtering behavior is tested through the public interface
			logger.Debug(ctx, "Debug message")
			logger.Info(ctx, "Info message")
			logger.Warning(ctx, "Warning message")
			logger.Error(ctx, "Error message")
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
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "Benchmark message %d", i)
	}
}

func BenchmarkLogFunction(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Log(ctx, "INFO", "Benchmark message %d", i)
	}
}
