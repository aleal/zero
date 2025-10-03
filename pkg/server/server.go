// Package server provides a simple, lean, and blazingly fast HTTP server library.
// It includes built-in middleware support, graceful shutdown, and a clean API
// for building HTTP services.
package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	zcontext "github.com/aleal/zero/internal/context"
	"github.com/aleal/zero/internal/handlers"
	"github.com/aleal/zero/internal/uuid"
	"github.com/aleal/zero/pkg/config"
	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/middlewares"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

// Server represents a Zero HTTP server instance
type server struct {
	*http.Server
	cfg *config.Config
	mux *http.ServeMux
	// handlers stores handlers for each method and pattern combination
	handlers    map[string]map[string]request.Handler
	middlewares []middlewares.Middleware
	ctx         context.Context
}

// Server is the interface for the Zero server
type Server interface {
	Get(pattern string, handler request.Handler)
	Post(pattern string, handler request.Handler)
	Put(pattern string, handler request.Handler)
	Delete(pattern string, handler request.Handler)
	Patch(pattern string, handler request.Handler)
	Handle(pattern string, method string, handler request.Handler)
	Middlewares(middlewares ...middlewares.Middleware)
	Start()
}

// NewServer creates a new Zero server instance with command line flags and environment variables
func NewServer(ctx context.Context) Server {
	// Load configuration from environment
	cfg := config.Default()

	// Parse command line flags
	var (
		host = flag.String("host", cfg.Host, "Server host")
		port = flag.Int("port", cfg.Port, "Server port")
	)
	flag.Parse()

	cfg.SetHost(*host)
	cfg.SetPort(*port)

	return NewServerWithConfig(ctx, cfg)
}

// NewServerWithConfig creates a new Zero server instance with configuration
func NewServerWithConfig(ctx context.Context, cfg *config.Config) Server {
	return NewServerWithConfigAndLogger(ctx, cfg, log.NewLogger())
}

// NewServerWithConfigAndLogger creates a new Zero server instance with configuration and logger
func NewServerWithConfigAndLogger(ctx context.Context, cfg *config.Config, logger log.Logger) Server {
	zctx := context.WithValue(ctx, zcontext.ZID, uuid.GenerateUUID())

	zctx = context.WithValue(zctx, zcontext.Logger, logger)

	mux := http.NewServeMux()

	// Use default config if none provided
	if cfg == nil {
		cfg = config.Default()
	}

	server := &server{
		Server: &http.Server{
			Addr:         cfg.GetAddr(),
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		mux:         mux,
		cfg:         cfg,
		handlers:    make(map[string]map[string]request.Handler),
		middlewares: make([]middlewares.Middleware, 0),
		ctx:         zctx,
	}

	server.registerDefaultMiddlewares()
	server.Get("/health", handlers.HealthCheckHandler())

	return server
}

// Middleware adds a middleware to the server
func (z *server) Middlewares(middlewares ...middlewares.Middleware) {
	z.middlewares = append(z.middlewares, middlewares...)
}

// Get registers a GET handler for the given pattern
func (z *server) Get(pattern string, handler request.Handler) {
	z.Handle(pattern, http.MethodGet, handler)
}

// Post registers a POST handler for the given pattern
func (z *server) Post(pattern string, handler request.Handler) {
	z.Handle(pattern, http.MethodPost, handler)
}

// Put registers a PUT handler for the given pattern
func (z *server) Put(pattern string, handler request.Handler) {
	z.Handle(pattern, http.MethodPut, handler)
}

// Delete registers a DELETE handler for the given pattern
func (z *server) Delete(pattern string, handler request.Handler) {
	z.Handle(pattern, http.MethodDelete, handler)
}

// Patch registers a PATCH handler for the given pattern
func (z *server) Patch(pattern string, handler request.Handler) {
	z.Handle(pattern, http.MethodPatch, handler)
}

// Handle registers a handler for the given pattern and method
func (z *server) Handle(pattern string, method string, handler request.Handler) {
	z.registerMethodHandler(pattern, method, handler)
}

// Start starts the server with graceful shutdown
func (z *server) Start() {
	logger := log.FromContext(z.ctx)
	quit := make(chan os.Signal, 1)
	// Start server in a goroutine
	go func() {
		printBanner()
		logger.Info(z.ctx, "Starting Zero server on %s", z.Server.Addr)
		if err := z.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(z.ctx, "Server error: %v", err)
			quit <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	logger.Info(z.ctx, "Shutting down server...")

	// Graceful shutdown
	ctxCancel, cancel := context.WithTimeout(z.ctx, 30*time.Second)
	defer cancel()

	if err := z.Shutdown(ctxCancel); err != nil && err != http.ErrServerClosed {
		logger.Fatal(z.ctx, "Server forced to shutdown: %v", err)
	}

	logger.Info(z.ctx, "Server exited 0")
}

// registerDefaultMiddlewares registers the default middlewares for the server
func (z *server) registerDefaultMiddlewares() {
	if z.cfg.EnableLogging {
		z.Middlewares(middlewares.Logging())
	}

	if z.cfg.EnableCORS {
		z.Middlewares(middlewares.CORS(z.cfg.AllowedOrigins))
	}

	if z.cfg.EnableRecovery {
		z.Middlewares(middlewares.Recovery())
	}
}

// normalizePattern normalizes the pattern for the given pattern
func (z *server) normalizePattern(pattern string) string {
	return fmt.Sprintf("/%s", strings.Trim(pattern, "/"))
}

// registerMethodHandler registers a handler for a specific HTTP method and pattern
func (z *server) registerMethodHandler(pattern, method string, handler request.Handler) {
	// Initialize the pattern map if it doesn't exist
	pattern = z.normalizePattern(pattern)

	logger := log.FromContext(z.ctx)
	logger.Info(z.ctx, "Registering handler for %s %s ...", method, pattern)

	registered := z.handlers[pattern] != nil
	if !registered {
		z.handlers[pattern] = make(map[string]request.Handler)
	}

	// Store the handler for this method and pattern and apply middlewares
	z.handlers[pattern][method] = middlewares.Chain(handler, z.middlewares...)

	// Register the pattern with the mux only once
	if !registered {
		z.mux.HandleFunc(pattern, z.methodRouter(pattern))
	}
}

// methodRouter creates a router that handles different HTTP methods for a pattern
func (z *server) methodRouter(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the handlers for this pattern
		handlers, exists := z.handlers[pattern]
		if !exists {
			http.NotFound(w, r)
			return
		}

		// Get the handler for the specific method
		handler, methodExists := handlers[r.Method]
		if !methodExists {
			methodNotAllowed(handlers, w, r)
			return
		}
		rctx := context.WithValue(z.ctx, zcontext.RequestID, uuid.GenerateUUID())
		// Call the appropriate handler
		handler(rctx, w, r)
	}
}

// methodNotAllowed returns a 405 response when the method is not allowed
func methodNotAllowed(handlers map[string]request.Handler, w http.ResponseWriter, r *http.Request) {
	// Method not allowed - return 405 with allowed methods
	allowedMethods := make([]string, 0, len(handlers))
	for method := range handlers {
		allowedMethods = append(allowedMethods, method)
	}
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	response.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method))
}
