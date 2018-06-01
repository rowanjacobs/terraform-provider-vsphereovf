package vsphereovf_test

import (
	. "github.com/onsi/ginkgo"

	"github.com/hashicorp/terraform/helper/resource"
)

var _ = Describe("OVF Template resource", func() {
	It("creates a basic vSphere template from an OVF template", func() {
		t := ginkgoTestWrapper()
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				acceptanceTestPreCheck(t)
			},
			// CheckDestroy: checkIfTemplateExistsInVSphere(false, true, "coreos_production_vmware_ovf"),
			Providers: acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVFTemplateResourceConfig,
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
			},
			// CheckDestroy: checkIfTemplateExistsInVSphere(false, true, "coreos_production_vmware_ova"),
			Providers: acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVATemplateResourceConfig,
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, true, "coreos_production_vmware_ova"),
					),
				},
			},
		})
	})
})

// TODO: this is specific to our vSphere environment.
// we should get folder, dc, ds, rp, and network name from env vars.

// TODO: what's the smallest OVF and OVA template available?

// TODO: in setup, download the coreOS templates, OR
// TODO: configure the location via environment variables
const basicVSphereOVFTemplateResourceConfig = `
resource "vsphereovf_template" "terraform-test-ovf" {
	name = "terraform-test-coreos-ovf"
	path = "../ignored/coreos_production_vmware_ova.ovf"
	folder = "khaleesi_templates"
	datacenter = "pizza-boxes-dc"
	resource_pool = "khaleesi"
	datastore = "vnx5600-pizza-2"
	template = true
	network_mappings {
		"VM Network" = "khaleesi"
	}
}
`

const basicVSphereOVATemplateResourceConfig = `
resource "vsphereovf_template" "terraform-test-ova" {
	path = "../ignored/coreos_production_vmware_ova.ova"
	folder = "khaleesi_templates"
	datacenter = "pizza-boxes-dc"
	resource_pool = "khaleesi"
	datastore = "vnx5600-pizza-2"
	template = true
	network_mappings {
		"VM Network" = "khaleesi"
	}
}
`
