package vsphereovf_test

import (
	"context"
	"os"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/onsi/ginkgo"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

// set up providers for resource tests
var acceptanceTestProvider *schema.Provider
var acceptanceTestProviders map[string]terraform.ResourceProvider

func init() {
	acceptanceTestProvider = vsphereovf.Provider().(*schema.Provider)
	acceptanceTestProviders = map[string]terraform.ResourceProvider{
		"vsphereovf": acceptanceTestProvider,
	}
}

// set up Ginkgo wrapper for running terraform's resource.Test method
// (this is so that terraform resource.Test failures will go red in Ginkgo)
type ginkgoTWrapper struct {
	ginkgo.GinkgoTInterface
}

func (g ginkgoTWrapper) Name() string {
	return ""
}

func ginkgoTestWrapper() resource.TestT {
	return ginkgoTWrapper{ginkgo.GinkgoT()}
}

// precheck for tests requiring actual infrastructure creds
func acceptanceTestPreCheck(t resource.TestT) {
	if v := os.Getenv("VSPHERE_USER"); v == "" {
		t.Fatal("VSPHERE_USER must be set for acceptance tests")
	}

	if v := os.Getenv("VSPHERE_PASSWORD"); v == "" {
		t.Fatal("VSPHERE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("VSPHERE_SERVER"); v == "" {
		t.Fatal("VSPHERE_SERVER must be set for acceptance tests")
	}
}

// utility methods to get vSphere clients and make vSphere client API calls
func getClient() *govmomi.Client {
	return acceptanceTestProvider.Meta().(*govmomi.Client)
}

func getTemplate(s *terraform.State, templatePath string) (*object.VirtualMachine, error) {
	finder := find.NewFinder(getClient().Client, false)

	// TODO: let users select a datacenter that isn't default
	dc, err := finder.DefaultDatacenter(context.Background())
	if err != nil {
		return nil, err
	}

	finder.SetDatacenter(dc)

	return finder.VirtualMachine(context.Background(), templatePath)
}
