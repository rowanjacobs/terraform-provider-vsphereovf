package search

import (
	"context"
	"fmt"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

const DefaultAPITimeout = 5 * time.Minute

type finder struct {
	*find.Finder
}

func NewFinder(client *govmomi.Client, dcPath string) (finder, error) {
	f := finder{find.NewFinder(client.Client, false)}
	dc, err := f.Datacenter(dcPath)
	if err != nil {
		return finder{}, fmt.Errorf("error retrieving datacenter '%s': %s", dcPath, err)
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
		return nil, fmt.Errorf("Finding folder: %s", err)
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
