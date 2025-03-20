package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cmd_run")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpDir)
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
`, tmpDir)

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	config, err := loadConfig(configPath)
	require.NoError(t, err)
	if err := run(*config); err != nil {
		t.Fatalf("Application execution failed: %v", err)
	}
}
