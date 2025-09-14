package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var revokeCmd = &cobra.Command{
	Use:   "revoke <certificate_name>",
	Short: "Revoke a certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		certName := args[0]
		cfg, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}

		certCfg, exists := cfg.Certificates[certName]
		if !exists {
			return fmt.Errorf("certificate '%s' not found", certName)
		}

		if err := os.Remove(certCfg.CertPath); err != nil {
			return fmt.Errorf("failed to delete certificate file: %v", err)
		}
		if err := os.Remove(certCfg.KeyPath); err != nil {
			return fmt.Errorf("failed to delete key file: %v", err)
		}

		fmt.Printf("Certificate '%s' revoked successfully.\n", certName)
		return nil
	},
}
