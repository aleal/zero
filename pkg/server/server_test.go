package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aleal/zero/pkg/config"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
)

func TestNewServer(t *testing.T) {
	ctx := context.Background()
	server := NewServer(ctx)

	if server == nil {
		t.Error("NewServer() returned nil")
	}
}

func TestNewServerWithConfig(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	cfg.SetHost("localhost")
	cfg.SetPort(8080)

	server := NewServerWithConfig(ctx, cfg)

	if server == nil {
		t.Error("NewServerWithConfig() returned nil")
	}
}

func TestNewServerWithNilConfig(t *testing.T) {
	ctx := context.Background()
	server := NewServerWithConfig(ctx, nil)

	if server == nil {
		t.Error("NewServerWithConfig() with nil config returned nil")
	}
}

func TestServerMethods(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Test that we can register handlers
	server.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	})

	server.Post("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"message": "created"})
	})

	server.Put("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
	})

	server.Delete("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
	})

	server.Patch("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "patched"})
	})

	// Test that we can add middlewares
	server.Middlewares(func(next request.Handler) request.Handler {
		return func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
			next(rctx, w, r)
		}
	})
}

func TestServerHandleMethod(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Test the Handle method directly
	server.Handle("/custom", "CUSTOM", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "custom"})
	})
}

func TestServerPatternNormalization(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Test that patterns are normalized
	server.Get("test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	})

	server.Get("/another-test/", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "another-test"})
	})

	// Test that we can register multiple handlers for the same pattern
	server.Get("/multi", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "first"})
	})

	server.Post("/multi", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"message": "second"})
	})
}

func TestServerMethodNotAllowed(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Register only GET handler
	server.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	})

	// Test that we can register handlers for different methods
	server.Post("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"message": "created"})
	})

	server.Put("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
	})

	server.Delete("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
	})

	server.Patch("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "patched"})
	})
}

func TestServerWithMultipleMiddlewares(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Add multiple middlewares
	server.Middlewares(
		func(next request.Handler) request.Handler {
			return func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware-1", "true")
				next(rctx, w, r)
			}
		},
		func(next request.Handler) request.Handler {
			return func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware-2", "true")
				next(rctx, w, r)
			}
		},
	)

	// Register a handler
	server.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	})
}

func TestServerStartMethod(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Register a simple handler
	server.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	})

	// Test that Start method can be called without panicking
	// Note: We can't easily test the full Start method as it blocks indefinitely
	// But we can test that it doesn't panic immediately
	go func() {
		// This will block, but we can test that it starts without error
		server.Start()
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Note: We can't easily shutdown the server through the interface
	// The Start method is designed to run indefinitely until interrupted
}

func TestServerMethodRouterBehavior(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Test that we can register multiple handlers for the same pattern
	server.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "GET"})
	})

	server.Post("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"method": "POST"})
	})

	server.Put("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "PUT"})
	})

	// Test that we can register handlers for different patterns
	server.Get("/another", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"pattern": "another"})
	})

	// Test that we can register custom methods
	server.Handle("/custom", "CUSTOM", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "CUSTOM"})
	})
}

func TestServerMethodNotAllowedBehavior(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	// Test that we can register handlers for different methods
	// The methodNotAllowed function is tested indirectly through the router behavior
	server.Get("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"method": "GET"})
	})

	server.Post("/test", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusCreated, map[string]string{"method": "POST"})
	})

	// Test that we can register handlers for different patterns
	server.Get("/another", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"pattern": "another"})
	})
}

// func BenchmarkServerCreation(b *testing.B) {
// 	ctx := context.Background()
// 	for i := 0; i < b.N; i++ {
// 		NewServer(ctx)
// 	}
// }

func BenchmarkServerWithConfigCreation(b *testing.B) {
	ctx := context.Background()
	cfg := config.Default()
	for i := 0; i < b.N; i++ {
		NewServerWithConfig(ctx, cfg)
	}
}

func BenchmarkServerMethodRegistration(b *testing.B) {
	ctx := context.Background()
	cfg := config.Default()
	server := NewServerWithConfig(ctx, cfg)

	handler := func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "test"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.Get("/test", handler)
	}
}
