package vsphereovf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/importer"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/mark"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/search"
	"github.com/vmware/govmomi"
)

func TemplateResource() *schema.Resource {
	return &schema.Resource{
		Create: CreateTemplate,
		Read:   resourceTemplateRead,
		Delete: resourceTemplateDelete,
		Update: resourceTemplateUpdate,
		Schema: map[string]*schema.Schema{
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
			"template": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func CreateTemplate(d *schema.ResourceData, m interface{}) error {
	ovfPath, err := filepath.Abs(d.Get("path").(string))
	ovfName := filepath.Base(ovfPath)
	d.SetId(fmt.Sprintf("%s/%s", d.Get("folder").(string), ovfName))

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

	i := importer.NewImporterFromClient(
		client,
		templateParentObjects.ResourcePool,
		templateParentObjects.Datastore,
	)
	importSpec, err := i.CreateImportSpec(string(ovfContents), templateParentObjects.Networks)
	if err != nil {
		return fmt.Errorf("error creating import spec: %s", err)
	}

	err = i.Import(importSpec, templateParentObjects.Folder, filepath.Dir(ovfPath))
	if err != nil {
		return err
	}

	name := strings.SplitN(ovfName, ".", 2)[0]
	if d.Get("template").(bool) {
		return mark.AsTemplate(client, d.Get("datacenter").(string), name)
	}
	return nil
}

func resourceTemplateRead(d *schema.ResourceData, m interface{}) error   { return nil }
func resourceTemplateDelete(d *schema.ResourceData, m interface{}) error { return nil }
func resourceTemplateUpdate(d *schema.ResourceData, m interface{}) error { return nil }
