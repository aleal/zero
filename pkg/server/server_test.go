package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aleal/zero/pkg/config"
	"github.com/aleal/zero/pkg/response"
)

func TestNewServer(t *testing.T) {
	server := New(context.Background())
	if server == nil {
		t.Error("New() returned nil")
	}
}

func TestNewServerWithConfig(t *testing.T) {
	cfg := config.Load()
	cfg.SetHost("0.0.0.0")
	cfg.SetPort(9090)

	server := New(context.Background(), WithConfig(cfg))
	if server == nil {
		t.Error("New(WithConfig) returned nil")
	}
}

func TestNewServerWithBadEnv(t *testing.T) {
	t.Setenv("ZERO_PORT", "not-a-number")

	server := New(context.Background())
	if server == nil {
		t.Error("New() should return server even with bad env")
	}
}

// --- HTTP round-trip tests ---

func TestMethodRouterRoundTrip(t *testing.T) {
	srv := New(context.Background())
	srv.Get("/items", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"action": "list"})
	})
	srv.Post("/items", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"action": "create"})
	})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"GET registered", http.MethodGet, "/items", http.StatusOK},
		{"POST registered", http.MethodPost, "/items", http.StatusCreated},
		{"PUT not registered returns 405", http.MethodPut, "/items", http.StatusMethodNotAllowed},
		{"DELETE not registered returns 405", http.MethodDelete, "/items", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Do: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestMethodNotAllowedHasAllowHeader(t *testing.T) {
	srv := New(context.Background())
	srv.Get("/resource", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, nil)
	})
	srv.Post("/resource", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, nil)
	})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/resource", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow == "" {
		t.Error("Allow header is empty on 405 response")
	}
}

func TestHandlerReturnsWorkingMux(t *testing.T) {
	srv := New(context.Background())
	srv.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"pong": "true"})
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestHealthCheckEndpoint(t *testing.T) {
	srv := New(context.Background())

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestShutdown(t *testing.T) {
	srv := New(context.Background(), WithPort(0))

	z := srv.(*zero)
	ln, err := net.Listen("tcp", z.Server.Addr)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}

	go z.Server.Serve(ln)

	resp, err := http.Get("http://" + ln.Addr().String() + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	_, err = http.Get("http://" + ln.Addr().String() + "/health")
	if err == nil {
		t.Error("expected connection error after shutdown")
	}
}

func TestMultipleMethodsSamePattern(t *testing.T) {
	srv := New(context.Background())
	srv.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "GET"})
	})
	srv.Post("/api/data", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"method": "POST"})
	})
	srv.Put("/api/data", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "PUT"})
	})
	srv.Delete("/api/data", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusNoContent, nil)
	})
	srv.Patch("/api/data", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "PATCH"})
	})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	methods := map[string]int{
		http.MethodGet:    http.StatusOK,
		http.MethodPost:   http.StatusCreated,
		http.MethodPut:    http.StatusOK,
		http.MethodDelete: http.StatusNoContent,
		http.MethodPatch:  http.StatusOK,
	}
	for method, wantStatus := range methods {
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, ts.URL+"/api/data", nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Do: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, wantStatus)
			}
		})
	}
}

func TestCustomHandleMethod(t *testing.T) {
	srv := New(context.Background())
	srv.Handle("/custom", "CUSTOM", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "CUSTOM"})
	})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	req, _ := http.NewRequest("CUSTOM", ts.URL+"/custom", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestPatternNormalization(t *testing.T) {
	srv := New(context.Background())
	srv.Get("no-leading-slash", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"ok": "true"})
	})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/no-leading-slash")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func BenchmarkServerCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(context.Background())
	}
}

func BenchmarkServerMethodRegistration(b *testing.B) {
	server := New(context.Background())

	handler := func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.Get(fmt.Sprintf("/test-%d", i), handler)
	}
}
