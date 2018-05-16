package lease_test

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease/leasefakes"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/vim25/types"
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

	Describe("UploadAll", func() {
		var (
			fileItems  []types.OvfFileItem
			itemUpload *leasefakes.FakeItemUpload
		)
		BeforeEach(func() {
			itemUpload = &leasefakes.FakeItemUpload{}
			leece = lease.Lease{
				ItemUpload: itemUpload,
				NFCLease:   nfcLease,
			}
			fileItem1 := types.OvfFileItem{
				Path: "first-path",
				Size: 123,
			}
			fileItem2 := types.OvfFileItem{
				Path: "second-path",
				Size: 456,
			}
			fileItems = []types.OvfFileItem{fileItem1, fileItem2}
			leaseInfo := nfc.LeaseInfo{
				Items: []nfc.FileItem{
					{
						OvfFileItem: fileItem1,
					},
					{
						OvfFileItem: fileItem2,
					},
				},
			}

			nfcLease.WaitReturns(&leaseInfo, nil)
		})

		It("delegates to nfc.Lease", func() {
			err := leece.UploadAll(fileItems, "some-dir")
			Expect(err).NotTo(HaveOccurred())

			Expect(nfcLease.WaitCallCount()).To(Equal(1))
			_, actualFileItems := nfcLease.WaitArgsForCall(0)
			Expect(actualFileItems).To(Equal(fileItems))

			Expect(itemUpload.UploadCallCount()).To(Equal(2))

			item1, path1 := itemUpload.UploadArgsForCall(0)
			Expect(item1.OvfFileItem).To(Equal(fileItems[0]))
			Expect(path1).To(Equal("some-dir/first-path"))

			item2, path2 := itemUpload.UploadArgsForCall(1)
			Expect(item2.OvfFileItem).To(Equal(fileItems[1]))
			Expect(path2).To(Equal("some-dir/second-path"))
		})

		Context("when we fail to to wait on the lease", func() {
			BeforeEach(func() {
				nfcLease.WaitReturns(nil, errors.New("kiwi"))
			})

			It("errors", func() {
				err := leece.UploadAll(fileItems, "some-dir")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when we fail to upload an item", func() {
			BeforeEach(func() {
				itemUpload.UploadReturns(errors.New("kumquat"))
			})

			It("errors", func() {
				err := leece.UploadAll(fileItems, "some-dir")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
