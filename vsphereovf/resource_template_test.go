package vsphereovf_test

import (
	. "github.com/onsi/ginkgo"

	"github.com/hashicorp/terraform/helper/resource"
)

var _ = Describe("OVF Template resource", func() {
	It("creates a basic vSphere template", func() {
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
