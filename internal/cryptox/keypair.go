package cryptox

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Keypair struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewKeypair(
	privateKey ed25519.PrivateKey,
	publicKey ed25519.PublicKey,
) (*Keypair, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid Ed25519 private key")
	}

	if len(publicKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid Ed25519 public key")
	}

	return &Keypair{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := validatePath(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid PEM private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	edKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("PEM is not an Ed25519 private key")
	}

	return edKey, nil
}

func LoadPublicKey(path string) (ed25519.PublicKey, error) {
	data, err := validatePath(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid PEM public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	edKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("PEM is not an Ed25519 public key")
	}

	return edKey, nil
}

func (k *Keypair) PrivateKey() ed25519.PrivateKey {
	return k.privateKey
}

func (k *Keypair) PublicKey() ed25519.PublicKey {
	return k.publicKey
}

func validatePath(path string) ([]byte, error) {
	clean := filepath.Clean(path)
	if strings.Contains(clean, "..") {
		return nil, errors.New("invalid path: path traversal detected")
	}

	abs, err := filepath.Abs(clean)
	if err != nil {
		return nil, err
	}

	// #nosec G304 -- path is validated and provided via trusted configuration
	return os.ReadFile(abs)
}
