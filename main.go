package main

import (
	"log"
	"chainsmith/config"
	"chainsmith/tls"
)

func run(configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	rootCert, rootKey, err := tls.GenerateCA(cfg.RootCAPath, cfg.RootCAPath+".key", nil, nil, true)
	if err != nil {
		return err
	}

	intermediateCert, intermediateKey, err := tls.GenerateCA(cfg.IntermediateCAPath, cfg.IntermediateCAPath+".key", rootCert, rootKey, false)
	if err != nil {
		return err
	}

	for name, certCfg := range cfg.Certificates {
		log.Printf("Generating certificate for %s...", name)
		if err := tls.GenerateCert(certCfg.CertPath, certCfg.KeyPath, intermediateCert, intermediateKey, certCfg.CommonName); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run("configs/config.yml"); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
