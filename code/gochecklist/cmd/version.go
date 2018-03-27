package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build number and versions injected at compile time
var (
	Build   string
	Version string
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build: %s\nVersion: %s\n", Build, Version)
	},
}
