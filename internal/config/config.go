package config

import (
	"log"
	"path"
	"time"

	"github.com/pgvillage-tools/chainsmith/pkg/tls"
	"github.com/spf13/viper"
)

// Config is the root object that will be loaded from a config yaml file
type Config struct {
	Chain              *tls.Chain                   `json:"chain"`
	RootCAPath         string                       `json:"root_ca_path"`
	RootExpiry         time.Duration                `json:"root_expiry"`
	IntermediateCAPath string                       `json:"intermediate_ca_path"`
	Certificates       map[string]CertificateConfig `json:"certificates"`
	Intermediates      tls.ClassicIntermediates     `json:"intermediates"`
	Subject            tls.Subject                  `json:"subject"`
	TmpDir             string                       `json:"tmpdir"`
}

// GetCaPaths derives path for ca cert and key from config settings
func (c Config) GetCaPaths() (string, string) {
	if c.RootCAPath != "" {
		return c.RootCAPath, c.RootCAPath + ".key"
	}
	if c.TmpDir != "" {
		rootCaPath := path.Join(c.TmpDir, "tls", "certs")
		rootCaCert := path.Join(rootCaPath, "cacert.pem")
		rootCaKey := path.Join(rootCaPath, "cakey.pem")
		return rootCaCert, rootCaKey
	}
	return "", ""
}

// GetIntermediatePaths derives path for ca cert and key from config settings
func (c Config) GetIntermediatePaths() (string, string) {
	if c.IntermediateCAPath != "" {
		return c.IntermediateCAPath, c.IntermediateCAPath + ".key"
	}
	if c.TmpDir != "" {
		intermediateCaPath := path.Join(c.TmpDir, "tls", "int_server")
		intermediateCaCert := path.Join(intermediateCaPath, "cacert.pem")
		intermediateCaKey := path.Join(intermediateCaPath, "cakey.pem")
		return intermediateCaCert, intermediateCaKey
	}
	return "", ""
}

// LoadConfig is used to load a yaml file and return a Config object
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // Allow environment variable overrides

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	log.Printf("Loaded configuration from %s", viper.ConfigFileUsed())
	return &cfg, nil
}

// AsChain can be used to derive all chain config from all config in the config
// file
func (c Config) AsChain() *tls.Chain {
	if c.Chain != nil {
		return c.Chain
	}
	return &tls.Chain{
		Subject:       c.Subject,
		Intermediates: c.Intermediates.AsIntermediates(),
		Expiry:        c.RootExpiry,
		Store:         c.TmpDir,
	}
}
