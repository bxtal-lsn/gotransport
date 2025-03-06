package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2" // Add this
	"github.com/briandowns/spinner"    // Add this
	"github.com/bxtal-lsn/gotransport/pkg/key"
	"github.com/fatih/color" // Add this
	"github.com/spf13/cobra"
)

var (
	keyOut    string
	keyLength int
)

func init() {
	// Keep your existing init function unchanged
	// ...
}

// Update this function
func runKeyCreate(cmd *cobra.Command, args []string) error {
	// Check if file already exists and confirm overwrite
	if _, err := os.Stat(keyOut); err == nil {
		printWarning("Key file %s already exists", keyOut)
		var confirm bool
		prompt := &survey.Confirm{
			Message: "Overwrite existing file?",
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}
		if !confirm {
			printInfo("Operation cancelled")
			return nil
		}
	}

	if verbose {
		printInfo("Creating RSA private key with %d bits...", keyLength)
	}

	// Show spinner while creating key
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Creating %d-bit RSA key...", keyLength)
	s.Color("cyan")
	s.Start()

	// Perform operation
	err := key.CreateRSAPrivateKeyAndSave(keyOut, keyLength)
	s.Stop()

	if err != nil {
		printError("Failed to create key: %v", err)
		return fmt.Errorf("create key error: %w", err)
	}

	printSuccess("RSA key created successfully!")
	info := color.New(color.FgHiWhite)
	info.Printf("  Path: %s\n", keyOut)
	info.Printf("  Length: %d bits\n", keyLength)
	return nil
}

