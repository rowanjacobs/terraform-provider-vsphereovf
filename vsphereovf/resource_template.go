package vsphereovf

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/importer"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/mark"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/search"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
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
				Optional: true,
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
			"template": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

// TODO: do this in a less hacky and brittle way
func ovf(path string) string {
	base := filepath.Base(path)
	if strings.HasSuffix(base, ".ova") {
		return fmt.Sprintf("%s.ovf", strings.TrimSuffix(base, ".ova"))
	}
	return base
}

func CreateTemplate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] creating template resource: %+v\n", d)

	ovfPath, err := filepath.Abs(d.Get("path").(string))
	ovfName := ovf(ovfPath)
	d.SetId(fmt.Sprintf("%s/%s", d.Get("folder").(string), ovfName))

	name, ok := d.Get("name").(string)
	if !ok || name == "" {
		name = strings.SplitN(ovfName, ".", 2)[0]
	}

	readerProvider, err := ovx.NewReaderProvider(ovfPath)
	if err != nil {
		return err
	}

	r, _, err := readerProvider.Reader(ovfName)
	if err != nil {
		return err
	}

	ovfContents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	d.Set("contents_sha", sha1.Sum(ovfContents))

	client := m.(*govmomi.Client)
	dcPath := d.Get("datacenter").(string)
	templateParentObjects, err := search.FetchParentObjects(
		client,
		dcPath,
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
		readerProvider,
	)
	importSpec, err := i.CreateImportSpec(name, string(ovfContents), templateParentObjects.Networks)
	if err != nil {
		return fmt.Errorf("error creating import spec: %s", err)
	}

	err = i.Import(importSpec, templateParentObjects.Folder, filepath.Dir(ovfPath))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] searching for VM '%s'\n", name)
	props, err := search.VMProperties(client, dcPath, name)
	if err != nil {
		return err
	}

	d.Set("uuid", props.Config.Uuid)

	log.Printf("[DEBUG] successfully created template resource: %+v\n", d)
	if d.Get("template").(bool) {
		return mark.AsTemplate(client, dcPath, name)
	}

	return nil
}

func resourceTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] reading template resource: %+v\n", d)
	client := m.(*govmomi.Client)

	name := d.Get("name").(string)

	if name == "" {
		name = strings.SplitN(strings.SplitN(d.Id(), "/", 2)[1], ".", 2)[0]
	}

	dcPath := d.Get("datacenter").(string)

	props, err := search.VMProperties(client, dcPath, name)
	if err != nil {
		if _, ok := err.(search.NotFoundError); ok {
			log.Printf("[WARN] template not found: %+v\n", err)
			// If the VM is not found, set the ID to "".
			// This causes Terraform to regenerate the resource.
			d.SetId("")
			return nil
		}
		if _, ok := err.(search.DatacenterNotFoundError); ok {
			log.Printf("[WARN] datacenter not found: %+v\n", err)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("uuid", props.Config.Uuid)

	log.Printf("[DEBUG] successfully read template resource: %+v\n", d)
	return nil
}

func resourceTemplateDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] reading template resource: %+v\n", d)
	client := m.(*govmomi.Client)

	name := d.Get("name").(string)

	if name == "" {
		name = strings.SplitN(strings.SplitN(d.Id(), "/", 2)[1], ".", 2)[0]
	}

	dcPath := d.Get("datacenter").(string)

	props, err := search.VMProperties(client, dcPath, name)
	if err != nil {
		if _, ok := err.(search.NotFoundError); ok {
			log.Printf("[WARN] template not found: %+v\n", err)
			// If the VM is not found, set the ID to "".
			// It probably doesn't exist.
			d.SetId("")
			return nil
		}
		if _, ok := err.(search.DatacenterNotFoundError); ok {
			log.Printf("[WARN] datacenter not found: %+v\n", err)
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] VM exists with properties %+v", props)

	// TODO: any chance this resource refers to a powered-on VM?
	// if so we need to do something like this:
	// if vprops.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOff {
	// 	timeout := d.Get("shutdown_wait_timeout").(int)
	// 	if err := virtualmachine.GracefulPowerOff(client, vm, timeout, true); err != nil {
	// 		return fmt.Errorf("error shutting down virtual machine: %s", err)
	// 	}
	// }

	vm, err := search.VMByUUID(client, props.Config.Uuid)
	if err != nil {
		log.Printf("[DEBUG] error searching for VM by UUID %s: %s", props.Config.Uuid, err)
		return err
	}

	if err := destroy(vm); err != nil {
		return fmt.Errorf("error destroying virtual machine: %s", err)
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Delete complete", name)

	return nil
}

func destroy(vm *object.VirtualMachine) error {
	log.Printf("[DEBUG] Deleting virtual machine %q", vm.InventoryPath)
	task, err := vm.Destroy(context.Background())
	if err != nil {
		return err
	}
	return task.Wait(context.Background())
}

func resourceTemplateUpdate(d *schema.ResourceData, m interface{}) error { return nil }
