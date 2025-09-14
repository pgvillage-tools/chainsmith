package tls

import (
	"time"
)

// Intermediates holds all intermediates that are configured
type Intermediates map[string]Intermediate

// Initialize can be used to generate, build and save certificates and private
// keys for all servers and clients of all intermediates
func (i Intermediates) Initialize() (Intermediates, error) {
}

// Intermediate holds the config of an intermediate, which can be either Server
// or Client (or both)
type Intermediate struct {
	// intermediateCert is the cert/privatekey of this intermediate
	intermediateCert Pair
	// certs is a collection of all certs signed by this intermediate
	certs             Pairs
	Expiry            time.Duration `json:"root_expiry"`
	ExtendedKeyUsages ExtKeyUsages  `json:"extendedKeyUsages"`
	KeyUsages         KeyUsages     `json:"keyUsages"`
	Servers           Servers       `json:"servers"`
	Clients           []string      `json:"clients"`
}

// Initialize can be used to generate, build and save certificates and private
// keys for all servers and clients of an intermediate
func (i Intermediate) Initialize() error {

}

// InitializeServers can be used to generate, build and save certificates and
// private keys for all servers an intermediate
func (i Intermediate) InitializeServers() error {

}

// InitializeClients can be used to generate, build and save certificates and
// private keys for all clients of an intermediate
func (i Intermediate) InitializeClients() error {

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
