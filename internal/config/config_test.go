package config

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", Ordered, func() {
	Context("When loading config", func() {
		var tmpDir string
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
		BeforeAll(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "configTest")
			Expect(err).Error().NotTo(HaveOccurred())
		})
		AfterAll(func() {
			Expect(os.RemoveAll(tmpDir)).Error().NotTo(HaveOccurred())
		})
		It("should return config as expected", func() {
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
		It("should raise error when file does not exist", func() {
			confPath := path.Join(tmpDir, "does_not_exist.yaml")
			_, err := LoadConfig(confPath)
			Expect(err).Error().To(HaveOccurred())
		})
		It("should raise error when file does not parse as yaml", func() {
			confPath := path.Join(tmpDir, "invalid.yaml")
			Expect(
				os.WriteFile(confPath, []byte("this is not valid yaml"), 0600),
			).Error().NotTo(HaveOccurred())
			_, err := LoadConfig(confPath)
			Expect(err).Error().To(HaveOccurred())
		})
	})
	Context("When converting a config to a Chain", func() {
		It("should return chain as expected", func() {
		})
	})
})
