// Package requestid provides lightweight request ID generation for log correlation.
// IDs use the format hostname-counter for uniqueness per host without crypto overhead.
package requestid

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
)

type contextKey struct{}

var (
	counter  atomic.Int64
	hostname string
)

func init() {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown"
	}
	hostname = h
}

// New generates a request ID in the format "hostname-counter".
func New() string {
	return fmt.Sprintf("%s-%d", hostname, counter.Add(1))
}

// WithContext stores a request ID in the context.
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// FromContext retrieves the request ID from context, or empty string if absent.
func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(contextKey{}).(string); ok {
		return id
	}
	return ""
}
