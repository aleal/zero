package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ParseJSONBody parses JSON from request body
func ParseJSONBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// ParseType defines the types that can be parsed from string values
type ParseType interface {
	int | int32 | int64 | float32 | float64 | bool | string
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(s string) string {
	// Remove null bytes and control characters
	s = strings.Map(func(r rune) rune {
		if r < 32 && r != 9 && r != 10 && r != 13 {
			return -1
		}
		return r
	}, s)

	return strings.TrimSpace(s)
}

// ParseString parses a string into a target type
func ParseString[T ParseType](valueStr string, target *T) error {
	valueStr = SanitizeString(valueStr)
	switch any(*target).(type) {
	case int:
		val, err := strconv.ParseInt(valueStr, 10, 0)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as int: %w", valueStr, err)
		}
		*target = any(int(val)).(T)
	case int32:
		val, err := strconv.ParseInt(valueStr, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as int32: %w", valueStr, err)
		}
		*target = any(int32(val)).(T)
	case int64:
		val, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as int64: %w", valueStr, err)
		}
		*target = any(int64(val)).(T) // val is already int64, so direct type assertion is safe.
	case float32:
		val, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as float32: %w", valueStr, err)
		}
		*target = any(float32(val)).(T)
	case float64:
		val, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as float64: %w", valueStr, err)
		}
		*target = any(float64(val)).(T) // val is already float64, so direct type assertion is safe.
	case bool:
		val, err := strconv.ParseBool(valueStr)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as bool: %w", valueStr, err)
		}
		*target = any(bool(val)).(T) // val is already bool, so direct type assertion is safe.
	case string:
		*target = any(valueStr).(T)
	default:
		return fmt.Errorf("unsupported type: %T", *target)
	}
	return nil
}
