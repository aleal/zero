package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aleal/zero/pkg/log"
	"github.com/aleal/zero/pkg/response"
)

func TestRecovery(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	recoveryHandler := Recovery()(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := log.SetLoggerToContext(req.Context(), log.NewLogger())
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	recoveryHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := `{"message":"success"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestRecoveryWithPanic(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	recoveryHandler := Recovery()(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := log.SetLoggerToContext(req.Context(), log.NewLogger())
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	recoveryHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 after panic, got %d", w.Code)
	}

	// Should return generic message, not panic details
	body := strings.TrimSpace(w.Body.String())
	if body != "internal server error" {
		t.Errorf("Expected generic error message, got %q", body)
	}
}

func TestRecoveryWithoutLoggerInContext(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	recoveryHandler := Recovery()(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// No logger in context — should not double-panic
	w := httptest.NewRecorder()
	recoveryHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRecoveryWithDifferentPanicTypes(t *testing.T) {
	tests := []struct {
		name       string
		panicValue any
	}{
		{"string panic", "test panic"},
		{"error panic", http.ErrServerClosed},
		{"int panic", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				panic(tt.panicValue)
			}

			recoveryHandler := Recovery()(handler)

			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			ctx := log.SetLoggerToContext(req.Context(), log.NewLogger())
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			recoveryHandler(w, req)

			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 500 after panic, got %d", w.Code)
			}

			body := strings.TrimSpace(w.Body.String())
			if body != "internal server error" {
				t.Errorf("Expected generic error, got %q", body)
			}
		})
	}
}

func TestRecoveryWithDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			}

			recoveryHandler := Recovery()(handler)

			req, err := http.NewRequest(method, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			ctx := log.SetLoggerToContext(req.Context(), log.NewLogger())
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			recoveryHandler(w, req)

			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 500 for %s, got %d", method, w.Code)
			}
		})
	}
}

func BenchmarkRecovery(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	recoveryHandler := Recovery()(handler)
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := context.Background()
	ctx = log.SetLoggerToContext(ctx, log.NewLogger())
	req = req.WithContext(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		recoveryHandler(w, req)
	}
}

func BenchmarkRecoveryWithPanic(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	recoveryHandler := Recovery()(handler)
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := context.Background()
	ctx = log.SetLoggerToContext(ctx, log.NewLogger())
	req = req.WithContext(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		recoveryHandler(w, req)
	}
}
