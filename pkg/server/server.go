// Package server provides a simple, lean, and blazingly fast HTTP server library.
// It includes built-in middleware support, graceful shutdown, and a clean API
// for building HTTP services.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aleal/zero/internal/handlers"
	"github.com/aleal/zero/pkg/config"
	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/middlewares"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

// zero represents a Zero HTTP server instance
type zero struct {
	*http.Server
	cfg *config.Config
	mux *http.ServeMux
	// handlers stores handlers for each method and pattern combination
	handlers    map[string]map[string]request.Handler
	middlewares []middlewares.Middleware
	ctx         context.Context
}

// Zero is the interface for the Zero server
type Zero interface {
	Get(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
	Post(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
	Put(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
	Delete(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
	Patch(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
	Handle(pattern string, method string, handler request.Handler, middlewares ...middlewares.Middleware)
	Handler() http.Handler
	Start() error
	Shutdown(ctx context.Context) error
}

// New creates a new Zero server instance with options
func New(ctx context.Context, options ...Option) Zero {
	cfg := config.Load()

	mux := http.NewServeMux()

	instance := &zero{
		mux:         mux,
		cfg:         cfg,
		handlers:    make(map[string]map[string]request.Handler),
		middlewares: make([]middlewares.Middleware, 0),
		ctx:         ctx,
	}

	for _, option := range options {
		option(instance)
	}

	// Build http.Server AFTER options so cfg mutations are reflected
	instance.Server = &http.Server{
		Addr:         instance.cfg.GetAddr(),
		Handler:      mux,
		ReadTimeout:  instance.cfg.ReadTimeout,
		WriteTimeout: instance.cfg.WriteTimeout,
		IdleTimeout:  instance.cfg.IdleTimeout,
	}

	instance.Get("/health", handlers.HealthCheckHandler())

	return instance
}

// Handler returns the server's HTTP handler for testing
func (z *zero) Handler() http.Handler {
	return z.mux
}

// Get registers a GET handler for the given pattern
func (z *zero) Get(pattern string, handler request.Handler, middlewares ...middlewares.Middleware) {
	z.Handle(pattern, http.MethodGet, handler, middlewares...)
}

// Post registers a POST handler for the given pattern
func (z *zero) Post(pattern string, handler request.Handler, middlewares ...middlewares.Middleware) {
	z.Handle(pattern, http.MethodPost, handler, middlewares...)
}

// Put registers a PUT handler for the given pattern
func (z *zero) Put(pattern string, handler request.Handler, middlewares ...middlewares.Middleware) {
	z.Handle(pattern, http.MethodPut, handler, middlewares...)
}

// Delete registers a DELETE handler for the given pattern
func (z *zero) Delete(pattern string, handler request.Handler, middlewares ...middlewares.Middleware) {
	z.Handle(pattern, http.MethodDelete, handler, middlewares...)
}

// Patch registers a PATCH handler for the given pattern
func (z *zero) Patch(pattern string, handler request.Handler, middlewares ...middlewares.Middleware) {
	z.Handle(pattern, http.MethodPatch, handler, middlewares...)
}

// Handle registers a handler for the given pattern and method
func (z *zero) Handle(pattern string, method string, handler request.Handler, middlewares ...middlewares.Middleware) {
	z.registerMethodHandler(pattern, method, handler, middlewares...)
}

// Start starts the server with graceful shutdown
func (z *zero) Start() error {
	printBanner()
	logger := log.FromContext(z.ctx)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("Starting Zero server", slog.String("address", z.Server.Addr))
		if err := z.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-errCh:
		return err
	case <-quit:
	}

	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := z.Server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.Any("error", err))
		return err
	}

	logger.Info("Server exited gracefully")
	return nil
}

// normalizePattern normalizes the pattern for the given pattern
func (z *zero) normalizePattern(pattern string) string {
	// Ensure leading slash, preserve trailing slash (Go 1.22+ subtree patterns)
	return "/" + strings.TrimLeft(pattern, "/")
}

// registerMethodHandler registers a handler for a specific HTTP method and pattern and applies the middlewares to the handler
func (z *zero) registerMethodHandler(pattern, method string, handler request.Handler, rmiddlewares ...middlewares.Middleware) {
	pattern = z.normalizePattern(pattern)

	logger := log.FromContext(z.ctx)
	logger.Info("Registering handler", slog.String("method", method), slog.String("pattern", pattern))

	registered := z.handlers[pattern] != nil
	if !registered {
		z.handlers[pattern] = make(map[string]request.Handler)
	}

	// Defensive copy to prevent slice aliasing
	allMiddlewares := make([]middlewares.Middleware, 0, len(z.middlewares)+len(rmiddlewares))
	allMiddlewares = append(allMiddlewares, z.middlewares...)
	allMiddlewares = append(allMiddlewares, rmiddlewares...)
	z.handlers[pattern][method] = middlewares.Chain(handler, allMiddlewares...)

	if !registered {
		z.mux.HandleFunc(pattern, z.methodRouter(pattern))
	}
}

// methodRouter creates a router that handles different HTTP methods for a pattern
func (z *zero) methodRouter(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers, exists := z.handlers[pattern]
		if !exists {
			http.NotFound(w, r)
			return
		}

		handler, methodExists := handlers[r.Method]
		if !methodExists {
			methodNotAllowed(handlers, w, r)
			return
		}

		handler(w, r)
	}
}

// methodNotAllowed returns a 405 response when the method is not allowed
func methodNotAllowed(handlers map[string]request.Handler, w http.ResponseWriter, r *http.Request) {
	allowedMethods := make([]string, 0, len(handlers))
	for method := range handlers {
		allowedMethods = append(allowedMethods, method)
	}
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	response.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method))
}
