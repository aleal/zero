package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/request"
)

type recoveryWriter struct {
	http.ResponseWriter
	written bool
}

func (rw *recoveryWriter) WriteHeader(code int) {
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *recoveryWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(b)
}

func (rw *recoveryWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// Recovery middleware recovers from panics
func Recovery() Middleware {
	return func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			rw := &recoveryWriter{ResponseWriter: w}
			defer func() {
				if err := recover(); err != nil {
					logger := log.FromContext(r.Context())
					logger.Error("Panic recovered", slog.Any("error", err))
					if !rw.written {
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				}
			}()
			next(rw, r)
		}
	}
}
