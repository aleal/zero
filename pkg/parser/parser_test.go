package parser

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with control characters",
			input:    "hello\x00world\x01\x02",
			expected: "helloworld",
		},
		{
			name:     "string with whitespace",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "string with tabs and newlines",
			input:    "hello\tworld\n",
			expected: "hello\tworld", // Newlines are trimmed by TrimSpace
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	t.Run("parse int", func(t *testing.T) {
		var result int
		err := ParseString("123", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != 123 {
			t.Errorf("ParseString() = %v, want %v", result, 123)
		}
	})

	t.Run("parse int32", func(t *testing.T) {
		var result int32
		err := ParseString("123", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != 123 {
			t.Errorf("ParseString() = %v, want %v", result, 123)
		}
	})

	t.Run("parse int64", func(t *testing.T) {
		var result int64
		err := ParseString("123456789", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != 123456789 {
			t.Errorf("ParseString() = %v, want %v", result, 123456789)
		}
	})

	t.Run("parse float32", func(t *testing.T) {
		var result float32
		err := ParseString("123.45", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != 123.45 {
			t.Errorf("ParseString() = %v, want %v", result, 123.45)
		}
	})

	t.Run("parse float64", func(t *testing.T) {
		var result float64
		err := ParseString("123.456789", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != 123.456789 {
			t.Errorf("ParseString() = %v, want %v", result, 123.456789)
		}
	})

	t.Run("parse bool", func(t *testing.T) {
		var result bool
		err := ParseString("true", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != true {
			t.Errorf("ParseString() = %v, want %v", result, true)
		}
	})

	t.Run("parse string", func(t *testing.T) {
		var result string
		err := ParseString("hello world", &result)
		if err != nil {
			t.Errorf("ParseString() error = %v", err)
		}
		if result != "hello world" {
			t.Errorf("ParseString() = %v, want %v", result, "hello world")
		}
	})

	t.Run("parse invalid int", func(t *testing.T) {
		var result int
		err := ParseString("not a number", &result)
		if err == nil {
			t.Error("ParseString() should return error for invalid int")
		}
	})

	t.Run("parse invalid int32", func(t *testing.T) {
		var result int32
		err := ParseString("not a number", &result)
		if err == nil {
			t.Error("ParseString() should return error for invalid int32")
		}
	})

	t.Run("parse invalid int64", func(t *testing.T) {
		var result int64
		err := ParseString("not a number", &result)
		if err == nil {
			t.Error("ParseString() should return error for invalid int64")
		}
	})

	t.Run("parse invalid float32", func(t *testing.T) {
		var result float32
		err := ParseString("not a number", &result)
		if err == nil {
			t.Error("ParseString() should return error for invalid float32")
		}
	})

	t.Run("parse invalid float64", func(t *testing.T) {
		var result float64
		err := ParseString("not a number", &result)
		if err == nil {
			t.Error("ParseString() should return error for invalid float64")
		}
	})

	t.Run("parse invalid bool", func(t *testing.T) {
		var result bool
		err := ParseString("not a bool", &result)
		if err == nil {
			t.Error("ParseString() should return error for invalid bool")
		}
	})
}
