package ovx_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx"
)

var _ = Describe("OVF Reader Provider", func() {
	var (
		provider ovx.ReaderProvider
	)

	Context("given an ovf file with adjacent assets", func() {
		BeforeEach(func() {
			wd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			testOvfPath := filepath.Join(wd, "testdata/test.ovf")
			provider, err = ovx.NewOVFReaderProvider(testOvfPath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a reader of that ovf file when asked", func() {
			r, size, err := provider.Reader("test.ovf")
			Expect(err).NotTo(HaveOccurred())
			Expect(size).To(BeNumerically(">", 0))

			Expect(gbytes.BufferReader(r)).To(gbytes.Say(".*"))
		})

		It("can return other files when asked", func() {
			r, size, err := provider.Reader("test.vmdk")
			Expect(err).NotTo(HaveOccurred())
			Expect(size).To(BeNumerically(">", 0))

			Expect(gbytes.BufferReader(r)).To(gbytes.Say(".*"))
		})

		Context("when asked for an asset that does not exist", func() {
			It("errors", func() {
				_, _, err := provider.Reader("nonexistentfile.wtf")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("given a ovf file that does not exist", func() {
		It("errors", func() {
			_, err := ovx.NewOVFReaderProvider("/this/path/is/not/real")
			Expect(err).To(HaveOccurred())
		})
	})
})
