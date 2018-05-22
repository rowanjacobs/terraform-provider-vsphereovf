package ovx_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx"
)

var _ = Describe("OVA Reader Provider", func() {
	var (
		provider ovx.ReaderProvider
	)

	Context("given a valid ova file", func() {
		BeforeEach(func() {
			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			testOvaPath := filepath.Join(wd, "testdata/test.ova")
			provider, err = ovx.NewOVAReaderProvider(testOvaPath)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when asked for the embedded ovf file", func() {
			It("returns a reader", func() {
				r, size, err := provider.Reader("test.ovf")
				Expect(err).NotTo(HaveOccurred())

				Expect(r).NotTo(BeNil())

				Eventually(gbytes.BufferReader(r)).Should(gbytes.Say("<xml></xml>"))

				Expect(size).To(BeNumerically(">", 0))
			})
		})

		Context("when asked for a file that is not embedded", func() {
			It("returns an error", func() {
				_, _, err := provider.Reader("invalid.wtf")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("given a nonexistent file", func() {
		It("returns an error", func() {
			_, err := ovx.NewOVAReaderProvider("/this/file/isn't/real")
			Expect(err).To(HaveOccurred())
		})
	})
})
