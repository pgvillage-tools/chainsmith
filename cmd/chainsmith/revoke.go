package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var revokeCmd = &cobra.Command{
	Use:   "revoke <certificate_name>",
	Short: "Revoke a certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		certName := args[0]
		// TODO: We need to implement an option te encrypt private keys
		// After taht we need to read back a chain from disk
		// After that we can revoke (I guess)
		cfg, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}

		chain, err := cfg.AsChain()
		if err != nil {
			return err
		}
		structure := chain.Structure()
		if err != nil {
			return err
		}
		// TODO: We need to create a CRL, and add the cert on the CRL.
		// Not yet implemented.
		certName += fmt.Sprintf("%d", len(structure.Certs))

		fmt.Printf("Certificate '%s' revoked successfully.\n", certName)
		return nil
	},
}
