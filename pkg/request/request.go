package request

import (
	"context"
	"net/http"

	"github.com/aleal/zero/pkg/parser"
)

// Handler is a function that handles a request
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// GetPathParam gets a path parameter
func GetPathParam(r *http.Request, key string) string {
	return GetPathParamOrDefault(r, key, "")
}

// GetPathParamOrDefault gets a path parameter with default value
func GetPathParamOrDefault(r *http.Request, key, defaultValue string) string {
	if value := r.PathValue(key); value != "" {
		return parser.SanitizeString(value)
	}
	return defaultValue
}

// GetQueryParam gets a query parameter
func GetQueryParam(r *http.Request, key string) string {
	return GetQueryParamOrDefault(r, key, "")
}

// GetQueryParamOrDefault gets a query parameter with default value
func GetQueryParamOrDefault(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return parser.SanitizeString(value)
	}
	return defaultValue
}
