package vsphereovf

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreos/etcd/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

func TemplateResource() *schema.Resource {
	return &schema.Resource{
		Create: CreateTemplate,
		Read:   resourceTemplateRead,
		Delete: resourceTemplateDelete,
		Update: resourceTemplateUpdate,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func CreateTemplate(d *schema.ResourceData, m interface{}) error {
	// Set the ID of the Terraform resource.
	// (Currently we just set the ID to the template name.)
	// TODO: (Path/name might be a better option for this.)
	name := d.Get("name").(string)
	d.SetId(name)

	ctx := context.TODO()
	// read and unmarshal the OVF file into an OVF envelope
	contents, err := ioutil.ReadFile(ovfPath)
	envelope, err := ovf.Unmarshal(bytes.NewReader(contents))

	// represent the OVF envelope's network mappings as a map[string]string
	for _, net := range envelope.Network.Networks {
		networks[net.Name] = net.Name
	}

	// make a CreateImportSpecParams, give it an empty network mapping
	isp := types.OvfCreateImportSpecParams{NetworkMapping: []types.OvfNetworkMapping{}}

	// populate the network mapping
	for src, dst := range networks {
		net, err := Network(client, dc, dst)
		isp.NetworkMapping = append(isp.NetworkMapping, types.OvfNetworkMapping{Name: src, Network: net.Reference()})
	}

	// create an ovf manager, use it to create an import spec out of our CreateImportSpecParams
	manager := ovf.NewManager(client.Client)
	importSpec, err := manager.CreateImportSpec(ctx, string(contents), resourcePool, dataStore, isp)

	// together, the following lines are a sort of blocking ImportVApp
	// use our import spec to get an nfc lease from our resource pool
	lease, err := resourcePool.ImportVApp(ctx, importSpec.ImportSpec, folder, nil)
	// use our import spec to get lease info (including list of item URLs) out of the nfc lease
	info, err := lease.Wait(ctx, importSpec.FileItem)

	// read ovf file from local path
	// might look like:
	file, err := os.Open(filepath.Join(filepath.Dir(ovfPath), item.Path))

	// use the nfc lease to upload the ovf file
	// looks kinda like this:
	err = lease.Upload(ctx, info.Items[0], file, soap.Upload{ContentLength: fileInfo.Size()})

	return nil
}

func resourceTemplateRead(d *schema.ResourceData, m interface{}) error   { return nil }
func resourceTemplateDelete(d *schema.ResourceData, m interface{}) error { return nil }
func resourceTemplateUpdate(d *schema.ResourceData, m interface{}) error { return nil }
