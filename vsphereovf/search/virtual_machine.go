package search

import (
	"context"
	"fmt"
	"log"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
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

func VMByPath(client *govmomi.Client, dc *object.Datacenter, path string) (*object.VirtualMachine, error) {
	finder := find.NewFinder(client.Client, false)
	finder.SetDatacenter(dc)

	return finder.VirtualMachine(context.Background(), path)
}

func VMByUUID(client *govmomi.Client, uuid string) (*object.VirtualMachine, error) {
	// the vCenter that we target in tests is old (pre-6.5).
	// container view is the only pre-6.5-supported way of searching by UUID which includes templates.
	// TODO: add support for using UUID search index for users with vSphere version 6.5 and later
	ctx := context.Background()

	// this code is slightly adapted, mostly taken from the vSphere Terraform provider, written by vancleuver.
	log.Printf("[DEBUG] Using ContainerView to look up UUID %q", uuid)

	m := view.NewManager(client.Client)

	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = v.Destroy(ctx); err != nil {
			log.Printf("[DEBUG] virtualMachineFromContainerView: Unexpected error destroying container view: %s", err)
		}
	}()

	var vms, results []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"config.uuid"}, &results)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result.Config == nil {
			continue
		}
		if result.Config.Uuid == uuid {
			vms = append(vms, result)
		}
	}

	switch {
	case len(vms) < 1:
		return nil, fmt.Errorf("virtual machine with UUID %q not found", uuid)
	case len(vms) > 1:
		return nil, fmt.Errorf("multiple virtual machines with UUID %q found", uuid)
	}

	vmRef := object.NewReference(client.Client, vms[0].Self)

	finder := find.NewFinder(client.Client, false)

	vm, err := finder.ObjectReference(ctx, vmRef.Reference())
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] VM %q found for UUID %q", vm.(*object.VirtualMachine).InventoryPath, uuid)

	return vm.(*object.VirtualMachine), nil
}
