package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		expected   map[string]interface{}
	}{
		{
			name:       "simple object",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "hello"},
			expected:   map[string]interface{}{"message": "hello"},
		},
		{
			name:       "complex object",
			statusCode: http.StatusCreated,
			data: map[string]interface{}{
				"id":   1,
				"name": "test",
				"tags": []string{"tag1", "tag2"},
			},
			expected: map[string]interface{}{
				"id":   float64(1),
				"name": "test",
				"tags": []interface{}{"tag1", "tag2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, tt.statusCode, tt.data)

			// Check status code
			if w.Code != tt.statusCode {
				t.Errorf("WriteJSON() status code = %v, want %v", w.Code, tt.statusCode)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSON() content type = %v, want application/json", contentType)
			}

			// Check response body
			var got map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Errorf("WriteJSON() response is not valid JSON: %v", err)
			}

			// Compare the response
			expectedJSON, _ := json.Marshal(tt.expected)
			gotJSON, _ := json.Marshal(got)
			if string(expectedJSON) != string(gotJSON) {
				t.Errorf("WriteJSON() = %v, want %v", string(gotJSON), string(expectedJSON))
			}
		})
	}
}

func TestWriteJSONWithMarshalingError(t *testing.T) {
	// Test WriteJSON with data that cannot be marshaled to JSON
	// This will trigger the else branch in WriteJSON
	w := httptest.NewRecorder()

	// Create a channel, which cannot be marshaled to JSON
	unmarshallableData := make(chan int)
	WriteJSON(w, http.StatusOK, unmarshallableData)

	// Should return 500 Internal Server Error when marshaling fails
	if w.Code != http.StatusInternalServerError {
		t.Errorf("WriteJSON() status code = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	// Should contain error message about internal server error
	body := w.Body.String()
	if body == "" {
		t.Error("WriteJSON() should return error message when marshaling fails")
	}
}

func TestWrite(t *testing.T) {
	t.Run("write with custom content type", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []byte("hello world")
		Write(w, "text/plain", http.StatusOK, data)

		if w.Code != http.StatusOK {
			t.Errorf("Write() status code = %v, want %v", w.Code, http.StatusOK)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "text/plain" {
			t.Errorf("Write() content type = %v, want text/plain", contentType)
		}

		body := w.Body.String()
		if body != "hello world" {
			t.Errorf("Write() body = %v, want hello world", body)
		}
	})
}

func TestWriteError(t *testing.T) {
	t.Run("write error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := http.ErrServerClosed
		WriteError(w, http.StatusInternalServerError, err)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("WriteError() status code = %v, want %v", w.Code, http.StatusInternalServerError)
		}

		// The actual error message includes "http: " prefix
		expected := "http: Server closed\n"
		if w.Body.String() != expected {
			t.Errorf("WriteError() body = %v, want %v", w.Body.String(), expected)
		}
	})
}

func TestWriteErrorMsg(t *testing.T) {
	t.Run("write error message", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteErrorMsg(w, http.StatusBadRequest, "Invalid input")

		if w.Code != http.StatusBadRequest {
			t.Errorf("WriteErrorMsg() status code = %v, want %v", w.Code, http.StatusBadRequest)
		}

		if w.Body.String() != "Invalid input\n" {
			t.Errorf("WriteErrorMsg() body = %v, want Invalid input", w.Body.String())
		}
	})
}

func TestInternalServerError(t *testing.T) {
	t.Run("write internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := http.ErrServerClosed
		InternalServerError(w, err)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("InternalServerError() status code = %v, want %v", w.Code, http.StatusInternalServerError)
		}

		// The actual error message includes "http: " prefix
		expected := "internal server error: http: Server closed\n"
		if w.Body.String() != expected {
			t.Errorf("InternalServerError() body = %v, want %v", w.Body.String(), expected)
		}
	})
}

func TestSetHeader(t *testing.T) {
	t.Run("set custom header", func(t *testing.T) {
		w := httptest.NewRecorder()
		SetHeader(w, "X-Custom-Header", "custom-value")

		headerValue := w.Header().Get("X-Custom-Header")
		if headerValue != "custom-value" {
			t.Errorf("SetHeader() header value = %v, want custom-value", headerValue)
		}
	})
}
