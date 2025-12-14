package cryptox

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"
)

const (
	upperAlphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes         = 32
	minNumberLen  = 4
)

// RandChars returns an n-length string of random uppercase letters and digits
// using crypto/rand.
func RandChars(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}

	out := make([]byte, n)
	maxLen := big.NewInt(int64(len(upperAlphaNum)))

	for i := range n {
		idx, err := rand.Int(rand.Reader, maxLen)
		if err != nil {
			return "", err
		}
		out[i] = upperAlphaNum[idx.Int64()]
	}

	return string(out), nil
}

// RandToken returns a URL-safe base64 token from 32 crypto-random bytes.
// Panics on CSPRNG failure.
func RandToken() string {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

// RandBytes returns n-length crypto-random bytes.
func RandBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Sha256 returns an SHA-256 from input.
func Sha256(input string) string {
	b := sha256.Sum256([]byte(input))

	return hex.EncodeToString(b[:])
}

// MaskNumber masks sensitive numbers showing only first 2 and last 2 digits.
// Returns format XX****XX where * represents masked digits.
func MaskNumber(n string) string {
	n = strings.ReplaceAll(n, " ", "")
	n = strings.ReplaceAll(n, "-", "")

	length := len(n)

	if length <= minNumberLen {
		return strings.Repeat("*", length)
	}

	first := n[:2]
	last := n[length-2:]

	middle := strings.Repeat("*", minNumberLen)

	return first + middle + last
}
