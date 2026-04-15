package server

import (
	"log/slog"
	"time"

	"github.com/aleal/zero/pkg/config"
	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/middlewares"
)

type Option func(*zero)

// WithConfig sets the config for the server
func WithConfig(cfg *config.Config) Option {
	return func(z *zero) {
		z.cfg = cfg
	}
}

// WithHost sets the host for the server
func WithHost(host string) Option {
	return func(z *zero) {
		z.cfg.Host = host
	}
}

// WithPort sets the port for the server
func WithPort(port int) Option {
	return func(z *zero) {
		z.cfg.Port = port
	}
}

// WithReadTimeout sets the read timeout for the server
func WithReadTimeout(timeout time.Duration) Option {
	return func(z *zero) {
		z.cfg.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets the write timeout for the server
func WithWriteTimeout(timeout time.Duration) Option {
	return func(z *zero) {
		z.cfg.WriteTimeout = timeout
	}
}

// WithIdleTimeout sets the idle timeout for the server
func WithIdleTimeout(timeout time.Duration) Option {
	return func(z *zero) {
		z.cfg.IdleTimeout = timeout
	}
}

// WithMaxJSONBodySize sets the max JSON body size for the server
func WithMaxJSONBodySize(size int64) Option {
	return func(z *zero) {
		z.cfg.MaxJSONBodySize = size
	}
}

// WithMaxUploadedFileSize sets the max uploaded file size for the server
func WithMaxUploadedFileSize(size int64) Option {
	return func(z *zero) {
		z.cfg.MaxUploadedFileSize = size
	}
}

// WithLogging applies the logging middleware to the server
// It applies the logging middleware to the server with the given logger
func WithLogging(logger *slog.Logger, priority middlewares.MiddlewarePriority) Option {
	return func(z *zero) {
		z.ctx = log.SetLoggerToContext(z.ctx, logger)
		z.middlewares = middlewares.Append(z.middlewares, middlewares.Logging(logger), priority)
	}
}

// WithCORS applies the CORS middleware to the server
// It applies the CORS middleware to the server with the given allowed origins
func WithCORS(origins []string, priority middlewares.MiddlewarePriority) Option {
	return func(z *zero) {
		z.middlewares = middlewares.Append(z.middlewares, middlewares.CORS(origins), priority)
	}
}

// WithRecovery applies the recovery middleware to the server
// It applies the recovery middleware to the server with the given priority
func WithRecovery(priority middlewares.MiddlewarePriority) Option {
	return func(z *zero) {
		z.middlewares = middlewares.Append(z.middlewares, middlewares.Recovery(), priority)
	}
}

// WithMiddleware applies a middleware to the server
// It applies the middleware to the server with the given middleware and priority
func WithMiddleware(middleware middlewares.Middleware, priority middlewares.MiddlewarePriority) Option {
	return func(z *zero) {
		z.middlewares = middlewares.Append(z.middlewares, middleware, priority)
	}
}

// WithDefaultLogging applies the default logging middlewares to the server
// It applies the logging middleware to the server with the default logger
// It is equivalent to calling WithLogging(log.NewLogger(), middlewares.MiddlewarePriorityLow)
func WithDefaultLogging() Option {
	return WithLogging(log.NewLogger(), middlewares.MiddlewarePriorityLow)
}

// WithDefaultCORS applies the default CORS middleware to the server
// It applies the CORS middleware to the server with the default allowed origins
// It is equivalent to calling WithCORS([]string{"*"}, middlewares.MiddlewarePriorityLow)
func WithDefaultCORS() Option {
	return WithCORS([]string{"*"}, middlewares.MiddlewarePriorityLow)
}

// WithDefaultRecovery applies the default recovery middleware to the server
// It applies the recovery middleware to the server
// It is equivalent to calling WithRecovery(middlewares.MiddlewarePriorityLow)
func WithDefaultRecovery() Option {
	return WithRecovery(middlewares.MiddlewarePriorityLow)
}

// WithDefaultMiddlewares applies the default middlewares to the server
// This is a convenience function that applies the default middlewares to the server
// It is equivalent to calling WithDefaultLogging(), WithDefaultCORS(), and WithDefaultRecovery() in sequence
func WithDefaultMiddlewares() Option {
	return func(z *zero) {
		WithDefaultLogging()(z)
		WithDefaultCORS()(z)
		WithDefaultRecovery()(z)
	}
}
