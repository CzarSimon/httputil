package crypto_test

import (
	"testing"

	"github.com/CzarSimon/httputil/crypto"
	"github.com/stretchr/testify/assert"
)

func TestRandomBytes(t *testing.T) {
	assert := assert.New(t)
	salts := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		salt, err := crypto.RandomBytes(16)
		assert.NoError(err)
		assert.Len(salt, 16)
		_, ok := salts[string(salt)]
		assert.False(ok)
		salts[string(salt)] = true
	}
}
