package main

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pgvillage-tools/chainsmith/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Generate CA and certificates based on the configuration file",
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}
		out, err := issue(cfg)
		if err != nil {
			return err
		}
		_, err = fmt.Print(out)
		return err
	},
}

type certs struct {
	Certs intBodies `json:"certs"`
	Keys  intBodies `json:"private_keys"`
}
type intBodies map[string]bodies
type bodies map[string]string

func issue(cfg *config.Config) ([]byte, error) {
	chain := cfg.AsChain()
	if err := chain.InitializeCA(); err != nil {
		return nil, err
	}
	if err := chain.InitializeIntermediates(); err != nil {
		return nil, err
	}
	structure := chain.Structure()
	return yaml.Marshal(structure)
}
