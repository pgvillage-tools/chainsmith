package tls

import (
	"time"
)

// Chain can hold all configuration for a chain.
type Chain struct {
	ca            Pair
	Subject       Subject       `json:"subject"`
	Intermediates Intermediates `json:"intermediates"`
	Expiry        time.Duration `json:"expiry"`
	// Path where all files are stored
	Store string `json:"store"`
	Keys  Key    `json:"keys"`
}

// Key represents a pair of private and public key used to encrypt and decrypt
// the private keys belonging to the certificates
type Key struct {
	// To decrypt
	PrivateKey string `json:"private"`
	// To encrypt
	PublicKey string
}

// InitializeCA can be used to generate, build and save the CA cert and private
// key
func (c *Chain) InitializeCA() error {
	if err := c.ca.Generate(c.Subject, c.Expiry); err != nil {
		return err
	} else if err := c.ca.Sign(c.ca); err != nil {
		return err
	} else if err := c.ca.Encode(); err != nil {
		return err
	}
	return c.ca.Save()
}

// InitializeIntermediates can be used to initiaze all initermediates belonging
// to this chain
func (c *Chain) InitializeIntermediates() error {
	var err error
	c.Intermediates, err = c.Intermediates.Initialize()
	return err
}
