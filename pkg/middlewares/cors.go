package middlewares

import (
	"net/http"

	"github.com/aleal/zero/pkg/request"
)

// CORS middleware adds CORS headers
func CORS(allowedOrigins []string) Middleware {
	wildcard := false
	for _, o := range allowedOrigins {
		if o == "*" {
			wildcard = true
			break
		}
	}

	return func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			matched := false

			if origin != "" {
				if wildcard {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					matched = true
				} else {
					for _, allowed := range allowedOrigins {
						if allowed == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							w.Header().Set("Access-Control-Allow-Credentials", "true")
							w.Header().Set("Vary", "Origin")
							matched = true
							break
						}
					}
				}
			}

			if matched {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "600")
			}

			if matched && r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}
