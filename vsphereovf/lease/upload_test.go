package lease_test

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease/leasefakes"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx/ovxfakes"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/vim25/types"
)

var _ = Describe("Lease", func() {
	var (
		leece    lease.Lease
		nfcLease *leasefakes.FakeNFCLease
		rp       *ovxfakes.FakeReaderProvider
		buffer   *gbytes.Buffer

		filePath string
	)

	BeforeEach(func() {
		buffer = gbytes.NewBuffer()
		_, err := buffer.Write([]byte("some contents"))
		Expect(err).NotTo(HaveOccurred())

		rp = &ovxfakes.FakeReaderProvider{}
		rp.ReaderReturns(buffer, 42, nil)

		nfcLease = &leasefakes.FakeNFCLease{}
		leece = lease.NewLease(nfcLease, rp)

		tempDir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		filePath = filepath.Join(tempDir, "some-temp-file")
		err = ioutil.WriteFile(filePath, []byte("some contents"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Upload", func() {
		It("delegates to nfc.Lease", func() {
			item := nfc.FileItem{}

			err := leece.Upload(item)
			Expect(err).NotTo(HaveOccurred())

			Expect(nfcLease.UploadCallCount()).To(Equal(1))
			ctx, receivedItem, reader, upload := nfcLease.UploadArgsForCall(0)

			Expect(ctx).To(Equal(context.Background()))
			Expect(receivedItem).To(Equal(item))
			Expect(reader).To(gbytes.Say("some contents"))
			Expect(upload.ContentLength).To(BeEquivalentTo(42))
		})

		Context("when readerprovider can't find the file", func() {
			BeforeEach(func() {
				rp.ReaderReturns(gbytes.NewBuffer(), -1, errors.New("no such file"))
			})
			It("returns an error", func() {
				err := leece.Upload(nfc.FileItem{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the nfc.Lease upload fails", func() {
			BeforeEach(func() {
				nfcLease.UploadReturns(errors.New("coconut"))
			})

			It("returns an error", func() {
				err := leece.Upload(nfc.FileItem{})
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

			Expect(itemUpload.UploadArgsForCall(0).OvfFileItem).To(Equal(fileItems[0]))
			Expect(itemUpload.UploadArgsForCall(1).OvfFileItem).To(Equal(fileItems[1]))
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
