package search

import (
	"context"
	"fmt"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

type Finder struct {
	*find.Finder
}

const DefaultAPITimeout = 5 * time.Minute

func NewFinder(client *govmomi.Client, dcPath string) (Finder, error) {
	finder := Finder{find.NewFinder(client.Client, false)}
	dc, err := finder.Datacenter(dcPath)
	if err != nil {
		return Finder{}, fmt.Errorf("error retrieving datacenter '%s': %s", dcPath, err)
	}

	finder.SetDatacenter(dc)
	return finder, nil
}

// adapted from tf vsphere provider internals

func (f Finder) Datastore(path string) (*object.Datastore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.DatastoreOrDefault(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding datastore: %s", err)
	}

	return obj, nil
}

func (f Finder) Datacenter(path string) (*object.Datacenter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.Datacenter(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding datacenter: %s", err)
	}

	return obj, nil
}

func (f Finder) Folder(path string) (*object.Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.Folder(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding folder: %s", err)
	}

	return obj, nil
}

func (f Finder) ResourcePool(path string) (*object.ResourcePool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.ResourcePoolOrDefault(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("Finding resource pool: %s", err)
	}

	return obj, nil
}

func (f Finder) Network(networkPath string) (types.ManagedObjectReference, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
	defer cancel()

	obj, err := f.Finder.Network(ctx, networkPath)
	if err != nil {
		return types.ManagedObjectReference{}, fmt.Errorf("Finding network: %s", err)
	}

	return obj.Reference(), nil
}
