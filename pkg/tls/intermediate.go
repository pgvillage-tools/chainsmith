package tls

import (
	"time"
)

// Intermediates holds all intermediates that are configured
type Intermediates map[string]Intermediate

// Initialize can be used to generate, build and save certificates and private
// keys for all servers and clients of all intermediates
func (i Intermediates) Initialize() (Intermediates, error) {
	for _, intermediate := range i {
		if err := intermediate.InitializeClients(); err != nil {
			return i, err
		}
		if err := intermediate.InitializeServers(); err != nil {
			return i, err
		}
	}
	return i, nil
}

// Intermediate holds the config of an intermediate, which can be either Server
// or Client (or both)
type Intermediate struct {
	cert              Pair
	Subject           *Subject `json:"subject"`
	children          Pairs
	Expiry            time.Duration `json:"root_expiry"`
	ExtendedKeyUsages ExtKeyUsages  `json:"extendedKeyUsages"`
	KeyUsages         KeyUsages     `json:"keyUsages"`
	Servers           Servers       `json:"servers"`
	Clients           []string      `json:"clients"`
}

// InitializeServers can be used to generate, build and save certificates and
// private keys for all servers an intermediate
func (i *Intermediate) InitializeServers() error {
	if i.children == nil {
		i.children = Pairs{}
	}
	for server, addresses := range i.Servers {
		subject := i.Subject.SetCommonName(server)
		pair := Pair{
			Cert: Cert{
				Subject:        &subject,
				Expiry:         i.Expiry,
				KeyUsage:       i.cert.Cert.KeyUsage,
				ExtKeyUsage:    i.cert.Cert.ExtKeyUsage,
				AlternateNames: addresses,
			},
		}
		err := pair.Process(i.cert)
		if err != nil {
			return err
		}
		i.children[server] = pair
	}
	return nil
}

// InitializeClients can be used to generate, build and save certificates and
// private keys for all clients of an intermediate
func (i *Intermediate) InitializeClients() error {
	if i.children == nil {
		i.children = Pairs{}
	}
	for _, client := range i.Clients {
		subject := i.Subject.SetCommonName(client)
		pair := Pair{
			Cert: Cert{
				Subject:     &subject,
				Expiry:      i.Expiry,
				KeyUsage:    i.cert.Cert.KeyUsage,
				ExtKeyUsage: i.cert.Cert.ExtKeyUsage,
			},
		}
		err := pair.Process(i.cert)
		if err != nil {
			return err
		}
		i.children[client] = pair
	}
	return nil
}

// ClassicIntermediates is a classical approach (list of structs with name in
// struct) to Intermediates (map of intermediates with name as key)
type ClassicIntermediates []ClassicIntermediate

// ClassicIntermediate exists for historic reasons
type ClassicIntermediate struct {
	Name string `json:"name"`
	Intermediate
}

// AsIntermediates converts a ClassicIntermediates into a Intermediates
func (cis ClassicIntermediates) AsIntermediates() Intermediates {
	i := Intermediates{}
	for _, ci := range cis {
		i[ci.Name] = ci.AsIntermediate()
	}
	return i
}

// AsIntermediate converts a ClassicIntermediate into a Intermediate
func (ci ClassicIntermediate) AsIntermediate() Intermediate {
	return Intermediate{
		Expiry:            ci.Expiry,
		ExtendedKeyUsages: ci.ExtendedKeyUsages,
		KeyUsages:         ci.KeyUsages,
		Servers:           ci.Servers,
		Clients:           ci.Clients,
	}
}

// Servers is a map holding servers, with addresses. The key will be used for
// the CommonName
type Servers map[string]ServerAddresses

// ServerAddresses is a list of DNS names and/or ip addresses to be used in the
// SAN field
type ServerAddresses []string
