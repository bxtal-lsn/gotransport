package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner" // Add this
	"github.com/bxtal-lsn/gotransport/pkg/cert"
	"github.com/fatih/color" // Add this
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

// Update this function
func runCACreate(cmd *cobra.Command, args []string) error {
	if config.CACert == nil {
		return fmt.Errorf("no CA certificate configuration found in config file")
	}

	if verbose {
		printInfo("Creating CA certificate...")
		fmt.Printf("CA Subject: %+v\n", config.CACert.Subject)
		fmt.Printf("Valid for: %d years\n", config.CACert.ValidForYears)
	}

	// Create spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Creating CA certificate..."
	s.Color("cyan")
	s.Start()

	// Perform operation
	err := cert.CreateCACert(config.CACert, caKey, caCert)

	// Stop spinner
	s.Stop()

	if err != nil {
		printError("Failed to create CA: %v", err)
		return fmt.Errorf("create CA error: %w", err)
	}

	printSuccess("CA created successfully!")
	color.New(color.FgHiWhite).Printf("Key: %s\n", caKey)
	color.New(color.FgHiWhite).Printf("Certificate: %s\n", caCert)
	return nil
}

