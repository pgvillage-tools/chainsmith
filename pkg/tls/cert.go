// Package tls takes care of all tls actions for a chain
package tls

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path"
	"time"
)

// Certs is a collection of Cert objects
type Certs []Cert

// Cert is an object representing a certificate
type Cert struct {
	cert  *x509.Certificate
	PEM   []byte `json:"pem"`
	Path  string `json:"path"`
	dirty bool
}

// Generate will generate a Certificate which still needs to be signed (a CSR)
func (c *Cert) Generate(subject Subject, expiry time.Duration) error {
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1),
		128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}
	if expiry < 24*time.Hour {
		expiry = 365 * 24 * time.Hour
	}

	now := time.Now()
	c.cert = &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject.AsPkixName(),
		NotBefore:    now,
		NotAfter:     now.Add(expiry),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	c.PEM = nil
	return nil
}

// Sign can be used to sign the cert (and will write to the PEM byte array)
func (c *Cert) Sign(privateKey PrivateKey, signer Pair) error {
	pubKey, err := privateKey.PublicKey()
	if err != nil {
		return err
	}
	certDER, err := x509.CreateCertificate(rand.Reader, c.cert, signer.Cert.cert,
		&pubKey, signer.PrivateKey.key)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	c.PEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	return nil
}

// Save can be used to save a Cert to disk
func (c *Cert) Save() error {
	if !c.dirty || c.Path == "" || len(c.PEM) == 0 {
		return nil
	}
	dir := path.Dir(c.Path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create path %s: %v", dir, err)
	}
	if err := os.WriteFile(c.Path, c.PEM, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %v", err)
	}
	c.dirty = false
	return nil
}
