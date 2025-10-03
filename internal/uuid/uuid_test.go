package uuid

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	// Test that UUID generation works
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()

	// Check that UUIDs are not empty
	if uuid1 == "" {
		t.Error("GenerateUUID() returned empty string")
	}
	if uuid2 == "" {
		t.Error("GenerateUUID() returned empty string")
	}

	// Check that UUIDs are different (very unlikely to be the same)
	if uuid1 == uuid2 {
		t.Error("GenerateUUID() returned same UUID twice")
	}

	// Check that UUIDs are valid hex strings
	if !isValidHex(string(uuid1)) {
		t.Errorf("GenerateUUID() returned invalid hex string: %s", uuid1)
	}
	if !isValidHex(string(uuid2)) {
		t.Errorf("GenerateUUID() returned invalid hex string: %s", uuid2)
	}

	// Check that UUIDs have the expected length (4 bytes = 8 hex chars)
	if len(string(uuid1)) != 8 {
		t.Errorf("GenerateUUID() returned UUID with wrong length: %d", len(string(uuid1)))
	}
	if len(string(uuid2)) != 8 {
		t.Errorf("GenerateUUID() returned UUID with wrong length: %d", len(string(uuid2)))
	}
}

func TestGenerateFallbackUUID(t *testing.T) {
	uuid := generateFallbackUUID()

	// Check that fallback UUID is not empty
	if uuid == "" {
		t.Error("generateFallbackUUID() returned empty string")
	}

	// Check that fallback UUID is the expected value
	if uuid != "00000000" {
		t.Errorf("generateFallbackUUID() returned %s, expected 00000000", uuid)
	}
}

func TestUUIDType(t *testing.T) {
	// Test that UUID is a string type
	var uuid UUID = "test"
	if string(uuid) != "test" {
		t.Error("UUID type conversion failed")
	}
}

// Helper function to check if a string is valid hex
func isValidHex(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func BenchmarkGenerateUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUUID()
	}
}
