package config

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cmd_run")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, tmpDir)
	configPath := path.Join(tmpDir, "config.yml")
	configContent := fmt.Sprintf(`
root_ca_path: "%[1]s/test_rootCA.crt"
intermediate_ca_path: "%[1]s/test_intermediateCA.crt"
certificates:
  server:
    cert_path: "%[1]s/test_server.crt"
    key_path: "%[1]s/test_server.key"
    common_name: "%[1]s/server.local"
  client:
    cert_path: "%[1]s/test_client.crt"
    key_path: "%[1]s/test_client.key"
    common_name: "%[1]s/client.local"
subject:
  C: NL/postalCode=1261 WZ
  CN: postgres
intermediates:
  - name: server
    extendedKeyUsages:
      - serverAuth
    keyUsages:
      - digitalSignature
    servers:
      server1:
        - server1.local
        - 1.2.3.4
      server2:
        - server2.local
        - 1.2.3.5
  - name: client
    clients:
      - client1
      - client2
      - client3
    extendedKeyUsages:
      - clientAuth
    keyUsages:
      - keyEncipherment
`, tmpDir)

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, config.RootCAPath, tmpDir+"/test_rootCA.crt")
	assert.Equal(t, config.IntermediateCAPath,
		tmpDir+"/test_intermediateCA.crt")
	assert.Len(t, config.Certificates, 2)
	require.Contains(t, config.Certificates, "server")
	server := config.Certificates["server"]
	assert.Equal(t, server.CertPath, tmpDir+"/test_server.crt")
	assert.Equal(t, server.KeyPath, tmpDir+"/test_server.key")
	assert.Equal(t, server.CommonName, tmpDir+"/server.local")

	require.Contains(t, config.Certificates, "client")
	client := config.Certificates["client"]
	assert.Equal(t, client.CertPath, tmpDir+"/test_client.crt")
	assert.Equal(t, client.KeyPath, tmpDir+"/test_client.key")
	assert.Equal(t, client.CommonName, tmpDir+"/client.local")

	assert.Len(t, config.RootSubject, 2)
	assert.Contains(t, config.RootSubject, "C")
	assert.Contains(t, config.RootSubject, "CN")

	assert.Len(t, config.Intermediates, 2)
	serverInt := config.Intermediates[0]
	assert.Equal(t, serverInt.Name, "server")
	assert.Len(t, serverInt.Servers, 2)
	assert.Len(t, serverInt.Clients, 0)
	assert.Contains(t, serverInt.ExtendedKeyUsages, "serverAuth")
	assert.Contains(t, serverInt.KeyUsages, "digitalSignature")

	clientInt := config.Intermediates[0]
	assert.Equal(t, clientInt.Name, "client")
	assert.Len(t, clientInt.Servers, 0)
	assert.Len(t, clientInt.Clients, 3)
	assert.Contains(t, clientInt.ExtendedKeyUsages, "clientAuth")
	assert.Contains(t, clientInt.KeyUsages, "keyEncipherment")
}
