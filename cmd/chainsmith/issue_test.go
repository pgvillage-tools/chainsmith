package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pgvillage-tools/chainsmith/internal/config"
)

var _ = Describe("cmd/chainsmith/issue", func() {
	var (
		cfg = config.Config{}
	)
	BeforeAll(func() {
	})
	BeforeEach(func() {
	})
	AfterEach(func() {
	})
	Context("When issuing certificates", func() {
		It("Should allow creation of a valid Paas", func() {
			issue(cfg)
		})
	})
})
