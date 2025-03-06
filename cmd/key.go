package cmd

import (
	"fmt"

	"github.com/bxtal-lsn/gotransport/pkg/key"
	"github.com/spf13/cobra"
)

var (
	keyOut    string
	keyLength int
)

func init() {
	// Create command
	keyCmd := &cobra.Command{
		Use:   "key",
		Short: "Create RSA key",
		Long:  `Create a new RSA private key for use with certificates`,
		RunE:  runKeyCreate,
	}

	// Add flags
	keyCmd.Flags().StringVarP(&keyOut, "key-out", "k", "key.pem", "destination path for key")
	keyCmd.Flags().IntVarP(&keyLength, "key-length", "l", 4096, "key length in bits")

	// Add to root command
	rootCmd.AddCommand(keyCmd)
}

func runKeyCreate(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Printf("Creating RSA private key with %d bits...\n", keyLength)
	}

	err := key.CreateRSAPrivateKeyAndSave(keyOut, keyLength)
	if err != nil {
		return fmt.Errorf("create key error: %w", err)
	}

	fmt.Printf("RSA key created successfully!\n")
	fmt.Printf("Key: %s\n", keyOut)
	fmt.Printf("Length: %d bits\n", keyLength)
	return nil
}
