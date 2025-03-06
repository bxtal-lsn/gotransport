package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// Config holds TLS configuration options
type Config struct {
	CAPath     string
	CertPath   string
	KeyPath    string
	ServerName string
	ClientAuth tls.ClientAuthType
	MinVersion uint16
}

// DefaultConfig returns a default secure TLS configuration
func DefaultConfig() Config {
	return Config{
		ClientAuth: tls.NoClientCert,
		MinVersion: tls.VersionTLS12,
	}
}

// NewServerTLSConfig creates a server TLS configuration
func NewServerTLSConfig(cfg Config) (*tls.Config, error) {
	// Set default min version if not specified
	if cfg.MinVersion == 0 {
		cfg.MinVersion = tls.VersionTLS12
	}

	// Create base TLS config
	tlsConfig := &tls.Config{
		MinVersion: cfg.MinVersion,
		ClientAuth: cfg.ClientAuth,
	}

	// Load server certificate and key if provided
	if cfg.CertPath != "" && cfg.KeyPath != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load server certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if provided
	if cfg.CAPath != "" && cfg.ClientAuth != tls.NoClientCert {
		caPool, err := loadCACert(cfg.CAPath)
		if err != nil {
			return nil, err
		}
		tlsConfig.ClientCAs = caPool
	}

	return tlsConfig, nil
}

// NewClientTLSConfig creates a client TLS configuration
func NewClientTLSConfig(cfg Config) (*tls.Config, error) {
	// Set default min version if not specified
	if cfg.MinVersion == 0 {
		cfg.MinVersion = tls.VersionTLS12
	}

	// Create base TLS config
	tlsConfig := &tls.Config{
		MinVersion:         cfg.MinVersion,
		InsecureSkipVerify: false, // Always verify server certs by default
	}

	// Set server name if provided
	if cfg.ServerName != "" {
		tlsConfig.ServerName = cfg.ServerName
	}

	// Load client certificate and key if provided
	if cfg.CertPath != "" && cfg.KeyPath != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if provided
	if cfg.CAPath != "" {
		caPool, err := loadCACert(cfg.CAPath)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = caPool
	}

	return tlsConfig, nil
}

// loadCACert loads a CA certificate from a file into a certificate pool
func loadCACert(caPath string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	return caPool, nil
}
