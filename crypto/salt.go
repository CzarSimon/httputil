package crypto

import (
	"crypto/rand"
	"fmt"
)

// RandomBytes generates an array of n random bytes.
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return b, nil
}
