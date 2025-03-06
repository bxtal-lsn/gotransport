package main

import (
	"fmt"
	"os"

	"github.com/bxtal-lsn/gotransport/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
