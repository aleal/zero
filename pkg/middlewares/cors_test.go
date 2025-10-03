package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleal/zero/pkg/response"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins []string
		method         string
		origin         string
		expectedOrigin string
	}{
		{
			name:           "wildcard origin",
			allowedOrigins: []string{"*"},
			method:         "GET",
			origin:         "https://example.com",
			expectedOrigin: "https://example.com",
		},
		{
			name:           "specific origin allowed",
			allowedOrigins: []string{"https://example.com", "https://test.com"},
			method:         "GET",
			origin:         "https://example.com",
			expectedOrigin: "https://example.com",
		},
		{
			name:           "specific origin not allowed",
			allowedOrigins: []string{"https://example.com"},
			method:         "GET",
			origin:         "https://malicious.com",
			expectedOrigin: "",
		},
		{
			name:           "no origin header",
			allowedOrigins: []string{"https://example.com"},
			method:         "GET",
			origin:         "",
			expectedOrigin: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple handler
			handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
				response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
			}

			// Apply CORS middleware
			corsHandler := CORS(tt.allowedOrigins)(handler)

			// Create test request
			req, err := http.NewRequest(tt.method, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			corsHandler(context.Background(), w, req)

			// Check CORS headers
			accessControlOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if accessControlOrigin != tt.expectedOrigin {
				t.Errorf("Expected Access-Control-Allow-Origin %s, got %s", tt.expectedOrigin, accessControlOrigin)
			}

			// Check other CORS headers
			accessControlMethods := w.Header().Get("Access-Control-Allow-Methods")
			if accessControlMethods == "" {
				t.Error("Access-Control-Allow-Methods header not set")
			}

			accessControlHeaders := w.Header().Get("Access-Control-Allow-Headers")
			if accessControlHeaders == "" {
				t.Error("Access-Control-Allow-Headers header not set")
			}

			// Check response status
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestCORSWithPreflight(t *testing.T) {
	// Create a simple handler
	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Apply CORS middleware
	corsHandler := CORS([]string{"https://example.com"})(handler)

	// Create preflight request
	req, err := http.NewRequest("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call the handler
	corsHandler(context.Background(), w, req)

	// Check preflight response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for preflight, got %d", w.Code)
	}

	// Check CORS headers
	accessControlOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if accessControlOrigin != "https://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin https://example.com, got %s", accessControlOrigin)
	}

	accessControlMethods := w.Header().Get("Access-Control-Allow-Methods")
	if accessControlMethods == "" {
		t.Error("Access-Control-Allow-Methods header not set")
	}

	accessControlHeaders := w.Header().Get("Access-Control-Allow-Headers")
	if accessControlHeaders == "" {
		t.Error("Access-Control-Allow-Headers header not set")
	}
}

func BenchmarkCORS(b *testing.B) {
	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	corsHandler := CORS([]string{"https://example.com"})(handler)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		corsHandler(context.Background(), w, req)
	}
}
