package tls

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PrivateKey", func() {

	/*
		BeforeAll(func() {
		})
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
	*/
	Context("When generating private keys", Ordered, func() {
		pKey := PrivateKey{}
		It("Should generate successfully", func() {
			Expect(pKey.key).To(BeNil())
			err := pKey.Generate()
			Expect(err).Error().NotTo(HaveOccurred())
			Expect(pKey.key).NotTo(BeNil())
		})
		It("Should not generate once generation has been done before", func() {
			Expect(pKey.key).NotTo(BeNil())
			preGenerateKey := pKey.key
			err := pKey.Generate()
			Expect(err).Error().NotTo(HaveOccurred())
			Expect(pKey.key).To(BeIdenticalTo(preGenerateKey))
		})
		It("Should encode successfully", func() {
			Expect(pKey.PEM).To(BeNil())
			err := pKey.Encode()
			Expect(err).Error().NotTo(HaveOccurred())
			Expect(pKey.PEM).NotTo(BeNil())
			Expect(string(pKey.PEM)).To(HavePrefix("-----BEGIN RSA PRIVATE KEY-----"))
			Expect(string(pKey.PEM)).To(HaveSuffix("-----END RSA PRIVATE KEY-----\n"))
		})
		It("Should save successfully", func() {
			var PEM []byte
			Expect(pKey.dirty).To(BeTrue())
			tmpDir, err := os.MkdirTemp("", "privateKeyTest")
			defer func() {
				Expect(os.RemoveAll(tmpDir)).Error().NotTo(HaveOccurred())
			}()
			Expect(err).Error().NotTo(HaveOccurred())
			pKeyPath := path.Join(tmpDir, "private_key")
			pKey.Path = pKeyPath
			err = pKey.Save()
			Expect(err).Error().NotTo(HaveOccurred())

			PEM, err = os.ReadFile(pKeyPath)
			Expect(err).Error().NotTo(HaveOccurred())
			Expect(PEM).To(Equal(pKey.PEM))
			Expect(string(PEM)).To(HavePrefix("-----BEGIN RSA PRIVATE KEY-----"))
			Expect(string(PEM)).To(HaveSuffix("-----END RSA PRIVATE KEY-----\n"))
		})
	})
})
