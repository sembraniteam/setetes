package argon2x

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/sembraniteam/setetes/internal/config"
	"github.com/sembraniteam/setetes/internal/cryptox"
	"golang.org/x/crypto/argon2"
)

const (
	partLen        = 5
	minPepperLen   = 32
	minMemory      = 1024 * 12
	minIterations  = 1
	minParallelism = 1
	minSaltLen     = 16
	minKeyLen      = 16
)

type (
	Raw struct {
		config Config
		salt   []byte
		hash   []byte
	}

	Config struct {
		pepper      string
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}

	hashParams struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		salt        []byte
		hash        []byte
	}
)

func New(config config.Config) Config {
	p := config.Password
	c := Config{
		pepper:      p.Pepper,
		memory:      p.Argon2.Memory,
		iterations:  p.Argon2.Iterations,
		parallelism: p.Argon2.Parallelism,
		saltLength:  p.Argon2.SaltLength,
		keyLength:   p.Argon2.KeyLength,
	}

	if err := c.validate(); err != nil {
		panic(err)
	}

	return c
}

func Default() Config {
	p := config.Get().Password
	c := Config{
		pepper:      p.Pepper,
		memory:      p.Argon2.Memory,
		iterations:  p.Argon2.Iterations,
		parallelism: p.Argon2.Parallelism,
		saltLength:  p.Argon2.SaltLength,
		keyLength:   p.Argon2.KeyLength,
	}

	if err := c.validate(); err != nil {
		panic(err)
	}

	return c
}

func (c *Config) validate() error {
	if len(c.pepper) < minPepperLen {
		return errors.New("pepper too short, minimum 32 length")
	}
	if c.memory < minMemory {
		return errors.New("memory too low, minimum 8MB")
	}
	if c.iterations < minIterations {
		return errors.New("iterations must be at least 1")
	}
	if c.parallelism < minParallelism {
		return errors.New("parallelism must be at least 1")
	}
	if c.saltLength < minSaltLen {
		return errors.New("salt too short, minimum 16 bytes")
	}
	if c.keyLength < minKeyLen {
		return errors.New("key too short, minimum 16 bytes")
	}

	return nil
}

func (c *Config) Hash(text, salt []byte) (*Raw, error) {
	if len(text) == 0 {
		return nil, errors.New("invalid password")
	}

	pepperedText := append(text, []byte(c.pepper)...)

	if salt == nil {
		var err error
		salt, err = cryptox.RandBytes(c.saltLength)
		if err != nil {
			return nil, err
		}
	}

	hash := argon2.IDKey(
		pepperedText,
		salt,
		c.iterations,
		c.memory,
		c.parallelism,
		c.keyLength,
	)

	return &Raw{
		config: *c,
		salt:   salt,
		hash:   hash,
	}, nil
}

func (c *Config) HashString(text []byte) (string, error) {
	raw, err := c.Hash(text, nil)
	if err != nil {
		return "", err
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(raw.salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(raw.hash)

	hash := fmt.Sprintf(
		"argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		c.memory,
		c.iterations,
		c.parallelism,
		b64Salt,
		b64Hash,
	)

	return hash, nil
}

func (r *Raw) Verify(text []byte) (bool, error) {
	raw, err := r.config.Hash(text, r.salt)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(r.hash, raw.hash) == 1, nil
}

func (c *Config) VerifyString(text []byte, hashString string) (bool, error) {
	params, err := parseHashString(hashString)
	if err != nil {
		return false, err
	}

	hashLen := len(params.hash)
	if hashLen > math.MaxUint32 {
		return false, errors.New("hash length overflow")
	}

	raw := &Raw{
		config: Config{
			pepper:      c.pepper,
			memory:      params.memory,
			iterations:  params.iterations,
			parallelism: params.parallelism,
			saltLength:  c.saltLength,
			keyLength:   uint32(hashLen),
		},
		salt: params.salt,
		hash: params.hash,
	}

	return raw.Verify(text)
}

func parseHashString(hashString string) (*hashParams, error) {
	if hashString == "" {
		return nil, errors.New("password hash is empty")
	}

	parts := strings.Split(hashString, "$")
	if len(parts) != partLen {
		return nil, fmt.Errorf(
			"invalid password hash format: expected %d parts separated by '$', got %d",
			partLen,
			len(parts),
		)
	}

	if parts[0] != "argon2id" {
		return nil, fmt.Errorf(
			"invalid hash algorithm: expected 'argon2id', got '%s'",
			parts[0],
		)
	}

	var version int
	if _, err := fmt.Sscanf(parts[1], "v=%d", &version); err != nil {
		return nil, fmt.Errorf("invalid version format: %w", err)
	}

	if version != argon2.Version {
		return nil, fmt.Errorf(
			"incompatible Argon2 version: expected %d, got %d",
			argon2.Version,
			version,
		)
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(
		parts[2],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	); err != nil {
		return nil, fmt.Errorf("invalid parameters format: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, fmt.Errorf("invalid salt encoding: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, fmt.Errorf("invalid hash encoding: %w", err)
	}

	return &hashParams{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		salt:        salt,
		hash:        hash,
	}, nil
}
