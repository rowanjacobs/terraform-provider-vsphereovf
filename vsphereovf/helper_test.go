package vsphereovf_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/onsi/ginkgo"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
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

const resourceTemplateFatal = "%s must be set for template resource acceptance tests"

func resourceTemplateTestPreCheck(t resource.TestT) {
	if v := os.Getenv("VSPHERE_OVA_PATH"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_OVA_PATH"))
	}

	if v := os.Getenv("VSPHERE_OVF_PATH"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_OVF_PATH"))
	}

	if v := os.Getenv("VSPHERE_FOLDER"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_FOLDER"))
	}

	if v := os.Getenv("VSPHERE_DATACENTER"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_DATACENTER"))
	}

	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_RESOURCE_POOL"))
	}

	if v := os.Getenv("VSPHERE_DATASTORE"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_DATASTORE"))
	}

	if v := os.Getenv("VSPHERE_NETWORK"); v == "" {
		t.Fatal(fmt.Sprintf(resourceTemplateFatal, "VSPHERE_NETWORK"))
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

func checkIfTemplateExistsInVSphere(expected bool, expectTemplate bool, fullPath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vm, err := getTemplate(s, fullPath)
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

		var o mo.VirtualMachine

		err = vm.Properties(context.Background(), vm.Reference(), []string{"config.template"}, &o)
		if err != nil {
			return err
		}

		if o.Config.Template != expectTemplate {
			if expectTemplate { //            v look at this alignment v
				return errors.New("expected VM to be template but it was not")
			}
			return errors.New("expected VM not to be template but it was")
		}

		return nil
	}
}
