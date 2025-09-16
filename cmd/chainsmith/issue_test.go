package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pgvillage-tools/chainsmith/internal/config"
	"github.com/pgvillage-tools/chainsmith/pkg/tls"
)

var _ = Describe("cmd/chainsmith/issue", func() {
	const (
		int1    = "servers"
		int2    = "clients"
		s1      = "server1"
		s1_ip1  = "1.2.3.4"
		s1_dns1 = "server1.mydomain"
		c1      = "user1"
		c2      = "sa2"
	)
	var (
		cfg = config.Config{
			Intermediates: tls.ClassicIntermediates{
				tls.ClassicIntermediate{
					Name: int1,
					Intermediate: tls.Intermediate{
						Servers: tls.Servers{
							s1: []string{
								s1_ip1,
								s1_dns1,
							},
						},
					},
				},
				tls.ClassicIntermediate{
					Name: int2,
					Intermediate: tls.Intermediate{
						Clients: []string{
							c1,
							c2,
						},
					},
				},
			},
		}
	)
	/*
		BeforeAll(func() {
		})
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
	*/
	Context("When issuing certificates with classic config", func() {
		It("Should allow creation of chain", func() {
			chain, err := issue(&cfg)
			Expect(err).Error().NotTo(HaveOccurred())
			Expect(chain).NotTo(BeNil())
			structure := chain.Structure()

			Expect(structure.Certs).To(HaveKey(int1))
			int1_structure := structure.Certs[int1]
			Expect(int1_structure).To(HaveKey(s1))

			Expect(structure.Certs).To(HaveKey(int2))
			int2_structure := structure.Certs[int2]
			Expect(int2_structure).To(HaveKey(c1))
			Expect(int2_structure).To(HaveKey(c2))
		})
	})
})
