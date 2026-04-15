package request

import (
	"net/http"
	"testing"
)

func TestGetPathParam(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{"existing path param", "id", "123", "123"},
		{"empty value returns empty", "id", "", ""},
		{"non-existing param", "name", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.value != "" {
				req.SetPathValue(tt.key, tt.value)
			}
			got := GetPathParam(req, tt.key)
			if got != tt.expected {
				t.Errorf("GetPathParam() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetQueryParam(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		key      string
		expected string
	}{
		{
			name:     "existing query param",
			query:    "?name=test&value=123",
			key:      "name",
			expected: "test",
		},
		{
			name:     "non-existing query param",
			query:    "?name=test",
			key:      "value",
			expected: "",
		},
		{
			name:     "empty query",
			query:    "",
			key:      "name",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			got := GetQueryParam(req, tt.key)
			if got != tt.expected {
				t.Errorf("GetQueryParam() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetPathParamOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{"existing param returns value", "id", "456", "default", "456"},
		{"missing param returns default", "name", "", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.value != "" {
				req.SetPathValue(tt.key, tt.value)
			}
			got := GetPathParamOrDefault(req, tt.key, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("GetPathParamOrDefault() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetQueryParamOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "existing query param",
			query:        "?name=test&value=123",
			key:          "name",
			defaultValue: "default",
			expected:     "test",
		},
		{
			name:         "non-existing query param",
			query:        "?name=test",
			key:          "value",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			got := GetQueryParamOrDefault(req, tt.key, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("GetQueryParamOrDefault() = %v, want %v", got, tt.expected)
			}
		})
	}
}
