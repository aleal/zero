package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// Get the health check handler
	handler := HealthCheckHandler()

	// Create a test request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the handler
	handler(context.Background(), w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Check that the response contains the expected fields
	if response["service"] != "zero" {
		t.Errorf("Expected service 'zero', got %s", response["service"])
	}

	if response["version"] != "0.0.1" {
		t.Errorf("Expected version '0.0.1', got %s", response["version"])
	}

	// Check that uptime is present and is a string
	if uptime, ok := response["uptime"].(string); !ok || uptime == "" {
		t.Errorf("Expected uptime to be a non-empty string, got %v", response["uptime"])
	}
}

func TestHealthCheckHandlerWithDifferentMethods(t *testing.T) {
	handler := HealthCheckHandler()

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler(context.Background(), w, req)

			// Health check should work with any method
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}
		})
	}
}

func BenchmarkHealthCheckHandler(b *testing.B) {
	handler := HealthCheckHandler()
	req, _ := http.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler(context.Background(), w, req)
	}
}
