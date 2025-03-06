package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Version information - these would typically be set during build
	Version    = "1.0.0"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display detailed version information about GoTransport.`,
		Run:   runVersionCmd,
	}

	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	// Format build date if it's a timestamp
	formattedDate := BuildDate
	if timestamp, err := time.Parse(time.RFC3339, BuildDate); err == nil {
		formattedDate = timestamp.Format("Jan 02, 2006 15:04:05 MST")
	}

	// Create color formatters
	titleStyle := color.New(color.FgHiCyan, color.Bold)
	valueStyle := color.New(color.FgHiWhite)

	// Print version info in a nice format
	fmt.Println()
	titleStyle.Println("GoTransport")
	fmt.Println("A modern TLS certificate & DNS management tool")
	fmt.Println()

	fmt.Printf("%-15s", "Version:")
	valueStyle.Printf("%s\n", Version)

	fmt.Printf("%-15s", "Git Commit:")
	valueStyle.Printf("%s\n", CommitHash)

	fmt.Printf("%-15s", "Built:")
	valueStyle.Printf("%s\n", formattedDate)

	fmt.Printf("%-15s", "Go Version:")
	valueStyle.Printf("%s\n", runtime.Version())

	fmt.Printf("%-15s", "OS/Arch:")
	valueStyle.Printf("%s/%s\n", runtime.GOOS, runtime.GOARCH)

	fmt.Println()
}
