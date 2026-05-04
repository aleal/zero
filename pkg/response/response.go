// Package response provides utilities for writing HTTP responses in the Zero server.
// It includes functions for writing JSON responses, error responses, and setting
// response headers.
package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.written {
		return
	}
	rw.statusCode = code
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK) // Default to 200 if Write is called first
	}
	return rw.ResponseWriter.Write(b)
}

// Unwrap returns the underlying ResponseWriter so http.NewResponseController
// can reach Flusher, Hijacker, etc.
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// WriteJSON writes a JSON response with the given status code and data
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	if jsonData, err := json.Marshal(data); err == nil {
		Write(w, "application/json", statusCode, jsonData)
	} else {
		InternalServerError(w, err)
	}
}

// Write writes a response with the given content type, status code, and data
func Write(w http.ResponseWriter, contentType string, statusCode int, data []byte) {
	SetHeader(w, "Content-Type", contentType)
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		slog.Debug("response write failed", slog.Any("error", err))
	}
}

// InternalServerError writes a generic 500 response to the client and logs the real error server-side
func InternalServerError(w http.ResponseWriter, err error) {
	slog.Error("internal server error", slog.Any("error", err))
	WriteErrorMsg(w, http.StatusInternalServerError, "internal server error")
}

// WriteError writes an error response with the given status code and error
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	WriteErrorMsg(w, statusCode, err.Error())
}

// WriteErrorMsg writes an error response with the given status code and message
func WriteErrorMsg(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

// SetHeader sets a response header with the given key and value
func SetHeader(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}

func StatusCode(w http.ResponseWriter) int {
	if rw, ok := w.(*responseWriter); ok {
		return rw.statusCode
	}
	return http.StatusOK
}
