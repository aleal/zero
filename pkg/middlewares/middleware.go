// Package middlewares provides HTTP middleware functionality for the Zero server.
// It includes common middleware like CORS, logging, and recovery, as well as
// utilities for chaining multiple middleware functions.
package middlewares

import (
	"github.com/aleal/zero/pkg/request"
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
