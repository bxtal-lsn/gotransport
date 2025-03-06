package cmd

import (
	"fmt"
	"os"

	configpkg "github.com/bxtal-lsn/gotransport/internal/config"
	"github.com/bxtal-lsn/gotransport/pkg/cert"
	"github.com/spf13/cobra"
)

type Config struct {
	CACert *cert.CACert          `yaml:"caCert"`
	Cert   map[string]*cert.Cert `yaml:"certs"`
}

var (
	cfgFilePath string
	config      Config
	verbose     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotransport",
	Short: "A modern TLS certificate management tool",
	Long: `gotransport is a command line tool for managing TLS certificates.
It provides functionality to create and manage CA certificates, 
server/client certificates, and RSA keys with a clean CLI interface.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", "tls.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

func initConfig() {
	if cfgFilePath == "" {
		cfgFilePath = "tls.yaml"
	}

	var err error
	config, err = configpkg.LoadConfig[Config](cfgFilePath)
	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		}
		// We don't exit here because not all commands require config
	}
}
