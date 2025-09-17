// Package main is the m,ain entrypoint for this commandline utility
package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/goccy/go-yaml"
	"github.com/pgvillage-tools/chainsmith/internal/config"
)

func init() {
	rootCmd.PersistentFlags().String("config", os.Getenv("CHAINSMITH_CONFIG"), "Path to the config file")
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		panic(fmt.Errorf("init failed: %w", err).Error())
	}
	rootCmd.AddCommand(issueCmd, listCmd, revokeCmd)
}

func loadConfig(configPath string) (*config.Config, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
