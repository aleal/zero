// Package context provides context key definitions for the Zero server.
// It defines the keys used to store request-specific data in the context.
package context

// ContextKey represents a key for storing values in the request context
type ContextKey string

const (
	// RequestID is the context key for storing the request ID
	RequestID ContextKey = "requestID"
	// ZID is the context key for storing the Zero instance ID
	ZID ContextKey = "zid"
	// Logger is the context key for storing the logger instance
	Logger ContextKey = "logger"
)
