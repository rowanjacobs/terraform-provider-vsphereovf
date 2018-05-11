package vsphereovf

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/govmomi"
)

var newGovmomiClient = govmomi.NewClient

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user name for vSphere API operations.",
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USER", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user password for vSphere API operations.",
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
			},
			"vsphere_server": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vSphere Server name for vSphere API operations.",
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_SERVER", nil),
			},
			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set, VMware vSphere client will permit unverifiable SSL certificates.",
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_ALLOW_UNVERIFIED_SSL", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vsphereovf_template": TemplateResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	user := d.Get("user").(string)
	password := d.Get("password").(string)
	host := d.Get("vsphere_server").(string)

	u, err := url.Parse("https://" + host + "/sdk")
	if err != nil {
		return nil, fmt.Errorf("error parsing vSphere server URL: %s", err)
	}

	u.User = url.UserPassword(user, password)

	insecure := d.Get("allow_unverified_ssl").(bool)
	return newGovmomiClient(context.Background(), u, insecure)
}
