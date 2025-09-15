package config

import (
	"log"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/pgvillage-tools/chainsmith/pkg/tls"
)

// Config is the root object that will be loaded from a config yaml file
type Config struct {
	Chain              *tls.Chain               `json:"chain"`
	RootCAPath         string                   `json:"root_ca_path"`
	RootExpiry         time.Duration            `json:"root_expiry"`
	RootKeyUsages      tls.KeyUsages            `json:"root_key_usages"`
	RootExtraKeyUsages tls.ExtKeyUsages         `json:"root_extra_key_usages"`
	RootSubject        tls.Subject              `json:"subject"`
	IntermediateCAPath string                   `json:"intermediate_ca_path"`
	Intermediates      tls.ClassicIntermediates `json:"intermediates"`
	TmpDir             string                   `json:"tmpdir"`
}

// LoadConfig is used to load a yaml file and return a Config object
func LoadConfig(configPath string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	log.Printf("Loaded configuration from %s", configPath)
	return &cfg, nil
}

// AsChain can be used to derive all chain config from all config in the config
// file
func (c Config) AsChain() (*tls.Chain, error) {
	if c.Chain != nil {
		return c.Chain, nil
	}
	keyUsages, err := c.RootKeyUsages.AsKeyUsage()
	if err != nil {
		return nil, err
	}
	extraKeyUsages, err := c.RootExtraKeyUsages.AsEKeyUsages()
	if err != nil {
		return nil, err
	}

	c.Chain = &tls.Chain{
		Root: tls.Pair{
			Cert: tls.Cert{
				Subject:     &c.RootSubject,
				Expiry:      c.RootExpiry,
				KeyUsage:    keyUsages,
				ExtKeyUsage: extraKeyUsages,
			},
		},
		Intermediates: c.Intermediates.AsIntermediates(),
		Store:         c.TmpDir,
	}
	return c.Chain, nil
}
