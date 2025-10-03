// Package response provides utilities for writing HTTP responses in the Zero server.
// It includes functions for writing JSON responses, error responses, and setting
// response headers.
package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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
	SetHeader(w, "Status", strconv.Itoa(statusCode))
	w.Write(data)
}

// InternalServerError writes a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, err error) {
	WriteError(w, http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err))
}

// WriteError writes an error response with the given status code and error
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	WriteErrorMsg(w, statusCode, err.Error())
}

// WriteErrorMsg writes an error response with the given status code and message
func WriteErrorMsg(w http.ResponseWriter, statusCode int, message string) {
	SetHeader(w, "Status", strconv.Itoa(statusCode))
	http.Error(w, message, statusCode)
}

// SetHeader sets a response header with the given key and value
func SetHeader(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}
