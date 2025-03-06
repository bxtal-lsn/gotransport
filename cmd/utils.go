package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

// Helper functions to print messages
func printSuccess(format string, args ...interface{}) {
	fmt.Print("✓ ")
	color.Green(format, args...)
}

func printInfo(format string, args ...interface{}) {
	fmt.Print("ℹ ")
	color.Cyan(format, args...)
}

func printWarning(format string, args ...interface{}) {
	fmt.Print("⚠ ")
	color.Yellow(format, args...)
}

func printError(format string, args ...interface{}) {
	fmt.Print("✗ ")
	color.Red(format, args...)
}
