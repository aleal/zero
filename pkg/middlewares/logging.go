package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/request"
)

// Logging middleware logs request details
func Logging() Middleware {
	return func(next request.Handler) request.Handler {
		return func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger := log.FromContext(rctx)
			if logger != nil {
				logger.Info(rctx, "%s %s %s started", r.Method, r.RequestURI, r.RemoteAddr)
			}
			next(rctx, w, r)
			if logger != nil {
				duration := time.Since(start)
				statusCode := w.Header().Get("Status")
				logger.Info(rctx, "%s %s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, statusCode, duration)
			}
		}
	}
}
