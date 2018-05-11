package vsphereovf_test

import (
	"errors"
	"fmt"
	"regexp"

	. "github.com/onsi/ginkgo"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var _ = Describe("OVF Template resource", func() {
	It("creates a basic vSphere template", func() {
		resource.Test(ginkgoTestWrapper(), resource.TestCase{
			// PreCheck: func() {
			// 	acceptanceTestPreCheck(t)
			// },
			// CheckDestroy: checkIfTemplateExistsInVSphere(false),
			Providers: acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVFTemplateResourceConfig,
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, "terraform-test-ovf"),
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
			fmt.Printf("received error %s\n", err)
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

const basicVSphereOVFTemplateResourceConfig = `
resource "vsphereovf_template" "terraform-test-ovf" {
	name = "some-template"
	path = "some-ovf-path"
}
`
