package search

import (
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

type TemplateParentObjects struct {
	Datacenter   *object.Datacenter
	Datastore    *object.Datastore
	Folder       *object.Folder
	ResourcePool *object.ResourcePool
	Networks     map[string]object.NetworkReference
}

func FetchParentObjects(client *govmomi.Client, dcPath, dsPath, folderPath, rpPath string, networkMap map[string]interface{}) (TemplateParentObjects, error) {
	finder, err := NewFinder(client, dcPath)
	if err != nil {
		return TemplateParentObjects{}, err
	}
	return FetchParentObjectsImpl(finder, dcPath, dsPath, folderPath, rpPath, networkMap)
}

func FetchParentObjectsImpl(finder Finder, dcPath, dsPath, folderPath, rpPath string, networkMap map[string]interface{}) (TemplateParentObjects, error) {
	dc, err := finder.Datacenter(dcPath)
	if err != nil {
		return TemplateParentObjects{}, fmt.Errorf("error retrieving datacenter '%s': %s", dcPath, err)
	}

	ds, err := finder.Datastore(dsPath)
	if err != nil {
		return TemplateParentObjects{}, fmt.Errorf("error retrieving datastore '%s': %s", dsPath, err)
	}

	folder, err := finder.Folder(folderPath)
	if err != nil {
		return TemplateParentObjects{}, fmt.Errorf("error retrieving folder '%s': %s", folderPath, err)
	}

	rp, err := finder.ResourcePool(rpPath)
	if err != nil {
		return TemplateParentObjects{}, fmt.Errorf("error retrieving resource pool '%s': %s", rpPath, err)
	}

	networks := map[string]object.NetworkReference{}
	for templateKey, proposedNetworkName := range networkMap {
		network, err := finder.Network(proposedNetworkName.(string))
		if err != nil {
			return TemplateParentObjects{}, fmt.Errorf("error retrieving network '%s': %s", proposedNetworkName, err)
		}
		networks[templateKey] = network
	}

	return TemplateParentObjects{
		Datacenter:   dc,
		Datastore:    ds,
		Folder:       folder,
		ResourcePool: rp,
		Networks:     networks,
	}, nil
}

//go:generate counterfeiter . Finder
type Finder interface {
	Datacenter(string) (*object.Datacenter, error)
	Datastore(string) (*object.Datastore, error)
	Folder(string) (*object.Folder, error)
	Network(string) (object.NetworkReference, error)
	ResourcePool(string) (*object.ResourcePool, error)
}
