package middlewares

import (
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
		expectHeaders  bool // whether Allow-Methods/Allow-Headers should be set
	}{
		{
			name:           "wildcard origin",
			allowedOrigins: []string{"*"},
			method:         "GET",
			origin:         "https://example.com",
			expectedOrigin: "*",
			expectHeaders:  true,
		},
		{
			name:           "specific origin allowed",
			allowedOrigins: []string{"https://example.com", "https://test.com"},
			method:         "GET",
			origin:         "https://example.com",
			expectedOrigin: "https://example.com",
			expectHeaders:  true,
		},
		{
			name:           "specific origin not allowed",
			allowedOrigins: []string{"https://example.com"},
			method:         "GET",
			origin:         "https://malicious.com",
			expectedOrigin: "",
			expectHeaders:  false,
		},
		{
			name:           "no origin header",
			allowedOrigins: []string{"https://example.com"},
			method:         "GET",
			origin:         "",
			expectedOrigin: "",
			expectHeaders:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
			}

			corsHandler := CORS(tt.allowedOrigins)(handler)

			req, err := http.NewRequest(tt.method, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			corsHandler(w, req)

			accessControlOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if accessControlOrigin != tt.expectedOrigin {
				t.Errorf("Expected Access-Control-Allow-Origin %q, got %q", tt.expectedOrigin, accessControlOrigin)
			}

			hasMethods := w.Header().Get("Access-Control-Allow-Methods") != ""
			hasHeaders := w.Header().Get("Access-Control-Allow-Headers") != ""

			if tt.expectHeaders && !hasMethods {
				t.Error("Access-Control-Allow-Methods header not set")
			}
			if tt.expectHeaders && !hasHeaders {
				t.Error("Access-Control-Allow-Headers header not set")
			}
			if !tt.expectHeaders && hasMethods {
				t.Error("Access-Control-Allow-Methods should not be set for unmatched origin")
			}
			if !tt.expectHeaders && hasHeaders {
				t.Error("Access-Control-Allow-Headers should not be set for unmatched origin")
			}

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestCORSCredentialsOnlyForExplicitOrigins(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Wildcard should NOT set credentials
	corsHandler := CORS([]string{"*"})(handler)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	corsHandler(w, req)

	if creds := w.Header().Get("Access-Control-Allow-Credentials"); creds != "" {
		t.Errorf("Wildcard CORS should not set credentials, got %q", creds)
	}

	// Explicit origin SHOULD set credentials
	corsHandler2 := CORS([]string{"https://example.com"})(handler)
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "https://example.com")
	w2 := httptest.NewRecorder()
	corsHandler2(w2, req2)

	if creds := w2.Header().Get("Access-Control-Allow-Credentials"); creds != "true" {
		t.Errorf("Explicit CORS should set credentials, got %q", creds)
	}

	if vary := w2.Header().Get("Vary"); vary != "Origin" {
		t.Errorf("Explicit CORS should set Vary: Origin, got %q", vary)
	}
}

func TestCORSMaxAge(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	corsHandler := CORS([]string{"https://example.com"})(handler)

	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	corsHandler(w, req)

	maxAge := w.Header().Get("Access-Control-Max-Age")
	if maxAge != "600" {
		t.Errorf("Expected Access-Control-Max-Age 600, got %q", maxAge)
	}
}

func TestCORSWithPreflight(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	corsHandler := CORS([]string{"https://example.com"})(handler)

	req, err := http.NewRequest("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	w := httptest.NewRecorder()
	corsHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for preflight, got %d", w.Code)
	}

	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "https://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin https://example.com, got %s", origin)
	}

	if methods := w.Header().Get("Access-Control-Allow-Methods"); methods == "" {
		t.Error("Access-Control-Allow-Methods header not set")
	}

	if headers := w.Header().Get("Access-Control-Allow-Headers"); headers == "" {
		t.Error("Access-Control-Allow-Headers header not set")
	}
}

func BenchmarkCORS(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	corsHandler := CORS([]string{"https://example.com"})(handler)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		corsHandler(w, req)
	}
}
