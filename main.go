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

	_, _, err = tls.GenerateCA(cfg.IntCAPath, cfg.IntCAPath+".key", rootCert, rootKey, false)
	if err != nil {
		log.Fatalf("Error generating Intermediate CA: %v", err)
	}
}
