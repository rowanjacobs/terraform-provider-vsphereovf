package vsphereovf_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"

	"github.com/hashicorp/terraform/helper/resource"
)

var _ = Describe("OVF Template resource", func() {
	// TODO: in setup, download the coreOS templates if env vars are empty

	It("creates a basic vSphere template from an OVF template", func() {
		t := ginkgoTestWrapper()
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				acceptanceTestPreCheck(t)
				resourceTemplateTestPreCheck(t)
			},
			CheckDestroy: checkIfTemplateExistsInVSphere(false, true, inventoryPath("terraform-test-coreos-ovf")),
			Providers:    acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVFTemplateResourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, true, inventoryPath("terraform-test-coreos-ovf")),
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
			CheckDestroy: checkIfTemplateExistsInVSphere(false, true, inventoryPath("coreos_production_vmware_ova")),
			Providers:    acceptanceTestProviders,
			Steps: []resource.TestStep{
				{
					Config: basicVSphereOVATemplateResourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkIfTemplateExistsInVSphere(true, true, inventoryPath("coreos_production_vmware_ova")),
					),
				},
			},
		})
	})

	// TODO: (this fails and I never could quite figure out why)
	Context("when another Terraform provider is being used", func() {
		It("can read a vSphere template created by this provider", func() {
			t := ginkgoTestWrapper()
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acceptanceTestPreCheck(t)
					resourceTemplateTestPreCheck(t)
				},
				CheckDestroy: checkIfTemplateExistsInVSphere(false, false, inventoryPath("terraform-test-coreos-vm")),
				Providers:    acceptanceTestProvidersWithVSphere,
				Steps: []resource.TestStep{
					{
						Config: readTemplateResourceConfig(),
						Check: resource.ComposeTestCheckFunc(
							checkIfTemplateExistsInVSphere(true, false, inventoryPath("terraform-test-coreos-vm")),
						),
					},
				},
			})
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

// partially copied from pivotal-cf/terraforming-vsphere,
// partially from the Terraform vsphere/r/virtual_machine docs
func readTemplateResourceConfig() string {
	template := `
data "vsphere_datacenter" "dc" {
  name = "%s"
}

data "vsphere_resource_pool" "pool" {
  name          = "%s"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "ds" {
  name          = "%s"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "network" {
  name          = "%s"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

resource "vsphereovf_template" "terraform-test-ova" {
	path = "%s"
	folder = "%s"
	datacenter = "${data.vsphere_datacenter.dc.name}"
	resource_pool = "${data.vsphere_resource_pool.pool.name}"
	datastore = "${data.vsphere_datastore.ds.name}"
	template = true
	network_mappings {
		"VM Network" = "${data.vsphere_network.network.name}"
	}
}

resource "vsphere_virtual_machine" "vm" {
  name = "terraform-test-coreos-vm"
  resource_pool_id = "${data.vsphere_resource_pool.pool.id}"
  datastore_id     = "${data.vsphere_datastore.ds.id}"

	num_cpus = 2
	memory = 1024
	guest_id = "other26xLinux64Guest"

  scsi_type = "pvscsi"

  network_interface {
    network_id   = "${data.vsphere_network.network.id}"
  }

  disk {
    name             = "disk0.vmdk"
    size             = "10"
  }

  clone {
    template_uuid = "${vsphereovf_template.terraform-test-ova.uuid}"
  }

  vapp {
    properties {
      "guestinfo.hostname"                        = "terraform-test.foobar.local"
      "guestinfo.interface.0.name"                = "ens192"
    }
  }
}
`
	return fmt.Sprintf(template,
		os.Getenv("VSPHERE_DATACENTER"),
		os.Getenv("VSPHERE_RESOURCE_POOL"),
		os.Getenv("VSPHERE_DATASTORE"),
		os.Getenv("VSPHERE_NETWORK"),
		os.Getenv("VSPHERE_OVA_PATH"),
		os.Getenv("VSPHERE_FOLDER"),
	)
}
