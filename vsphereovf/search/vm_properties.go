package search

import (
	"context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

func VMProperties(client *govmomi.Client, dcPath, name string) (*mo.VirtualMachine, error) {
	f, err := NewFinder(client, dcPath)
	if err != nil {
		return nil, err
	}

	vm, err := f.VirtualMachine(name)
	if err != nil {
		return nil, err
	}

	var props mo.VirtualMachine
	if err := vm.Properties(context.Background(), vm.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}
