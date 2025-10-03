// Package handlers provides built-in HTTP handlers for the Zero server.
// It includes common handlers like health checks and other utility endpoints.
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

// HealthCheckHandler returns a health check response
func HealthCheckHandler() request.Handler {
	var start = time.Now()
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]any{
			"service": "zero",
			"version": "0.0.1",
			"uptime":  fmt.Sprintf("%v", time.Since(start)),
		})
	}
}
