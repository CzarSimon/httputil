package crypto_test

import (
	"testing"

	"github.com/CzarSimon/httputil/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAESCipher(t *testing.T) {
	assert := assert.New(t)

	data := "my secret data"
	key := []byte("some-secret-key-")

	var cipher crypto.Cipher = crypto.NewAESCipher(key)

	ciphertext, err := cipher.Encrypt([]byte(data))
	assert.NoError(err)
	assert.NotNil(ciphertext)
	assert.NotEqual(data, string(ciphertext))

	plaintext, err := cipher.Decrypt(ciphertext)
	assert.NoError(err)
	assert.NotNil(plaintext)
	assert.Equal(data, string(plaintext))
}
