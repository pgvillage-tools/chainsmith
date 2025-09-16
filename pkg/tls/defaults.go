package tls

import (
	"crypto/x509"
	"time"
)

var (
	DefaultExpiry = 365 * 24 * time.Hour

	// DefaultKeyUsage is used when no KeyUsages are set
	DefaultKeyUsage = x509.KeyUsageDataEncipherment |
		x509.KeyUsageDigitalSignature |
		x509.KeyUsageKeyEncipherment

	// DefaultExtendedKeyUsages is a list of extended Key usages to be used when not
	// specified in the config
	DefaultExtendedKeyUsages = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageEmailProtection,
		x509.ExtKeyUsageServerAuth,
	}
	// DefaultSubject is used when no subject is set
	DefaultSubject = Subject{
		Country:            "NL",
		CommonName:         "chainsmith",
		Locality:           "Blarocum",
		Organisation:       "Nibble-IT",
		OrganisationalUnit: "postgres",
		PostalCode:         "1261 WZ",
		State:              "Utrecht",
		StreetAddress:      "Binnendelta 1-U2",
	}
)
