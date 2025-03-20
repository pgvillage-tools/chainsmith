package tls

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCA(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tls_ca")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpDir)
	rootCertPath := path.Join(tmpDir, "test_rootCA.crt")
	rootKeyPath := path.Join(tmpDir, "test_rootCA.key")

	rootCert, rootKey, err := GenerateCA(rootCertPath, rootKeyPath, nil, nil, true)
	assert.NoError(t, err)
	assert.NotNil(t, rootCert)
	assert.NotNil(t, rootKey)
}

func TestGenerateIntermediateCA(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tls_intermediate")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, tmpDir)
	rootCertPath := path.Join(tmpDir, "test_rootCA.crt")
	rootKeyPath := path.Join(tmpDir, "test_rootCA.key")
	rootCert, rootKey, err := GenerateCA(rootCertPath, rootKeyPath, nil, nil, true)
	require.NoError(t, err)
	intermediateCertPath := path.Join(tmpDir, "test_intermediateCA.crt")
	intermediateKeyPath := path.Join(tmpDir, "test_intermediateCA.key")
	intermediateCert, intermediateKey, err := GenerateCA(intermediateCertPath, intermediateKeyPath, rootCert, rootKey, false)
	assert.NoError(t, err)
	assert.NotNil(t, intermediateCert)
	assert.NotNil(t, intermediateKey)
}

func TestGenerateCert(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "tls_intermediate")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, tmpDir)
	rootCertPath := path.Join(tmpDir, "test_rootCA.crt")
	rootKeyPath := path.Join(tmpDir, "test_rootCA.key")
	rootCert, rootKey, _ := GenerateCA(rootCertPath, rootKeyPath, nil, nil, true)
	intermediateCertPath := path.Join(tmpDir, "test_intermediateCA.crt")
	intermediateKeyPath := path.Join(tmpDir, "test_intermediateCA.key")
	intermediateCert, intermediateKey, err := GenerateCA(intermediateCertPath, intermediateKeyPath, rootCert, rootKey, false)
	require.NoError(t, err)
	require.NotNil(t, intermediateCert)
	require.NotNil(t, intermediateKey)

	testCertPath := path.Join(tmpDir, "test_server.crt")
	testKeyPath := path.Join(tmpDir, "test_server.key")
	err = GenerateCert(testCertPath, testKeyPath, intermediateCert, intermediateKey, "server.local")
	assert.NoError(t, err)
}
