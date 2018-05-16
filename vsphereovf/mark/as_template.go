package mark

import (
	"context"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/search"
	"github.com/vmware/govmomi"
)

func AsTemplate(client *govmomi.Client, dcPath, name string) error {
	finder, err := search.NewFinder(client, dcPath)
	if err != nil {
		return err
	}

	vm, err := finder.VirtualMachine(name)
	if err != nil {
		return err
	}

	return vm.MarkAsTemplate(context.Background())
}
