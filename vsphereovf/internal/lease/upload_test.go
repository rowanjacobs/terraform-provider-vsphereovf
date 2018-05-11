package lease_test

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphere/internal/lease"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphere/internal/lease/leasefakes"
	"github.com/vmware/govmomi/nfc"
)

var _ = Describe("Lease", func() {
	var (
		leece    lease.Lease
		nfcLease *leasefakes.FakeNFCLease

		filePath string
	)

	BeforeEach(func() {
		nfcLease = &leasefakes.FakeNFCLease{}
		leece = lease.NewLease(nfcLease)

		tempDir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		filePath = filepath.Join(tempDir, "some-temp-file")
		err = ioutil.WriteFile(filePath, []byte("some contents"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Upload", func() {
		It("delegates to nfc.Lease", func() {
			item := nfc.FileItem{}

			err := leece.Upload(item, filePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(nfcLease.UploadCallCount()).To(Equal(1))
			ctx, receivedItem, reader, upload := nfcLease.UploadArgsForCall(0)

			Expect(ctx).To(Equal(context.Background()))
			Expect(receivedItem).To(Equal(item))
			Expect(reader).To(BeAssignableToTypeOf(&os.File{}))
			Expect(reader.(*os.File).Name()).To(Equal(filePath))
			Expect(upload.ContentLength).To(Equal(int64(len([]byte("some contents")))))
		})

		Context("when the file does not exist", func() {
			It("returns an error", func() {
				err := leece.Upload(nfc.FileItem{}, "/this/path/is/not/real")
				Expect(err).To(MatchError("open /this/path/is/not/real: no such file or directory"))
			})
		})

		Context("when the nfc.Lease upload fails", func() {
			BeforeEach(func() {
				nfcLease.UploadReturns(errors.New("coconut"))
			})

			It("returns an error", func() {
				err := leece.Upload(nfc.FileItem{}, filePath)
				Expect(err).To(MatchError("Lease upload: coconut"))
			})
		})
	})

	// Describe("Import", func() {
	// 	var fileItems []nfc.FileItem
	// 	BeforeEach(func() {
	// 		fileItems = []nfc.FileItem{
	// 			nfc.FileItem{},
	// 			nfc.FileItem{},
	// 		}
	// 		leaseInfo := nfc.LeaseInfo{
	// 			Items: fileItems,
	// 		}

	// 		nfcLease.WaitReturns(&leaseInfo, nil)
	// 	})

	// 	It("delegates to nfc.Lease", func() {
	// 		actualItems, err := leece.Import(spec, nil)
	// 	})
	// })
})
