package main

import (
	"os"

	"github.com/infocus7/imp/cmd"
	"github.com/pterm/pterm"
)

func main() {
	// Initialize PTerm
	pterm.EnableDebugMessages()

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
