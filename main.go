package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	CertFilePath = "certs/rootCA.crt"
	KeyFilePath  = "certs/rootCA.key"
)

// GenerateRootCA creates a self-signed root certificate and private key
func GenerateRootCA() error {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	now := time.Now()
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"NL"},
			Organization:       []string{"Mannem Solutions"},
			OrganizationalUnit: []string{"Chainsmith TLS chain maker"},
			CommonName:        "chainsmith",
		},
		NotBefore:             now,
		NotAfter:              now.Add(3650 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	if err := ioutil.WriteFile(CertFilePath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write cert file: %v", err)
	}

	if err := ioutil.WriteFile(KeyFilePath, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %v", err)
	}

	log.Println("Root CA certificate and key generated successfully")
	return nil
}

func main() {
	if err := os.MkdirAll("certs", 0755); err != nil {
		log.Fatalf("Failed to create certs directory: %v", err)
	}

	if err := GenerateRootCA(); err != nil {
		log.Fatalf("Error generating Root CA: %v", err)
	}
}
