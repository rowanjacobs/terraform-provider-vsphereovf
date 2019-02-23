package search

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
)

const DefaultAPITimeout = 5 * time.Minute

type finder struct {
	*find.Finder
	c *vim25.Client
}

func NewFinder(client *govmomi.Client, dcPath string) (finder, error) {
	f := finder{
		Finder: find.NewFinder(client.Client, false),
		c:      client.Client,
	}
	dc, err := f.Datacenter(dcPath)
	if err != nil {
		return finder{}, DatacenterNotFoundError{dc: dcPath, message: err.Error()}
	}

	f.SetDatacenter(dc)
	return f, nil
}

// adapted from tf vsphere provider internals

func (f finder) Datastore(path string) (*object.Datastore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.DatastoreOrDefault(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding datastore: %s", err)
	}

	return obj, nil
}

func (f finder) Datacenter(path string) (*object.Datacenter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.Datacenter(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding datacenter: %s", err)
	}

	return obj, nil
}

func (f finder) Folder(path string) (*object.Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.Folder(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding folder '%s': %s", path, err)
	}

	return obj, nil
}

func (f finder) ResourcePool(path string) (*object.ResourcePool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.ResourcePoolOrDefault(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding resource pool: %s", err)
	}

	return obj, nil
}

func (f finder) Network(networkPath string) (object.NetworkReference, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.Network(ctx, networkPath)
	if err != nil {
		return nil, fmt.Errorf("Finding network: %s", err)
	}

	return obj, nil
}

// takes VM or template inventory path, which has the format:
// /$dcname/vm/$foldername/$vmname
// this path will be more complex if there are nested folders.
// TODO: you might be able to use $folder1/$folder2/.../$folderN as the folder name but I haven't tried it.
func (f finder) VirtualMachine(inventoryPath string) (*object.VirtualMachine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	si := object.NewSearchIndex(f.c)
	log.Printf("[DEBUG] searching for vm at %s\n", inventoryPath)
	ref, err := si.FindByInventoryPath(ctx, inventoryPath)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, fmt.Errorf("could not find VM at inventory path '%s'", inventoryPath)
	}
	return ref.(*object.VirtualMachine), nil
}
