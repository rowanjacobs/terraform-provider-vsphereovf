package importer

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type Importer struct {
	ovfManager   OVFManager
	resourcePool ResourcePool
	datastore    mo.Reference
}

//go:generate counterfeiter . ResourcePool
type ResourcePool interface {
	ImportVApp(context.Context, types.BaseImportSpec, *object.Folder, *object.HostSystem) (Lease, error)
	Reference() types.ManagedObjectReference
}

//go:generate counterfeiter . Lease
type Lease interface {
	Upload(nfc.FileItem, string) error
	UploadAll([]types.OvfFileItem, string) error
}

type ResourcePoolImpl struct {
	*object.ResourcePool
}

func (r ResourcePoolImpl) ImportVApp(ctx context.Context, importSpec types.BaseImportSpec, folder *object.Folder, hostSystem *object.HostSystem) (Lease, error) {
	nfcLease, err := r.ResourcePool.ImportVApp(ctx, importSpec, folder, hostSystem)
	if err != nil {
		return nil, err
	}
	return lease.NewLease(lease.NFCLeaseImpl{nfcLease}), nil
}

//go:generate counterfeiter . OVFManager
type OVFManager interface {
	CreateImportSpec(context.Context, string, mo.Reference, mo.Reference, types.OvfCreateImportSpecParams) (*types.OvfCreateImportSpecResult, error)
}

// to keep ovfManager private, we have two separate constructors.
// this one takes an OVFManager, which can be faked.
func NewImporter(manager OVFManager, resourcePool ResourcePool, datastore mo.Reference) Importer {
	return Importer{
		ovfManager:   manager,
		resourcePool: resourcePool,
		datastore:    datastore,
	}
}

// to keep ovfManager private, we have two separate constructors.
// this one uses a library method ovf.NewManager to create an ovf.Manager from the govmomi.Client.
func NewImporterFromClient(client *govmomi.Client, resourcePool *object.ResourcePool, datastore mo.Reference) Importer {
	return NewImporter(
		ovf.NewManager(client.Client),
		ResourcePoolImpl{resourcePool},
		datastore,
	)
}

func (i Importer) CreateImportSpec(ovfContents string, networkRemap map[string]object.NetworkReference) (*types.OvfCreateImportSpecResult, error) {
	envelope, err := ovf.Unmarshal(bytes.NewBufferString(ovfContents))
	if err != nil {
		return nil, fmt.Errorf("invalid ovf: %s", err)
	}

	networkMapping := []types.OvfNetworkMapping{}
	for _, net := range envelope.Network.Networks {
		if networkRemap[net.Name] == nil {
			continue
		}

		networkMapping = append(networkMapping, types.OvfNetworkMapping{
			Name:    net.Name,
			Network: networkRemap[net.Name].Reference(),
		})
	}
	params := types.OvfCreateImportSpecParams{NetworkMapping: networkMapping}

	importSpec, err := i.ovfManager.CreateImportSpec(context.Background(), ovfContents, i.resourcePool.Reference(), i.datastore, params)
	if err != nil {
		return nil, err
	}
	if len(importSpec.Error) > 0 {
		return nil, fmt.Errorf("SOAP API fault: %s", importSpec.Error[0].LocalizedMessage)
	}
	return importSpec, nil
}

func (i Importer) Import(createImportSpecResult *types.OvfCreateImportSpecResult, folder *object.Folder, dir string) error {
	l, err := i.resourcePool.ImportVApp(context.TODO(), createImportSpecResult.ImportSpec, folder, nil)
	if err != nil {
		return err
	}

	return l.UploadAll(createImportSpecResult.FileItem, dir)
}
