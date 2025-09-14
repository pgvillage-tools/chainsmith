// Package config can be used to parse a `yaml` file into a config struct to be
// used in the rest of the application
package config

// CertificateConfig can be used to set some defaults regarding certificates
type CertificateConfig struct {
	CertPath   string `json:"cert_path"`
	KeyPath    string `json:"key_path"`
	CommonName string `json:"common_name"`
}
