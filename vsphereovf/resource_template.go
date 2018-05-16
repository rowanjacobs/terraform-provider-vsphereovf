package vsphereovf

import (
	"fmt"
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
	templateParentObjects, err := search.FetchParentObjects(
		client,
		d.Get("datacenter").(string),
		d.Get("datastore").(string),
		d.Get("folder").(string),
		d.Get("resource_pool").(string),
		d.Get("network_mappings").(map[string]interface{}),
	)
	if err != nil {
		return err
	}

	i := importer.NewImporterFromClient(client, importer.ResourcePoolImpl{templateParentObjects.ResourcePool}, templateParentObjects.Datastore)
	importSpec, err := i.CreateImportSpec(string(ovfContents), templateParentObjects.Networks)
	if err != nil {
		return fmt.Errorf("error creating import spec: %s", err)
	}

	return i.Import(importSpec, templateParentObjects.Folder, filepath.Dir(ovfPath))
}

func resourceTemplateRead(d *schema.ResourceData, m interface{}) error   { return nil }
func resourceTemplateDelete(d *schema.ResourceData, m interface{}) error { return nil }
func resourceTemplateUpdate(d *schema.ResourceData, m interface{}) error { return nil }
