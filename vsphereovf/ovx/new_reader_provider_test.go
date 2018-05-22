package ovx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx"
)

var _ = Describe("New Reader Provider", func() {
	Context("given an OVF file", func() {
		It("returns an OVF reader provider", func() {
			rp, err := ovx.NewReaderProvider("testdata/test.ovf")
			Expect(err).NotTo(HaveOccurred())

			Expect(rp).To(BeAssignableToTypeOf(ovx.OVFReaderProvider{}))
		})
	})

	Context("given an OVA file", func() {
		It("returns an OVA reader provider", func() {
			rp, err := ovx.NewReaderProvider("testdata/test.ova")
			Expect(err).NotTo(HaveOccurred())

			Expect(rp).To(BeAssignableToTypeOf(ovx.OVAReaderProvider{}))
		})
	})

	Context("given a non-OVF, non-OVA file", func() {
		It("returns an error", func() {
			_, err := ovx.NewReaderProvider("testdata/test.vmdk")
			Expect(err).To(MatchError("file 'test.vmdk' does not have .ova or .ovf extension"))
		})
	})
})
