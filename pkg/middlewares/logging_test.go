package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/response"
)

func TestLogging(t *testing.T) {
	// Create a simple handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Apply logging middleware
	loggingHandler := Logging(log.NewLogger())(handler)

	// Create test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Create context with logger
	logger := log.NewLogger()
	ctx := log.SetLoggerToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	// Call the handler
	loggingHandler(w, req)

	// Check that the response was successful
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check that the response body is correct
	expectedBody := `{"message":"success"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestLoggingWithDifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"success", http.StatusOK},
		{"created", http.StatusCreated},
		{"bad request", http.StatusBadRequest},
		{"not found", http.StatusNotFound},
		{"internal server error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a handler that returns the specified status code
			handler := func(w http.ResponseWriter, r *http.Request) {
				response.WriteJSON(w, tt.statusCode, map[string]string{"status": tt.name})
			}

			// Apply logging middleware
			loggingHandler := Logging(log.NewLogger())(handler)

			// Create test request
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create context with logger
			logger := log.NewLogger()
			ctx := log.SetLoggerToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			// Call the handler
			loggingHandler(w, req)

			// Check that the status code is correct
			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestLoggingWithDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Create a simple handler
			handler := func(w http.ResponseWriter, r *http.Request) {
				response.WriteJSON(w, http.StatusOK, map[string]string{"method": method})
			}

			// Apply logging middleware
			loggingHandler := Logging(log.NewLogger())(handler)

			// Create test request
			req, err := http.NewRequest(method, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create context with logger
			logger := log.NewLogger()
			ctx := log.SetLoggerToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			// Call the handler
			loggingHandler(w, req)

			// Check that the response was successful
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}
		})
	}
}

func BenchmarkLogging(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	loggingHandler := Logging(log.NewLogger())(handler)
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create context with logger
	ctx := context.Background()
	logger := log.NewLogger()
	ctx = log.SetLoggerToContext(ctx, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		loggingHandler(w, req)
	}
}
