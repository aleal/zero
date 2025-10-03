// Package uuid provides UUID generation functionality for the Zero server.
// It generates cryptographically secure identifiers for request tracking.
package uuid

import (
	"crypto/rand"
	"fmt"
)

// UUID represents a unique identifier string
type UUID string

// GenerateUUID generates a cryptographically secure UUID for request identification.
// This is a simplified version that generates a 4-byte hex string for internal use.
func GenerateUUID() UUID {
	uuid := make([]byte, 4)

	// Read 4 random bytes from the cryptographically secure source.
	if _, err := rand.Read(uuid); err != nil {
		return generateFallbackUUID()
	}

	return UUID(fmt.Sprintf("%x", uuid[0:4]))
}

func generateFallbackUUID() UUID {
	return UUID("00000000")
}
