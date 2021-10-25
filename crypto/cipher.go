package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// AES256KeySize key size for AES-256
const AES256KeySize = 32

// Common encryption / decryption errors.
var (
	ErrToShortCiphertext = errors.New("ciphertext too short")
)

// GenerateAESKey generates a random AES-256 key.
func GenerateAESKey() ([]byte, error) {
	key, err := RandomBytes(AES256KeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}

	return key, nil
}

// Cipher symetric chipher interface.
type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// NewCipher creates a new cipher using the default implementation.
func NewCipher(key []byte) Cipher {
	return &AESCipher{
		key: key,
	}
}

// AESCipher symetric cipher implementation using AES-GCM.
type AESCipher struct {
	key []byte
}

// NewAESCipher creates a new AES-GCM cipher with the given key.
func NewAESCipher(key []byte) *AESCipher {
	return &AESCipher{
		key: key,
	}
}

// Encrypt encrypts the provided plaintext with the chiphers key.
func (c *AESCipher) Encrypt(plaintext []byte) ([]byte, error) {
	gcm, err := createGCMChipher(c.key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts the provided ciphertexts with the chiphers key.
func (c *AESCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	gcm, err := createGCMChipher(c.key)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrToShortCiphertext
	}

	return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
}

func createGCMChipher(key []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(c)
}
