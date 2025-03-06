package cmd

import (
	"fmt"

	"github.com/bxtal-lsn/gotransport/pkg/cert"
	"github.com/spf13/cobra"
)

var (
	caKey  string
	caCert string
)

func init() {
	// Create command
	caCmd := &cobra.Command{
		Use:   "ca",
		Short: "Create CA certificate",
		Long:  `Create a Certificate Authority (CA) certificate and private key`,
		RunE:  runCACreate,
	}

	// Add flags
	caCmd.Flags().StringVarP(&caKey, "key-out", "k", "ca.key", "destination path for CA key")
	caCmd.Flags().StringVarP(&caCert, "cert-out", "o", "ca.crt", "destination path for CA certificate")

	// Add to root command
	rootCmd.AddCommand(caCmd)
}

func runCACreate(cmd *cobra.Command, args []string) error {
	if config.CACert == nil {
		return fmt.Errorf("no CA certificate configuration found in config file")
	}

	if verbose {
		fmt.Printf("Creating CA certificate...\n")
		fmt.Printf("CA Subject: %+v\n", config.CACert.Subject)
		fmt.Printf("Valid for: %d years\n", config.CACert.ValidForYears)
	}

	err := cert.CreateCACert(config.CACert, caKey, caCert)
	if err != nil {
		return fmt.Errorf("create CA error: %w", err)
	}

	fmt.Printf("CA created successfully!\n")
	fmt.Printf("Key: %s\n", caKey)
	fmt.Printf("Certificate: %s\n", caCert)
	return nil
}
