package main

import (
	"fmt"
	"log"
	"path"

	"github.com/pgvillage-tools/chainsmith/internal/config"
	"github.com/pgvillage-tools/chainsmith/pkg/tls"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Generate CA and certificates based on the configuration file",
	RunE: func(_ *cobra.Command, _ []string) error {
		config, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}
		certs, err := issue(*config)
		if err != nil {
			return err
		}
		out, err := yaml.Marshal(certs)
		if err != nil {
			return err
		}
		_, err = fmt.Print(out)
		if err != nil {
			return err
		}
		return nil
	},
}

type certs struct {
	Certs intBodies `json:"certs"`
	Keys  intBodies `json:"private_keys"`
}
type intBodies map[string]bodies
type bodies map[string]string

func issue(cfg config.Config) (*certs, error) {
	// certs/private_keys [intermediate_name] [common_name] cert
	certBodies := intBodies{}
	keyBodies := intBodies{}
	caCertPath, caKeyPath := cfg.GetCaPaths()
	keyUsages, err := tls.DefaultKeyUsages.AsKeyUsage()
	extKeyUsages, err := tls.DefaultExtendedKeyUsages.AsEKeyUsages()
	rootCert, rootKey, err := tls.GenerateCA(
		caCertPath,
		caKeyPath,
		nil,
		nil,
		cfg.Subject.AsPkixName(),
		cfg.RootExpiry,
		keyUsages,
		extKeyUsages,
	)
	if err != nil {
		return nil, err
	}
	if len(cfg.Intermediates) > 0 {
		var intermediateCaCert, intermediateCaKey string
		for _, intermediate := range cfg.Intermediates {
			if cfg.TmpDir != "" {
				intermediateCaPath := path.Join(cfg.TmpDir, "tls", "int_"+intermediate.Name)
				intermediateCaCert = path.Join(intermediateCaPath, "cacert.pem")
				intermediateCaKey = path.Join(intermediateCaPath, "cakey.pem")
			}
			ku, err := intermediate.KeyUsages.AsKeyUsage()
			if err != nil {
				return nil, err
			}
			eku, err := intermediate.ExtendedKeyUsages.AsEKeyUsages()
			if err != nil {
				return nil, err
			}
			intermediateCert, intermediateKey, err := tls.GenerateCA(
				intermediateCaCert,
				intermediateCaKey,
				rootCert,
				rootKey,
				cfg.Subject.SetCommonName(intermediate.Name).AsPkixName(),
				intermediate.Expiry,
				ku,
				eku,
			)
			if err != nil {
				return nil, err
			}
			certBodies[intermediate.Name] = intermediateCert.
		}
	} else {
		intermediateCertPath, intermediateKeyPath := cfg.GetIntermediatePaths()
		intermediateCert, intermediateKey, err := tls.GenerateCA(
			intermediateCertPath,
			intermediateKeyPath,
			rootCert,
			rootKey,
			false,
		)
		if err != nil {
			return nil, err
		}
		for name, certCfg := range cfg.Certificates {
			log.Printf("Generating certificate for %s...", name)
			if err := tls.GenerateCert(certCfg.CertPath, certCfg.KeyPath, intermediateCert, intermediateKey, certCfg.CommonName); err != nil {
				return nil, err
			}
		}
	}

	out := &certs{
		Keys:  keyBodies,
		Certs: certBodies,
	}
	return out, nil
}
