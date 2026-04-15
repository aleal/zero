package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/requestid"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.written {
		return
	}
	rw.statusCode = code
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK) // Default to 200 if Write is called first
	}
	return rw.ResponseWriter.Write(b)
}

// Unwrap returns the underlying ResponseWriter so http.NewResponseController
// can reach Flusher, Hijacker, etc.
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// Logging middleware logs request details
func Logging(logger *slog.Logger) Middleware {
	return func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rid := requestid.New()
			logger := logger.With(
				slog.String("requestId", rid),
				slog.String("method", r.Method),
				slog.String("requestURI", r.RequestURI),
				slog.String("remoteAddr", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
			)
			logger.Info("Request started")
			rctx := requestid.WithContext(r.Context(), rid)
			rctx = log.SetLoggerToContext(rctx, logger)
			r = r.WithContext(rctx)
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK, written: false}
			next(rw, r)
			duration := time.Since(start).Microseconds()
			logger.Info("Request completed", slog.Int64("durationMicros", duration), slog.Int("statusCode", rw.statusCode))
		}
	}
}
