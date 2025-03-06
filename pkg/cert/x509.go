package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/bxtal-lsn/gotransport/pkg/key"
)

// CreateCACert creates a new Certificate Authority certificate and key
func CreateCACert(ca *CACert, keyFilePath, caCertFilePath string) error {
	// Create certificate template
	template := &x509.Certificate{
		SerialNumber: ca.Serial,
		Subject: pkix.Name{
			Country:            removeEmptyString([]string{ca.Subject.Country}),
			Organization:       removeEmptyString([]string{ca.Subject.Organization}),
			OrganizationalUnit: removeEmptyString([]string{ca.Subject.OrganizationalUnit}),
			Locality:           removeEmptyString([]string{ca.Subject.Locality}),
			Province:           removeEmptyString([]string{ca.Subject.Province}),
			StreetAddress:      removeEmptyString([]string{ca.Subject.StreetAddress}),
			PostalCode:         removeEmptyString([]string{ca.Subject.PostalCode}),
			CommonName:         ca.Subject.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(ca.ValidForYears, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create certificate and key
	keyBytes, certBytes, err := createCert(template, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	// Write key file
	if err := os.WriteFile(keyFilePath, keyBytes, 0o600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	// Write certificate file
	if err := os.WriteFile(caCertFilePath, certBytes, 0o644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	return nil
}

// CreateCert creates a new certificate signed by a CA
func CreateCert(cert *Cert, caKey []byte, caCert []byte, keyFilePath, certFilePath string) error {
	// Create certificate template
	template := &x509.Certificate{
		SerialNumber: cert.Serial,
		Subject: pkix.Name{
			Country:            removeEmptyString([]string{cert.Subject.Country}),
			Organization:       removeEmptyString([]string{cert.Subject.Organization}),
			OrganizationalUnit: removeEmptyString([]string{cert.Subject.OrganizationalUnit}),
			Locality:           removeEmptyString([]string{cert.Subject.Locality}),
			Province:           removeEmptyString([]string{cert.Subject.Province}),
			StreetAddress:      removeEmptyString([]string{cert.Subject.StreetAddress}),
			PostalCode:         removeEmptyString([]string{cert.Subject.PostalCode}),
			CommonName:         cert.Subject.CommonName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(cert.ValidForYears, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
		DNSNames:    removeEmptyString(cert.DNSNames),
	}

	// Parse CA key
	caKeyParsed, err := key.PrivateKeyPemToRSA(caKey)
	if err != nil {
		return fmt.Errorf("failed to parse CA key: %w", err)
	}

	// Parse CA certificate
	caCertParsed, err := PemToX509(caCert)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Create certificate and key
	keyBytes, certBytes, err := createCert(template, caKeyParsed, caCertParsed)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Write key file
	if err := os.WriteFile(keyFilePath, keyBytes, 0o600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	// Write certificate file
	if err := os.WriteFile(certFilePath, certBytes, 0o644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	return nil
}

// createCert is a helper function that creates a certificate and key pair
func createCert(template *x509.Certificate, caKey *rsa.PrivateKey, caCert *x509.Certificate) ([]byte, []byte, error) {
	var (
		derBytes []byte
		certOut  bytes.Buffer
		keyOut   bytes.Buffer
		err      error
	)

	// Create private key
	privateKey, err := key.CreateRSAPrivateKey(4096)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create private key: %w", err)
	}

	// Create certificate based on whether it's a CA or not
	if template.IsCA {
		// Self-signed certificate for CA
		derBytes, err = x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	} else {
		// Certificate signed by CA
		derBytes, err = x509.CreateCertificate(rand.Reader, template, caCert, &privateKey.PublicKey, caKey)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate to PEM
	if err = pem.Encode(&certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, fmt.Errorf("failed to encode certificate: %w", err)
	}

	// Encode key to PEM
	if err = pem.Encode(&keyOut, key.RSAPrivateKeyToPEM(privateKey)); err != nil {
		return nil, nil, fmt.Errorf("failed to encode key: %w", err)
	}

	return keyOut.Bytes(), certOut.Bytes(), nil
}

// removeEmptyString filters out empty strings from a slice
func removeEmptyString(input []string) []string {
	if len(input) == 0 {
		return []string{}
	}

	if len(input) == 1 && input[0] == "" {
		return []string{}
	}

	result := make([]string, 0, len(input))
	for _, s := range input {
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}
