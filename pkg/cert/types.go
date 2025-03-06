package cert

import "math/big"

// CACert represents a Certificate Authority configuration
type CACert struct {
	Serial        *big.Int    `yaml:"serial"`
	ValidForYears int         `yaml:"validForYears"`
	Subject       CertSubject `yaml:"subject"`
}

// Cert represents a certificate configuration
type Cert struct {
	Serial        *big.Int    `yaml:"serial"`
	ValidForYears int         `yaml:"validForYears"`
	Subject       CertSubject `yaml:"subject"`
	DNSNames      []string    `yaml:"dnsNames"`
}

// CertSubject represents the subject fields of a certificate
type CertSubject struct {
	Country            string `yaml:"country"`
	Organization       string `yaml:"organization"`
	OrganizationalUnit string `yaml:"organizationalUnit"`
	Locality           string `yaml:"locality"`
	Province           string `yaml:"province"`
	StreetAddress      string `yaml:"streetAddress"`
	PostalCode         string `yaml:"postalCode"`
	SerialNumber       string `yaml:"serialNumber"`
	CommonName         string `yaml:"commonName"`
}
