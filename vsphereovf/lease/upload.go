package lease

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type Lease struct {
	ItemUpload
	NFCLease NFCLease
}

func NewLease(nfcLease NFCLease) Lease {
	return Lease{itemUploadImpl{nfcLease}, nfcLease}
}

//go:generate counterfeiter . NFCLease
type NFCLease interface {
	Upload(context.Context, nfc.FileItem, io.Reader, soap.Upload) error
	Wait(context.Context, []types.OvfFileItem) (*nfc.LeaseInfo, error)
	StartUpdater(context.Context, *nfc.LeaseInfo) *nfc.LeaseUpdater
	Complete(context.Context) error
}

func (l Lease) UploadAll(fileItems []types.OvfFileItem, dir string) error {
	ctx := context.Background()
	leaseInfo, err := l.NFCLease.Wait(ctx, fileItems)
	if err != nil {
		return err
	}

	updater := l.NFCLease.StartUpdater(ctx, leaseInfo)
	if updater != nil {
		defer updater.Done() // acceptance fails if this doesn't happen
	}

	for _, i := range leaseInfo.Items {
		err := l.Upload(i, filepath.Join(dir, i.Path))
		if err != nil {
			return err
		}
	}

	return l.NFCLease.Complete(ctx)
}

//go:generate counterfeiter . ItemUpload
type ItemUpload interface {
	Upload(nfc.FileItem, string) error
}

type itemUploadImpl struct {
	nfcLease NFCLease
}

// TODO: maybe this could be private? it needs an open lease...
func (i itemUploadImpl) Upload(item nfc.FileItem, fullPath string) error {
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	err = i.nfcLease.Upload(context.Background(), item, file, soap.Upload{ContentLength: fileInfo.Size()})
	if err != nil {
		return fmt.Errorf("Lease upload: %s", err)
	}

	return nil
}
