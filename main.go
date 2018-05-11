package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vsphereovf.Provider})
}
