package config

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("When loading config", func() {
		var config1 = `
chain:
  intermediates:
    server:
      servers:
        server1.local:
          - server1.local
    client:
      clients:
      - client1
  root:
    cert:
      subject:
        CN: chainsmith
tmpdir: /tmp/pgvillage/chainsmith/chain1/
`
		var config2 = `
intermediates:
  - name: server
    servers:
      server1.local:
        - server1.local
  - clients:
    - client1
    name: client
subject:
  CN: chainsmith
tmpdir: /tmp/pgvillage/chainsmith/chain1/
`
		It("should return config as expected", func() {
			tmpDir, err := os.MkdirTemp("", "configTest")
			Expect(err).Error().NotTo(HaveOccurred())
			defer func() {
				Expect(os.RemoveAll(tmpDir)).Error().NotTo(HaveOccurred())
			}()
			confPath := path.Join(tmpDir, "config.yaml")
			for _, config := range []string{config1, config2} {
				Expect(os.WriteFile(confPath, []byte(config), 0600)).Error().NotTo(HaveOccurred())
				config, err := LoadConfig(confPath)
				Expect(err).Error().NotTo(HaveOccurred())
				Expect(config).NotTo(BeNil())
				//Expect(config.TmpDir).To(Equal(os.Getenv("TMPDIR")))
				Expect(config.TmpDir).To(Equal("/tmp/pgvillage/chainsmith/chain1/"))
				chain, err := config.AsChain()
				Expect(err).Error().NotTo(HaveOccurred())
				Expect(chain.Intermediates).To(HaveLen(2))
				Expect(chain.Root.Cert.Subject.CommonName).To(Equal("chainsmith"))
			}
		})
	})
	Context("When converting a config to a Chain", func() {
		It("should return chain as expected", func() {
		})
	})
})
