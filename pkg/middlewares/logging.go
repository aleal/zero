package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/requestid"
	"github.com/aleal/zero/pkg/response"
)

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
			next(w, r)
			duration := time.Since(start).Microseconds()
			statusCode := response.StatusCode(w)
			logger.Info("Request completed", slog.Int64("durationMicros", duration), slog.Int("statusCode", statusCode))
		}
	}
}
