package lease

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type Lease struct {
	nfcLease NFCLease
}

func NewLease(nfcLease NFCLease) Lease {
	return Lease{nfcLease: nfcLease}
}

//go:generate counterfeiter . NFCLease
type NFCLease interface {
	Upload(context.Context, nfc.FileItem, io.Reader, soap.Upload) error
	Wait(context.Context, []types.OvfFileItem) (*nfc.LeaseInfo, error)
}

//go:generate counterfeiter . ResourcePool
// type ResourcePool interface {
// 	ImportVApp(context.Context, types.BaseImportSpec, *object.Folder, *object.HostSystem) (NFCLease, error)
// }

func (l Lease) Upload(item nfc.FileItem, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	err = l.nfcLease.Upload(context.Background(), item, file, soap.Upload{ContentLength: fileInfo.Size()})
	if err != nil {
		return fmt.Errorf("Lease upload: %s", err)
	}

	return nil
}

// func (l Lease) Import(spec *types.OvfCreateImportSpecResult, folder *object.Folder) ([]nfc.FileItem, error) {
// 	return []nfc.FileItem{}, nil
// }
