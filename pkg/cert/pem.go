package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// PemToX509 converts PEM encoded certificate bytes to an x509.Certificate
func PemToX509(input []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(input)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("PEM block is not a certificate (type: %s)", block.Type)
	}

	return x509.ParseCertificate(block.Bytes)
}

// X509ToPem converts an x509.Certificate to PEM encoded bytes
func X509ToPem(cert *x509.Certificate) ([]byte, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate is nil")
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	return pem.EncodeToMemory(block), nil
}
