// Package middlewares provides HTTP middleware functionality for the Zero server.
// It includes common middleware like CORS, logging, and recovery, as well as
// utilities for chaining multiple middleware functions.
package middlewares

import (
	"github.com/aleal/zero/pkg/request"
)

type MiddlewarePriority int

const (
	// MiddlewarePriorityLow is the lowest priority for a middleware to be executed at the end of the chain
	MiddlewarePriorityLow MiddlewarePriority = iota
	// MiddlewarePriorityHigh is the highest priority for a middleware to be executed at the beginning of the chain
	MiddlewarePriorityHigh
)

// Middleware represents a middleware function
type Middleware func(next request.Handler) (handler request.Handler)

// Chain applies multiple middleware functions to a handler
func Chain(handler request.Handler, middlewares ...Middleware) request.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// Append appends a middleware to the chain at the given priority
func Append(middlewares []Middleware, middleware Middleware, priority MiddlewarePriority) []Middleware {
	switch priority {
	case MiddlewarePriorityHigh:
		return append([]Middleware{middleware}, middlewares...)
	case MiddlewarePriorityLow:
		fallthrough
	default:
		return append(middlewares, middleware)
	}
}
