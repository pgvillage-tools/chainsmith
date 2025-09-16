package main

import (
	"fmt"

	"github.com/pgvillage-tools/chainsmith/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chainsmith",
	Short: "Chainsmith - A simple certificate chain manager",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Use --help to see available commands.")
	},
	Version: version.GetAppVersion(),
}
