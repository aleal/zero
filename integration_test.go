package zero

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleal/zero/pkg/parser"
	"github.com/aleal/zero/pkg/request"
	"github.com/aleal/zero/pkg/response"
	"github.com/aleal/zero/pkg/server"
)

func TestIntegrationServer(t *testing.T) {
	ctx := context.Background()
	z := server.New(ctx)

	z.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "GET success"})
	})

	z.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]any
		if err := parser.ParseJSONBody(r.Body, &data); err != nil {
			response.WriteError(w, http.StatusBadRequest, err)
			return
		}
		response.WriteJSON(w, http.StatusCreated, map[string]any{
			"message": "POST success",
			"data":    data,
		})
	})

	z.Put("/test", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "PUT success"})
	})

	z.Delete("/test", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "DELETE success"})
	})

	z.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		name := request.GetQueryParamOrDefault(r, "name", "default")
		age := request.GetQueryParamOrDefault(r, "age", "0")
		response.WriteJSON(w, http.StatusOK, map[string]string{
			"name": name,
			"age":  age,
		})
	})

	// Use the actual zero mux
	testServer := httptest.NewServer(z.Handler())
	defer testServer.Close()

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

	t.Run("POST request with JSON", func(t *testing.T) {
		data := map[string]any{"name": "test", "value": 123}
		jsonData, _ := json.Marshal(data)

		resp, err := http.Post(testServer.URL+"/test", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		if result["message"] != "POST success" {
			t.Errorf("Expected message 'POST success', got %s", result["message"])
		}
	})

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

		if result["service"] != "zero" {
			t.Errorf("Expected service 'zero', got %s", result["service"])
		}
	})

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

func BenchmarkIntegrationRequests(b *testing.B) {
	ctx := context.Background()
	z := server.New(ctx)

	z.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "success"})
	})

	testServer := httptest.NewServer(z.Handler())
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
