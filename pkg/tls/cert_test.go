package tls

import (
	"fmt"
	"net"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cert", func() {
	const (
		email_s  = "user@example.org"
		domain_s = "someserver1.subdomain.example.org"
		ipv4_1s  = "1.2.3.4"
		ipv4_2s  = "127.0.0.1"
		ipv6_1s  = "::1"
		ipv6_2s  = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		ipv6_3s  = "2001:db8::1:0:0:1"
		ipv6_4s  = "::ffff:192.0.2.128"
		uri_1s   = "https://private.api.example.org"
		uri_2s   = "http://public.api.example.org:8080"

		// This becomes a uri
		//broken_email_s  = "user@@example.org"
		//broken_domain_s = "someserver1.subdomain.example.org."
		//broken_uri_s    = "https:private.api.example.org"
		broken_ipv4_s = "1.2.3:4"
		broken_ipv6_s = ipv6_2s + ":"
	)
	var (
		emails           = []string{email_s}
		domains          = []string{domain_s}
		ipv4_slist       = []string{ipv4_1s, ipv4_2s}
		ipv4_1           = net.ParseIP(ipv4_1s)
		ipv4_2           = net.ParseIP(ipv4_2s)
		ipv6_slist       = []string{ipv6_1s, ipv6_2s, ipv6_3s, ipv6_4s}
		ipv6_1           = net.ParseIP(ipv6_1s)
		ipv6_2           = net.ParseIP(ipv6_2s)
		ipv6_3           = net.ParseIP(ipv6_3s)
		ipv6_4           = net.ParseIP(ipv6_4s)
		uri_slist        = []string{uri_1s, uri_2s}
		uri_must_compile = func(s string) *url.URL {
			uri, err := url.Parse(s)
			if err != nil {
				panic(fmt.Sprintf("%s could not be compiled: %e", s, err))
			}
			return uri
		}
		uri_1           = uri_must_compile(uri_1s)
		uri_2           = uri_must_compile(uri_2s)
		workingAltNames = append(
			append(
				append(
					append(emails, domains...),
					ipv4_slist...,
				),
				ipv6_slist...,
			),
			uri_slist...)
		invalidAltNames = []string{
			broken_ipv4_s,
			broken_ipv6_s,
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
	Context("When splitting Alternate Names ", Ordered, func() {
		It("Should work with valid values", func() {
			result, err := splitAlternateNames(workingAltNames)
			Expect(err).Error().NotTo(HaveOccurred())

			Expect(result.dnsNames).To(HaveLen(len(domains)))
			Expect(result.dnsNames).To(ContainElements(domain_s))

			Expect(result.eMailAddresses).To(HaveLen(len(emails)))
			Expect(result.eMailAddresses).To(ContainElements(email_s))

			Expect(result.ipAddresses).To(HaveLen(len(ipv4_slist) + len(ipv6_slist)))
			Expect(result.ipAddresses).To(ContainElements(ipv4_1, ipv4_2, ipv6_1, ipv6_2, ipv6_3, ipv6_4))

			Expect(result.uris).To(HaveLen(len(uri_slist)))
			Expect(result.uris).To(ContainElements(uri_2))
			Expect(result.uris).To(ContainElements(uri_1))
		})
		It("Should not work with invalid values", func() {
			for _, invalid := range invalidAltNames {
				Expect(fmt.Fprintf(GinkgoWriter, "DEBUG - Test: %v\n", invalid)).
					Error().NotTo(HaveOccurred())
				myAltNames := append(workingAltNames, invalid)
				_, err := splitAlternateNames(myAltNames)
				Expect(err).Error().To(HaveOccurred())
			}
		})
	})
})
