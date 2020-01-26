package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// Common errors
var (
	ErrHashMissmatch = errors.New("hashes do not match")
	ErrInvalidKey    = errors.New("invalid key")
)

// Hasher interface for hashing and verifying hased values.
type Hasher interface {
	Hash(plaintext, salt []byte) ([]byte, error)
	Verify(plaintext, salt, hashtext []byte) error
}

// ScryptHasher hasher implementation using the SCRYPT key derivation function.
type ScryptHasher struct {
	key scryptKey
}

// NewScryptHasher creates a new scrypt hasher
func NewScryptHasher(N, r, p, keyLen int) Hasher {
	return &ScryptHasher{
		key: scryptKey{
			N:      N,
			r:      r,
			p:      p,
			keyLen: keyLen,
		},
	}
}

// DefaultScryptHasher creates scrypt hasher with the default key parameters
func DefaultScryptHasher() Hasher {
	return NewScryptHasher(16384, 8, 1, 32)
}

// Hash derives a key using Scrypt.
func (h *ScryptHasher) Hash(plaintext, salt []byte) ([]byte, error) {
	key, err := h.deriveKey(plaintext, salt, h.key)
	if err != nil {
		return nil, err
	}

	return []byte(key.String()), nil
}

// Verify verifies that plaintext and salt hash to the given hashtext.
func (h *ScryptHasher) Verify(plaintext, salt, hashtext []byte) error {
	key, err := parseScryptKey(hashtext)
	if err != nil {
		return err
	}

	dk, err := h.deriveKey(plaintext, salt, key)
	if err != nil {
		return err
	}

	if !bytes.Equal(key.hash, dk.hash) {
		return ErrHashMissmatch
	}

	return nil
}

func (h *ScryptHasher) deriveKey(plaintext, salt []byte, key scryptKey) (scryptKey, error) {
	hash, err := scrypt.Key(plaintext, salt, key.N, key.r, key.p, key.keyLen)
	if err != nil {
		return scryptKey{}, err
	}

	return scryptKey{
		N:      h.key.N,
		p:      h.key.p,
		r:      h.key.r,
		keyLen: h.key.keyLen,
		hash:   hash,
	}, nil
}

// Sha256Hasher hasher implementation using sha256 as a hash function.
type Sha256Hasher struct{}

// Hash computes hmac of plaintext and salt and calculates a SHA-256 hash on the output.
func (h *Sha256Hasher) Hash(plaintext, salt []byte) ([]byte, error) {
	mac, err := Hmac(plaintext, salt)
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(mac)
	return sum[:], nil
}

// Verify verifies that plaintext and salt hash to the given hashtext.
func (h *Sha256Hasher) Verify(plaintext, salt, hashtext []byte) error {
	hash, err := h.Hash(plaintext, salt)
	if err != nil {
		return err
	}

	if !bytes.Equal(hashtext, hash) {
		return ErrHashMissmatch
	}

	return nil
}

// Hmac computes a hmac using SHA-256.
func Hmac(message, key []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	_, err := h.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hmac: %w", err)
	}

	return h.Sum(nil), nil
}

type scryptKey struct {
	N      int
	p      int
	r      int
	keyLen int
	hash   []byte
}

func (k scryptKey) String() string {
	hash := hex.EncodeToString(k.hash)
	return fmt.Sprintf("SCRYPT$%d$%d$%d$%d$%s", k.N, k.r, k.p, k.keyLen, hash)
}

func parseScryptKey(b []byte) (scryptKey, error) {
	var key scryptKey
	c := strings.Split(string(b), "$")
	if len(c) != 6 {
		return key, fmt.Errorf("%w: unexpected number of components", ErrInvalidKey)
	}

	if c[0] != "SCRYPT" {
		return key, fmt.Errorf("%w: not a SCRYPT key", ErrInvalidKey)
	}

	N, err := parseN(c[1])
	if err != nil {
		return key, err
	}

	r, err := strconv.Atoi(c[2])
	if err != nil {
		return key, fmt.Errorf("%w: invalid r, %v", ErrInvalidKey, err)
	}

	p, err := strconv.Atoi(c[3])
	if err != nil {
		return key, fmt.Errorf("%w: invalid p, %v", ErrInvalidKey, err)
	}

	keyLen, err := strconv.Atoi(c[4])
	if err != nil {
		return key, fmt.Errorf("%w: invalid keylen, %v", ErrInvalidKey, err)
	}
	if keyLen <= 0 {
		return key, fmt.Errorf("%w: invalid key length, must be greater than 0", ErrInvalidKey)
	}

	hash, err := hex.DecodeString(c[5])
	if err != nil {
		return key, fmt.Errorf("%w: invalid hash, %v", ErrInvalidKey, err)
	}
	if keyLen != len(hash) {
		return key, fmt.Errorf("%w: invalid hash, does not match key length", ErrInvalidKey)
	}

	key.N = N
	key.p = p
	key.r = r
	key.keyLen = keyLen
	key.hash = hash

	return key, nil
}

func parseN(str string) (int, error) {
	N, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid N, %v", ErrInvalidKey, err)
	}

	// N = 2^x =>
	// ln(N) = ln(2^x) =>
	// ln(N) = x * ln(2) =>
	// x = ln(N) / ln(2) // x should then be a positive integer
	fN := float64(N)
	exponent := math.Log(fN) / math.Log(2)

	if exponent <= 0 || exponent != float64(int(exponent)) {
		return 0, fmt.Errorf("%w: invalid N, should satisfy N = 2^x where x is a positive integer", ErrInvalidKey)
	}

	return N, nil
}
