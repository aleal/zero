package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/response"
)

func TestRecovery(t *testing.T) {
	// Create a simple handler
	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Apply recovery middleware
	recoveryHandler := Recovery()(handler)

	// Create test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Create context with logger
	ctx := context.Background()
	logger := log.NewLogger()
	ctx = log.SetLoggerToContext(ctx, logger)

	// Call the handler
	recoveryHandler(ctx, w, req)

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

func TestRecoveryWithPanic(t *testing.T) {
	// Create a handler that panics
	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	// Apply recovery middleware
	recoveryHandler := Recovery()(handler)

	// Create test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Create context with logger
	ctx := context.Background()
	logger := log.NewLogger()
	ctx = log.SetLoggerToContext(ctx, logger)

	// Call the handler - this should not panic due to recovery middleware
	recoveryHandler(ctx, w, req)

	// Check that we got an internal server error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 after panic, got %d", w.Code)
	}

	// Check that we got an error response
	if w.Body.String() == "" {
		t.Error("Expected error response body, got empty")
	}
}

func TestRecoveryWithDifferentPanicTypes(t *testing.T) {
	tests := []struct {
		name        string
		panicValue  interface{}
		expectedMsg string
	}{
		{"string panic", "test panic", "test panic"},
		{"error panic", http.ErrServerClosed, "server closed"},
		{"int panic", 42, "42"},
		{"nil panic", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a handler that panics with the specified value
			handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
				panic(tt.panicValue)
			}

			// Apply recovery middleware
			recoveryHandler := Recovery()(handler)

			// Create test request
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create context with logger
			ctx := context.Background()
			logger := log.NewLogger()
			ctx = log.SetLoggerToContext(ctx, logger)

			// Call the handler - this should not panic due to recovery middleware
			recoveryHandler(ctx, w, req)

			// Check that we got an internal server error
			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 500 after panic, got %d", w.Code)
			}

			// Check that we got an error response
			if w.Body.String() == "" {
				t.Error("Expected error response body, got empty")
			}
		})
	}
}

func TestRecoveryWithDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Create a handler that panics
			handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			}

			// Apply recovery middleware
			recoveryHandler := Recovery()(handler)

			// Create test request
			req, err := http.NewRequest(method, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create context with logger
			ctx := context.Background()
			logger := log.NewLogger()
			ctx = log.SetLoggerToContext(ctx, logger)

			// Call the handler
			recoveryHandler(ctx, w, req)

			// Check that we got an internal server error
			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 500 for %s, got %d", method, w.Code)
			}
		})
	}
}

func BenchmarkRecovery(b *testing.B) {
	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	recoveryHandler := Recovery()(handler)
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create context with logger
	ctx := context.Background()
	logger := log.NewLogger()
	ctx = log.SetLoggerToContext(ctx, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		recoveryHandler(ctx, w, req)
	}
}

func BenchmarkRecoveryWithPanic(b *testing.B) {
	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	recoveryHandler := Recovery()(handler)
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create context with logger
	ctx := context.Background()
	logger := log.NewLogger()
	ctx = log.SetLoggerToContext(ctx, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		recoveryHandler(ctx, w, req)
	}
}
