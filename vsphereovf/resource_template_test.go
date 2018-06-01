package vsphereovf_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"

	"github.com/hashicorp/terraform/helper/resource"
)

var _ = Describe("OVF Template resource", func() {
	// TODO: what's the smallest OVF and OVA template available?
	// TODO: in setup, download the coreOS templates if env vars are empty

	It("creates a basic vSphere template from an OVF template", func() {
		t := ginkgoTestWrapper()
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				acceptanceTestPreCheck(t)
				resourceTemplateTestPreCheck(t)
			},
			// CheckDestroy: checkIfTemplateExistsInVSphere(false, true, "coreos_production_vmware_ovf"),
			Providers: acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVFTemplateResourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, true, "terraform-test-coreos-ovf"),
					),
				},
			},
		})
	})

	It("creates a basic vSphere template from an OVA template", func() {
		t := ginkgoTestWrapper()
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				acceptanceTestPreCheck(t)
				resourceTemplateTestPreCheck(t)
			},
			// CheckDestroy: checkIfTemplateExistsInVSphere(false, true, "coreos_production_vmware_ova"),
			Providers: acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVATemplateResourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, true, "coreos_production_vmware_ova"),
					),
				},
			},
		})
	})
})

func basicVSphereOVFTemplateResourceConfig() string {
	template := `
resource "vsphereovf_template" "terraform-test-ovf" {
	name = "terraform-test-coreos-ovf"
	path = "%s"
	folder = "%s"
	datacenter = "%s"
	resource_pool = "%s"
	datastore = "%s"
	template = true
	network_mappings {
		"VM Network" = "%s"
	}
}
`
	return fmt.Sprintf(template,
		os.Getenv("VSPHERE_OVF_PATH"),
		os.Getenv("VSPHERE_FOLDER"),
		os.Getenv("VSPHERE_DATACENTER"),
		os.Getenv("VSPHERE_RESOURCE_POOL"),
		os.Getenv("VSPHERE_DATASTORE"),
		os.Getenv("VSPHERE_NETWORK"),
	)
}

func basicVSphereOVATemplateResourceConfig() string {
	template := `
resource "vsphereovf_template" "terraform-test-ova" {
	path = "%s"
	folder = "%s"
	datacenter = "%s"
	resource_pool = "%s"
	datastore = "%s"
	template = true
	network_mappings {
		"VM Network" = "%s"
	}
}
`
	return fmt.Sprintf(template,
		os.Getenv("VSPHERE_OVA_PATH"),
		os.Getenv("VSPHERE_FOLDER"),
		os.Getenv("VSPHERE_DATACENTER"),
		os.Getenv("VSPHERE_RESOURCE_POOL"),
		os.Getenv("VSPHERE_DATASTORE"),
		os.Getenv("VSPHERE_NETWORK"),
	)
}
