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
		fmt.Println("Issued Certificates:")
		for name, certCfg := range cfg.Certificates {
			fmt.Printf("- %s (%s)\n", name, certCfg.CommonName)
		}
		return nil
	},
}
