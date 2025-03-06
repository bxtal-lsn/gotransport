package cmd

import (
	"fmt"
	"os"

	"github.com/bxtal-lsn/gotransport/pkg/cert"
	"github.com/spf13/cobra"
)

var (
	certKeyPath string
	certPath    string
	certName    string
)

func init() {
	// Create command
	certCmd := &cobra.Command{
		Use:   "cert",
		Short: "Create certificate",
		Long:  `Create a certificate signed by your CA`,
		RunE:  runCertCreate,
	}

	// Add flags
	certCmd.Flags().StringVarP(&certKeyPath, "key-out", "k", "server.key", "destination path for certificate key")
	certCmd.Flags().StringVarP(&certPath, "cert-out", "o", "server.crt", "destination path for certificate")
	certCmd.Flags().StringVarP(&certName, "name", "n", "", "name of the certificate in the config file")
	certCmd.Flags().StringVar(&caKey, "ca-key", "ca.key", "CA key path to sign certificate")
	certCmd.Flags().StringVar(&caCert, "ca-cert", "ca.crt", "CA cert path for certificate")

	// Mark required flags
	certCmd.MarkFlagRequired("name")
	certCmd.MarkFlagRequired("ca-key")
	certCmd.MarkFlagRequired("ca-cert")

	// Add to root command
	rootCmd.AddCommand(certCmd)
}

func runCertCreate(cmd *cobra.Command, args []string) error {
	// Read CA key
	caKeyBytes, err := os.ReadFile(caKey)
	if err != nil {
		return fmt.Errorf("CA key read error: %w", err)
	}

	// Read CA cert
	caCertBytes, err := os.ReadFile(caCert)
	if err != nil {
		return fmt.Errorf("CA cert read error: %w", err)
	}

	// Check if certificate exists in config
	certConfig, ok := config.Cert[certName]
	if !ok {
		return fmt.Errorf("certificate '%s' not found in configuration", certName)
	}

	if verbose {
		fmt.Printf("Creating certificate '%s'...\n", certName)
		fmt.Printf("Subject: %+v\n", certConfig.Subject)
		fmt.Printf("DNS Names: %v\n", certConfig.DNSNames)
		fmt.Printf("Valid for: %d years\n", certConfig.ValidForYears)
	}

	// Create certificate
	err = cert.CreateCert(certConfig, caKeyBytes, caCertBytes, certKeyPath, certPath)
	if err != nil {
		return fmt.Errorf("create certificate error: %w", err)
	}

	fmt.Printf("Certificate '%s' created successfully!\n", certName)
	fmt.Printf("Key: %s\n", certKeyPath)
	fmt.Printf("Certificate: %s\n", certPath)
	return nil
}
