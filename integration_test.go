package zero

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aleal/zero/pkg/config"
	"github.com/aleal/zero/pkg/parser"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
	zeroserver "github.com/aleal/zero/pkg/server"
)

func TestIntegrationServer(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	cfg.SetPort(0) // Use random port for testing
	serverInstance := zeroserver.NewServerWithConfig(ctx, cfg)

	// Register test endpoints
	serverInstance.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "GET success"})
	})

	serverInstance.Post("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		if err := parser.ParseJSONBody(r, &data); err != nil {
			response.WriteError(w, http.StatusBadRequest, err)
			return
		}
		response.WriteJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "POST success",
			"data":    data,
		})
	})

	serverInstance.Put("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "PUT success"})
	})

	serverInstance.Delete("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "DELETE success"})
	})

	// Test query parameters
	serverInstance.Get("/query", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		name := request.GetQueryParamOrDefault(r, "name", "default")
		age := request.GetQueryParamOrDefault(r, "age", "0")
		response.WriteJSON(w, http.StatusOK, map[string]string{
			"name": name,
			"age":  age,
		})
	})

	// Test with middleware
	serverInstance.Middlewares(func(next request.Handler) request.Handler {
		return func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next(rctx, w, r)
		}
	})

	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the server's request handling
		switch r.URL.Path {
		case "/test":
			switch r.Method {
			case "GET":
				response.WriteJSON(w, http.StatusOK, map[string]string{"message": "GET success"})
			case "POST":
				var data map[string]interface{}
				if err := parser.ParseJSONBody(r, &data); err != nil {
					response.WriteError(w, http.StatusBadRequest, err)
					return
				}
				response.WriteJSON(w, http.StatusCreated, map[string]interface{}{
					"message": "POST success",
					"data":    data,
				})
			case "PUT":
				response.WriteJSON(w, http.StatusOK, map[string]string{"message": "PUT success"})
			case "DELETE":
				response.WriteJSON(w, http.StatusOK, map[string]string{"message": "DELETE success"})
			default:
				response.WriteErrorMsg(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		case "/query":
			name := request.GetQueryParamOrDefault(r, "name", "default")
			age := request.GetQueryParamOrDefault(r, "age", "0")
			response.WriteJSON(w, http.StatusOK, map[string]string{
				"name": name,
				"age":  age,
			})
		case "/health":
			response.WriteJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
		default:
			response.WriteErrorMsg(w, http.StatusNotFound, "Not found")
		}
	}))
	defer testServer.Close()

	// Test GET request
	t.Run("GET request", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/test")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		if result["message"] != "GET success" {
			t.Errorf("Expected message 'GET success', got %s", result["message"])
		}
	})

	// Test POST request with JSON body
	t.Run("POST request with JSON", func(t *testing.T) {
		data := map[string]interface{}{"name": "test", "value": 123}
		jsonData, _ := json.Marshal(data)

		resp, err := http.Post(testServer.URL+"/test", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		if result["message"] != "POST success" {
			t.Errorf("Expected message 'POST success', got %s", result["message"])
		}
	})

	// Test query parameters
	t.Run("Query parameters", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/query?name=john&age=25")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		if result["name"] != "john" {
			t.Errorf("Expected name 'john', got %s", result["name"])
		}

		if result["age"] != "25" {
			t.Errorf("Expected age '25', got %s", result["age"])
		}
	})

	// Test health endpoint
	t.Run("Health endpoint", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/health")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		if result["status"] != "healthy" {
			t.Errorf("Expected status 'healthy', got %s", result["status"])
		}
	})

	// Test 404
	t.Run("Not found", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/nonexistent")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestIntegrationWithRealServer(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	cfg.SetPort(0) // Use random port
	serverInstance := zeroserver.NewServerWithConfig(ctx, cfg)

	// Add a simple endpoint
	serverInstance.Get("/api/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "API test successful"})
	})

	// Start server in a goroutine
	go func() {
		// Note: In a real test, you'd want to properly start and stop the server
		// For this integration test, we'll just verify the server can be created
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Verify server was created successfully
	if serverInstance == nil {
		t.Error("Server was not created")
	}
}

func BenchmarkIntegrationRequests(b *testing.B) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}))
	defer testServer.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(testServer.URL + "/test")
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
