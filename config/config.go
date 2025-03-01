package config

import (
	"os"
	"log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RootCAPath string `yaml:"root_ca_path"`
	IntCAPath  string `yaml:"intermediate_ca_path"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
