package vsphereovf_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf"
)

var _ = Describe("ovf Template resource", func() {
	var govmomiClientStub *vsphereovf.GovmomiClient = nil

	Describe("CreateTemplate", func() {
		It("sets ID on the resource data object", func() {
			resourceData := vsphereovf.TemplateResource().Data(nil)
			err := resourceData.Set("name", "some-name")
			Expect(err).NotTo(HaveOccurred())

			err = vsphereovf.CreateTemplate(resourceData, govmomiClientStub)
			Expect(err).NotTo(HaveOccurred())
			Expect(resourceData.Id()).To(Equal("some-name"))
		})

		It("", func() {
		})
	})
})
