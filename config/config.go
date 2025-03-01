package config

import (
	"os"
	"testing"
	"gopkg.in/yaml.v3"
)

type CertificateConfig struct {
	CertPath   string `yaml:"cert_path"`
	KeyPath    string `yaml:"key_path"`
	CommonName string `yaml:"common_name"`
}

type Config struct {
	RootCAPath         string                        `yaml:"root_ca_path"`
	IntermediateCAPath string                        `yaml:"intermediate_ca_path"`
	Certificates       map[string]CertificateConfig `yaml:"certificates"`
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

// Unit Test for Config
func TestLoadConfig(t *testing.T) {
	tempConfig := "test_config.yml"
	configContent := `
root_ca_path: "test_rootCA.crt"
intermediate_ca_path: "test_intermediateCA.crt"
certificates:
  server:
    cert_path: "test_server.crt"
    key_path: "test_server.key"
    common_name: "server.local"
  client:
    cert_path: "test_client.crt"
    key_path: "test_client.key"
    common_name: "client.local"
`

	err := os.WriteFile(tempConfig, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(tempConfig)

	cfg, err := LoadConfig(tempConfig)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.RootCAPath != "test_rootCA.crt" {
		t.Errorf("Expected root_ca_path to be 'test_rootCA.crt', got %s", cfg.RootCAPath)
	}

	if _, exists := cfg.Certificates["server"]; !exists {
		t.Error("Expected server certificate configuration, but not found")
	}
}
