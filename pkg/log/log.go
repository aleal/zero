package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	zcontext "github.com/aleal/zero/internal/context"
	"github.com/aleal/zero/internal/uuid"
)

// Log levels
const (
	DebugLevel   = "DEBUG"
	InfoLevel    = "INFO"
	WarningLevel = "WARNING"
	ErrorLevel   = "ERROR"
	FatalLevel   = "FATAL"
	PanicLevel   = "PANIC"
	debugValue   = 0
	infoValue    = 1
	warningValue = 2
	errorValue   = 3
	fatalValue   = 4
	panicValue   = 5
)

var (
	levels = map[string]int{
		"DEBUG":   debugValue,
		"INFO":    infoValue,
		"WARNING": warningValue,
		"ERROR":   errorValue,
		"FATAL":   fatalValue,
		"PANIC":   panicValue,
	}
)

type logger struct {
	level   int
	logFunc LogFunc
}

type Logger interface {
	Debug(rctx context.Context, format string, a ...any)
	Info(rctx context.Context, format string, a ...any)
	Warning(rctx context.Context, format string, a ...any)
	Error(rctx context.Context, format string, a ...any)
	Fatal(rctx context.Context, format string, a ...any)
	Panic(rctx context.Context, format string, a ...any)
}

type LogFunc func(rctx context.Context, level, format string, a ...any)

// NewLogger creates a new logger with the default log function
func NewLogger() Logger {
	return NewLoggerWithFunc(Log)
}

// NewLoggerWithFunc creates a new logger with a custom log function
func NewLoggerWithFunc(logFunc LogFunc) Logger {
	return &logger{
		level:   getLevel(),
		logFunc: logFunc,
	}
}

// FromContext gets the logger from the context
func FromContext(rctx context.Context) Logger {
	return rctx.Value(zcontext.Logger).(Logger)
}

// SetLoggerToContext sets the logger to the context
func SetLoggerToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, zcontext.Logger, logger)
}

// Log logs a message with the time and request ID
func Log(rctx context.Context, level, format string, a ...any) {
	format = fmt.Sprintf("%s %s%s %s", level, rctx.Value(zcontext.ZID), getRequestID(rctx), format)
	log.Printf(format, a...)
}

// getRequestID extracts the request ID from the context, returning empty UUID if not found
func getRequestID(rctx context.Context) uuid.UUID {
	requestID := rctx.Value(zcontext.RequestID)
	if requestID == nil {
		return uuid.UUID("")
	}
	return requestID.(uuid.UUID)
}

func (l *logger) Debug(rctx context.Context, format string, a ...any) {
	if l.isDebug() {
		l.logFunc(rctx, DebugLevel, format, a...)
	}
}

func (l *logger) Info(rctx context.Context, format string, a ...any) {
	if l.isInfo() {
		l.logFunc(rctx, InfoLevel, format, a...)
	}
}

func (l *logger) Warning(rctx context.Context, format string, a ...any) {
	if l.isWarning() {
		l.logFunc(rctx, WarningLevel, format, a...)
	}
}

func (l *logger) Error(rctx context.Context, format string, a ...any) {
	if l.isError() {
		l.logFunc(rctx, ErrorLevel, format, a...)
	}
}

func (l *logger) Fatal(rctx context.Context, format string, a ...any) {
	if l.isFatal() {
		l.logFunc(rctx, FatalLevel, format, a...)
		os.Exit(1)
	}
}

func (l *logger) Panic(rctx context.Context, format string, a ...any) {
	if l.isPanic() {
		l.logFunc(rctx, PanicLevel, format, a...)
		panic(fmt.Sprintf(format, a...))
	}
}

func (l *logger) isDebug() bool {
	return l.level < infoValue
}

func (l *logger) isInfo() bool {
	return l.level < warningValue
}

func (l *logger) isWarning() bool {
	return l.level < errorValue
}

func (l *logger) isError() bool {
	return l.level < fatalValue
}

func (l *logger) isFatal() bool {
	return l.level < panicValue
}

func (l *logger) isPanic() bool {
	return l.level <= panicValue
}

func getLevel() int {
	if level := os.Getenv("ZERO_LOG_LEVEL"); level != "" {
		if l, ok := levels[strings.ToUpper(level)]; ok {
			return l
		}
	}
	return infoValue
}
