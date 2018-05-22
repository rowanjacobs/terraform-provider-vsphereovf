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
			// CheckDestroy: checkIfTemplateExistsInVSphere(false),
			Providers: acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVFTemplateResourceConfig,
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, true, "coreos_production_vmware_ovf"),
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
			// CheckDestroy: checkIfTemplateExistsInVSphere(false),
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
const basicVSphereOVFTemplateResourceConfig = `
resource "vsphereovf_template" "terraform-test-ovf" {
	path = "../ignored/coreos_production_vmware_ovf.ovf"
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
