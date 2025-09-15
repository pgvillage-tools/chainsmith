package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all issued certificates",
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}
		chain, err := cfg.AsChain()
		if err != nil {
			return err
		}
		structure := chain.Structure()
		fmt.Println("Issued Certificates:")
		for iName, certs := range structure.Certs {
			fmt.Printf("- intermediate: %s\n", iName)
			for name := range certs {
				fmt.Printf("  - %s\n", name)
			}
		}
		return nil
	},
}
