package vsphereovf_test

import (
	"errors"
	"regexp"

	. "github.com/onsi/ginkgo"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
						checkIfTemplateExistsInVSphere(true, "coreos_production_vmware_ova"),
					),
				},
			},
		})
	})
})

func checkIfTemplateExistsInVSphere(expected bool, fullPath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getTemplate(s, fullPath)
		if err != nil {
			if ok, _ := regexp.MatchString("virtual machine with UUID \"[-a-f0-9]+\" not found", err.Error()); ok && !expected {
				// Expected missing
				return nil
			}
			return err
		}
		if !expected {
			return errors.New("expected VM to be missing")
		}
		return nil
	}
}

// TODO: this is specific to our vSphere environment.
// we should get folder, dc, ds, rp, and network name from env vars.
const basicVSphereOVFTemplateResourceConfig = `
resource "vsphereovf_template" "terraform-test-ovf" {
	name = "coreos_production_vmware_ova"
	path = "../ignored/coreos_production_vmware_ova.ovf"
	folder = "khaleesi_templates"
	datacenter = "pizza-boxes-dc"
	resource_pool = "khaleesi"
	datastore = "vnx5600-pizza-2"
	network_mappings {
		"VM Network" = "khaleesi"
	}
}
`
