// main.go
package main

import (
	"log"
	"chainsmith/config"
	"chainsmith/tls"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	rootCert, rootKey, err := tls.GenerateCA(cfg.RootCAPath, cfg.RootCAPath+".key", nil, nil, true)
	if err != nil {
		log.Fatalf("Error generating Root CA: %v", err)
	}

	intermediateCert, intermediateKey, err := tls.GenerateCA(cfg.IntermediateCAPath, cfg.IntermediateCAPath+".key", rootCert, rootKey, false)
	if err != nil {
		log.Fatalf("Error generating Intermediate CA: %v", err)
	}

	for name, certCfg := range cfg.Certificates {
		log.Printf("Generating certificate for %s...", name)
		if err := tls.GenerateCert(certCfg.CertPath, certCfg.KeyPath, intermediateCert, intermediateKey, certCfg.CommonName); err != nil {
			log.Fatalf("Error generating %s certificate: %v", name, err)
		}
	}
}
