package lease

import (
	"context"
	"fmt"
	"io"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type Lease struct {
	ItemUpload
	NFCLease NFCLease
}

func NewLease(nfcLease NFCLease, rp ovx.ReaderProvider) Lease {
	return Lease{itemUploadImpl{nfcLease, rp}, nfcLease}
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

	// items are OVA contents: OVF and VMDK files
	for _, i := range leaseInfo.Items {
		err := l.Upload(i)
		if err != nil {
			return err
		}
	}

	return l.NFCLease.Complete(ctx)
}

//go:generate counterfeiter . ItemUpload
type ItemUpload interface {
	Upload(nfc.FileItem) error
}

type itemUploadImpl struct {
	nfcLease       NFCLease
	readerProvider ovx.ReaderProvider
}

// TODO: maybe this could be private? it needs an open lease...
func (i itemUploadImpl) Upload(item nfc.FileItem) error {
	r, size, err := i.readerProvider.Reader(item.Path)
	if err != nil {
		return err
	}

	err = i.nfcLease.Upload(context.Background(), item, r, soap.Upload{ContentLength: size})
	if err != nil {
		return fmt.Errorf("Lease upload: %s", err)
	}

	return nil
}
