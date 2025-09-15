package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pgvillage-tools/chainsmith/internal/config"
	"github.com/pgvillage-tools/chainsmith/pkg/tls"
)

var _ = Describe("cmd/chainsmith/issue", func() {
	var (
		cfg = config.Config{
			Intermediates: tls.ClassicIntermediates{
				tls.ClassicIntermediate{
					Name: "servers",
					Intermediate: tls.Intermediate{
						Servers: tls.Servers{
							"server1": []string{
								"1.2.3.4",
								"server1.mydomain",
							},
						},
					},
				},
				tls.ClassicIntermediate{
					Name: "clients",
					Intermediate: tls.Intermediate{
						Clients: []string{
							"user1",
							"sa2",
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
	Context("When issuing certificates", func() {
		It("Should allow creation of chain", func() {
			out, err := issue(&cfg)
			Expect(err).Error().NotTo(HaveOccurred())
			Expect(out).NotTo(BeEmpty())
			yaml := string(out)
			Expect(yaml).To(Equal("yaml"))
		})
	})
})
