// Package main demonstrates how to use the Zero HTTP server library.
// This example shows how to create a server, add routes, and handle different
// types of HTTP requests.
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/parser"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
	zero "github.com/aleal/zero/pkg/server"
)

// main initializes and starts the Zero HTTP server with example routes
func main() {
	ctx := context.Background()
	// Create server instance with flags and environment variables
	server := zero.NewServer(ctx)

	// Add routes
	setupRoutes(server)

	// Start server
	server.Start()
}

// setupRoutes configures all the HTTP routes for the example server
func setupRoutes(server zero.Server) {
	// Health check (already added by NewServer)
	// server.Get("/health", handlers.HealthCheckHandler)

	// API routes
	server.Get("/hello", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		logger := log.FromContext(rctx)
		logger.Info(rctx, "Hello from Zero!")
		response.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Hello from Zero!",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
		logger.Info(rctx, "Done!")
	})

	server.Get("/users", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		users := []map[string]interface{}{
			{"id": 1, "name": "Alice", "email": "alice@example.com"},
			{"id": 2, "name": "Bob", "email": "bob@example.com"},
			{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
		}
		response.WriteJSON(w, http.StatusOK, users)
	})

	server.Get("/users/{id}", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		// In a real application, you would parse the ID from the URL
		response.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"id":    request.GetPathParam(r, "id"),
			"name":  "Alice",
			"email": "alice@example.com",
		})
	})

	server.Post("/users", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := parser.ParseJSONBody(r, &user); err != nil {
			response.WriteErrorMsg(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// In a real application, you would save the user to a database
		response.WriteJSON(w, http.StatusCreated, map[string]interface{}{
			"id":    4,
			"name":  user.Name,
			"email": user.Email,
		})
	})

	server.Get("/status", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]any{
			"status":    "running",
			"uptime":    "1h 23m 45s",
			"requests":  1234,
			"memory":    "15.2 MB",
			"cpu_usage": "2.1%",
		})
	})

	// Static file serving example
	server.Get("/static", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "examples/server/static/index.html")
	})

	// 404 handler for unmatched routes
	server.Get("/", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		// Serve a simple HTML page
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Zero HTTP Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        code { background: #e0e0e0; padding: 2px 4px; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Zero HTTP Server</h1>
        <p>A simple, lean, and blazingly fast HTTP server library built with pure Go.</p>
        
        <h2>Available Endpoints:</h2>
        <div class="endpoint">
            <strong>GET</strong> <code>/health</code> - Health check endpoint
        </div>
        <div class="endpoint">
            <strong>GET</strong> <code>/hello</code> - Hello world endpoint
        </div>
        <div class="endpoint">
            <strong>GET</strong> <code>/users</code> - List users
        </div>
        <div class="endpoint">
            <strong>GET</strong> <code>/users/:id</code> - Get user by ID
        </div>
        <div class="endpoint">
            <strong>POST</strong> <code>/users</code> - Create new user
        </div>
        <div class="endpoint">
            <strong>GET</strong> <code>/status</code> - Server status
        </div>
        
        <h2>Try it out:</h2>
        <p>Use curl or your browser to test the endpoints:</p>
        <pre><code>curl http://localhost:8080/api/hello</code></pre>
    </div>
</body>
</html>
        `))
	})
}
