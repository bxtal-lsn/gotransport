package cmd

import (
	"fmt"
	"os"
	"strings"

	configpkg "github.com/bxtal-lsn/gotransport/internal/config"
	"github.com/bxtal-lsn/gotransport/pkg/cert"
	"github.com/fatih/color" // Add this import
	"github.com/spf13/cobra"
)

// Add this logo constant
const logo = `
  _____      _____                                     _   
 / ____|    |_   _|                                   | |  
| |  __  ___  | |_ __ __ _ _ __  ___ _ __   ___  _ __| |_ 
| | |_ |/ _ \ | | '__/ _' | '_ \/ __| '_ \ / _ \| '__| __|
| |__| | (_) || | | | (_| | | | \__ \ |_) | (_) | |  | |_ 
 \_____|\___/___|_|  \__,_|_| |_|___/ .__/ \___/|_|   \__|
                                    | |                    
                                    |_|                    
`

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
	// Only print the logo if we're running the main command (not subcommands)
	if len(os.Args) <= 1 || (len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-")) {
		color.Cyan(logo)
		fmt.Println("Modern TLS Certificate & DNS Management Tool")
		fmt.Println()
	}

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

// Add these helper functions for help styling - we'll use them later
func styleHeading(heading string) string {
	return color.New(color.FgHiCyan, color.Bold).Sprint(heading)
}

func styleExample(example string) string {
	lines := strings.Split(example, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "#") {
				lines[i] = "  " + color.New(color.FgHiBlack).Sprint(trimmed)
			} else {
				lines[i] = "  " + color.New(color.FgCyan).Sprint("$ ") + trimmed
			}
		}
	}
	return strings.Join(lines, "\n")
}
