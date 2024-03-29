package crypto_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/CzarSimon/httputil/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSha256Hasher(t *testing.T) {
	assert := assert.New(t)

	var hasher crypto.Hasher = &crypto.Sha256Hasher{}

	plaintext := "go24/crypto.Hahser super secret key"
	salt := []byte("a11b5a2b-db3c-4105-8c1a-0918f10eda21")
	expectedHash := "42b226d4c3a2ed010052d70475d01b625152011e0261c521b5d9a7ff662a61c9"

	expectedHashBytes, err := hex.DecodeString(expectedHash)
	assert.NoError(err)

	hash, err := hasher.Hash([]byte(plaintext), salt)
	assert.NoError(err)
	assert.Equal(expectedHash, hex.EncodeToString(hash))

	err = hasher.Verify([]byte(plaintext), salt, expectedHashBytes)
	assert.NoError(err)

	err = hasher.Verify([]byte("some other plaintext value"), salt, expectedHashBytes)
	assert.Equal(crypto.ErrHashMissmatch, err)

	err = hasher.Verify([]byte(plaintext), salt, []byte("some other hash value"))
	assert.Equal(crypto.ErrHashMissmatch, err)

	err = hasher.Verify([]byte(plaintext), []byte("some other salt value"), expectedHashBytes)
	assert.Equal(crypto.ErrHashMissmatch, err)
}

func TestScryptHasher(t *testing.T) {
	assert := assert.New(t)

	var hasher crypto.Hasher = crypto.DefaultScryptHasher()

	plaintext := "625181dbfb5c6100cdacd97f3ba32ab4"
	salt, err := hex.DecodeString("478c1d403dec20707cf487f81c06d646")
	assert.NoError(err)
	expectedHash := "SCRYPT$16384$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739"
	expectedHashBytes := []byte(expectedHash)

	hash, err := hasher.Hash([]byte(plaintext), salt)
	assert.NoError(err)
	assert.Equal(expectedHash, string(hash))

	validKeys := []string{
		expectedHash,
		"SCRYPT$8192$8$1$32$b0aaba3db0727f3cff18088b64aa69eea8c19d911fd1c66e9e54cc7bde7ea1fe",
		"SCRYPT$32768$8$1$32$76986441c8bccd072eee0076b959835f3663cd6102c0b1b7683322c4f0a8d0ba",
		"SCRYPT$8192$4$1$32$7ffd6a49fe8079e352ef8fd0f5abfe27c23019f6bc09918d0b6a158d46e7a816",
		"SCRYPT$8192$4$1$16$7ffd6a49fe8079e352ef8fd0f5abfe27",
		"SCRYPT$8192$4$1$64$7ffd6a49fe8079e352ef8fd0f5abfe27c23019f6bc09918d0b6a158d46e7a8164be4931396266cdc7d89a178e4c9eae3ad602a8c5e3d57fc0dd8b459828b5c44",
	}

	for i, key := range validKeys {
		err = hasher.Verify([]byte(plaintext), salt, []byte(key))
		assert.NoError(err, fmt.Sprintf("Key number %d should be valid", i+1))
	}

	err = hasher.Verify([]byte("some other plaintext value"), salt, expectedHashBytes)
	assert.Equal(crypto.ErrHashMissmatch, err)

	wrongHash := []byte("SCRYPT$16384$8$1$32$44b2245a75516a5c3dc2ebe016a2ba564f0f20cb678076d89d0adcde599206e0")
	err = hasher.Verify([]byte(plaintext), salt, wrongHash)
	assert.Equal(crypto.ErrHashMissmatch, err)

	wrongParams2 := []byte("SCRYPT$16384$4$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739")
	err = hasher.Verify([]byte(plaintext), salt, wrongParams2)
	assert.Equal(crypto.ErrHashMissmatch, err)

	wrongParams3 := []byte("SCRYPT$8192$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739")
	err = hasher.Verify([]byte(plaintext), salt, wrongParams3)
	assert.Equal(crypto.ErrHashMissmatch, err)

	err = hasher.Verify([]byte(plaintext), []byte("some other salt value"), expectedHashBytes)
	assert.Equal(crypto.ErrHashMissmatch, err)
}

func TestParseScryptKey(t *testing.T) {
	assert := assert.New(t)

	var hasher crypto.Hasher = crypto.DefaultScryptHasher()

	plaintext := "625181dbfb5c6100cdacd97f3ba32ab4"
	salt, err := hex.DecodeString("478c1d403dec20707cf487f81c06d646")
	assert.NoError(err)

	invalidKeys := []string{
		"SCRYPT$16384$8$1$64$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"SCRYPT$16384$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a6973",
		"SCRYPT$16384$8$1$0$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"SCRYPT$16385$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"SCRYPT$20000$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"SCRYPT$16383$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"PBKDF2$16384$8$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"SCRYPT$16384$not-an-int$1$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
		"SCRYPT$16384$8$not-an-int$32$b8059f5d26826ef3af0faa424a8fc0f51f80bd62aa46ada056f7174e08a69739",
	}

	for i, key := range invalidKeys {
		err := hasher.Verify([]byte(plaintext), salt, []byte(key))
		assert.True(errors.Is(err, crypto.ErrInvalidKey), fmt.Sprintf("Key number %d should be invalid", i+1))
	}
}

func TestArgon2Hasher(t *testing.T) {
	assert := assert.New(t)

	var hasher crypto.Hasher = crypto.DefaultArgon2Hasher()

	plaintext := "625181dbfb5c6100cdacd97f3ba32ab4"
	salt, err := hex.DecodeString("478c1d403dec20707cf487f81c06d646")
	assert.NoError(err)
	expectedHash := "ARGON2ID$1$65536$4$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4"
	expectedHashBytes := []byte(expectedHash)

	hash, err := hasher.Hash([]byte(plaintext), salt)
	assert.NoError(err)
	assert.Equal(expectedHash, string(hash))

	validKeys := []string{
		expectedHash,
		"ARGON2ID$1$65536$4$16$1848e997d2f9a0cdfa8975a505ebaf73",
		"ARGON2ID$2$65536$4$32$f7ace26bad493acb466cce51ca0511aed6e563e3f3e69e970b3f758c66668ead",
		"ARGON2ID$1$32768$4$32$1f49786b8203202b4d2406a698e610cf58a533743cd55415c5f06e7e408f53cd",
		"ARGON2ID$1$65536$2$32$929fd9baa95f70c31d340975c91343d4e189340729f5d41f3ce1616358dd36d6",
		"ARGON2ID$1$65536$4$64$323c6d2236ba89b8cafcad21275838331ad83e5bac21a494c009ea167663fd126a2001e28cbee5174811d765090800eefe0433c3ed63b2fe3774f474fc735afd",
	}

	for i, key := range validKeys {
		err = hasher.Verify([]byte(plaintext), salt, []byte(key))
		assert.NoError(err, fmt.Sprintf("Key number %d should be valid", i+1))
	}

	err = hasher.Verify([]byte("some other plaintext value"), salt, expectedHashBytes)
	assert.Equal(crypto.ErrHashMissmatch, err)

	// Wrong hash value
	wrongHash := []byte("ARGON2ID$1$65536$4$32$8bfa337b938e7e1c137d5af2484667469ec7aa2d932183f0df7c6b43e040eb3d")
	err = hasher.Verify([]byte(plaintext), salt, wrongHash)
	assert.Equal(crypto.ErrHashMissmatch, err)

	// Wrong time value
	wrongParams2 := []byte("ARGON2ID$2$65536$4$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4")
	err = hasher.Verify([]byte(plaintext), salt, wrongParams2)
	assert.Equal(crypto.ErrHashMissmatch, err)

	// Wrong memory value
	wrongParams3 := []byte("ARGON2ID$1$32768$4$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4")
	err = hasher.Verify([]byte(plaintext), salt, wrongParams3)
	assert.Equal(crypto.ErrHashMissmatch, err)

	// Wrong thread count
	wrongParams4 := []byte("ARGON2ID$1$65536$2$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4")
	err = hasher.Verify([]byte(plaintext), salt, wrongParams4)
	assert.Equal(crypto.ErrHashMissmatch, err)

	err = hasher.Verify([]byte(plaintext), []byte("some other salt value"), expectedHashBytes)
	assert.Equal(crypto.ErrHashMissmatch, err)
}

func TestParseArgon2Key(t *testing.T) {
	assert := assert.New(t)

	var hasher crypto.Hasher = crypto.DefaultScryptHasher()

	plaintext := "625181dbfb5c6100cdacd97f3ba32ab4"
	salt, err := hex.DecodeString("478c1d403dec20707cf487f81c06d646")
	assert.NoError(err)

	invalidKeys := []string{
		"ARGON2ID$-1$65536$4$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4",
		"SCRYPT$1$65536$4$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4",
		"ARGON2ID$1$65536$4$32$1848e997d2f9a0cdfa8975a505ebaf73",
		"ARGON2ID$1$not-an-uint32$4$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4",
		"ARGON2ID$1$65536$131072$32$fe8c4ec2941d08c0588db1daa1030762065838357ff6e400bbaa00edeb1ff6a4",
	}

	for i, key := range invalidKeys {
		err := hasher.Verify([]byte(plaintext), salt, []byte(key))
		assert.True(errors.Is(err, crypto.ErrInvalidKey), fmt.Sprintf("Key number %d should be invalid", i+1))
	}
}
