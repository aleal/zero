package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

func TestChain(t *testing.T) {
	// Create a simple handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Create middleware that adds a header
	middleware1 := func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "true")
			next(w, r)
		}
	}

	// Create middleware that adds another header
	middleware2 := func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "true")
			next(w, r)
		}
	}

	// Chain the middlewares
	chainedHandler := Chain(handler, middleware1, middleware2)

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the chained handler
	chainedHandler(w, req)

	// Check that both middleware headers were added
	if w.Header().Get("X-Middleware-1") != "true" {
		t.Error("Middleware 1 header not found")
	}

	if w.Header().Get("X-Middleware-2") != "true" {
		t.Error("Middleware 2 header not found")
	}

	// Check that the response was successful
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestChainWithNoMiddlewares(t *testing.T) {
	// Create a simple handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Chain with no middlewares
	chainedHandler := Chain(handler)

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the chained handler
	chainedHandler(w, req)

	// Check that the response was successful
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestChainOrder(t *testing.T) {
	var executionOrder []string

	// Create a handler that records execution
	handler := func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	// Create middleware that records execution order
	middleware1 := func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		}
	}

	middleware2 := func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		}
	}

	// Chain the middlewares
	chainedHandler := Chain(handler, middleware1, middleware2)

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the chained handler
	chainedHandler(w, req)

	// Check execution order - middlewares are applied in reverse order
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected %d execution steps, got %d", len(expectedOrder), len(executionOrder))
		return
	}

	for i, step := range expectedOrder {
		if executionOrder[i] != step {
			t.Errorf("Step %d: expected %s, got %s", i, step, executionOrder[i])
		}
	}
}

func BenchmarkChain(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	}

	middleware := func(next request.Handler) request.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chainedHandler := Chain(handler, middleware, middleware, middleware)
		chainedHandler(w, req)
	}
}
