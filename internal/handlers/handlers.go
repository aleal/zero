// Package handlers provides built-in HTTP handlers for the Zero server.
// It includes common handlers like health checks and other utility endpoints.
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aleal/zero/pkg/metadata"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

// HealthCheckHandler returns a health check response
func HealthCheckHandler() request.Handler {
	var start = time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(start).Truncate(time.Second)

		response.WriteJSON(w, http.StatusOK, map[string]any{
			"service": "zero",
			"version": metadata.GetVersion(),
			"uptime":  fmt.Sprintf("%v", uptime),
		})
	}
}
