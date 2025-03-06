package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// CreateRSAPrivateKey generates a new RSA private key with the specified bit size
func CreateRSAPrivateKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, fmt.Errorf("key length must be at least 2048 bits for security reasons")
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// RSAPrivateKeyToPEM converts an RSA private key to a PEM block
func RSAPrivateKeyToPEM(privateKey *rsa.PrivateKey) *pem.Block {
	if privateKey == nil {
		return nil
	}
	return &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
}

// CreateRSAPrivateKeyAndSave generates a new RSA private key and saves it to a file
func CreateRSAPrivateKeyAndSave(path string, bits int) error {
	// Generate key
	privateKey, err := CreateRSAPrivateKey(bits)
	if err != nil {
		return err
	}

	// Create or open file with secure permissions
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open key file: %w", err)
	}
	defer f.Close()

	// Encode and write key
	block := RSAPrivateKeyToPEM(privateKey)
	if err := pem.Encode(f, block); err != nil {
		return fmt.Errorf("failed to encode key: %w", err)
	}

	return nil
}

// PrivateKeyPemToRSA converts a PEM encoded RSA private key to an rsa.PrivateKey
func PrivateKeyPemToRSA(input []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(input)
	if block == nil {
		return nil, fmt.Errorf("failed to parse key PEM")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("PEM block is not an RSA private key (type: %s)", block.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA private key: %w", err)
	}

	return privateKey, nil
}
