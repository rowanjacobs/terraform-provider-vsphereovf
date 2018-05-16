package vsphereovf

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/internal/importer"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/internal/search"
	"github.com/vmware/govmomi"
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
			"network_mappings": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"folder": &schema.Schema{
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
	d.SetId(d.Get("name").(string))
	// stuff above this comment is unit tested

	// TODO: find a way to unit test the rest of this function.
	ovfPath, err := filepath.Abs(d.Get("path").(string))
	ovfContents, err := ioutil.ReadFile(ovfPath)

	client := m.(*govmomi.Client)
	finder, err := search.NewFinder(client, d.Get("datacenter").(string))
	if err != nil {
		return err
	}

	resourcePool, err := finder.ResourcePool(d.Get("resource_pool").(string))
	if err != nil {
		return err
	}

	datastore, err := finder.Datastore(d.Get("datastore").(string))
	if err != nil {
		return err
	}

	folder, err := finder.Folder(d.Get("folder").(string))
	if err != nil {
		return err
	}

	i := importer.NewImporterFromClient(client, finder, importer.ResourcePoolImpl{resourcePool}, datastore)
	importSpec, err := i.CreateImportSpec(string(ovfContents), d.Get("network_mappings").(map[string]interface{}))
	if err != nil {
		return err
	}

	return i.Import(importSpec, folder, filepath.Dir(ovfPath))
}

func resourceTemplateRead(d *schema.ResourceData, m interface{}) error   { return nil }
func resourceTemplateDelete(d *schema.ResourceData, m interface{}) error { return nil }
func resourceTemplateUpdate(d *schema.ResourceData, m interface{}) error { return nil }
