package importer_test

import (
	"bytes"
	"encoding/xml"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/internal/importer"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/internal/importer/importerfakes"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/types"
)

var _ = Describe("Importer", func() {
	var (
		importr      importer.Importer
		ovfManager   *importerfakes.FakeOVFManager
		resourcePool *importerfakes.FakeResourcePool
		lease        *importerfakes.FakeLease

		resourcePoolRef types.ManagedObjectReference
		datastore       types.ManagedObjectReference

		ovfContents  string
		networkRemap map[string]object.NetworkReference
	)

	BeforeEach(func() {
		resourcePoolRef = types.ManagedObjectReference{
			Type:  "resourcePool",
			Value: "some-resource-pool",
		}
		datastore = types.ManagedObjectReference{
			Type:  "datastore",
			Value: "some-datastore",
		}

		resourcePool = &importerfakes.FakeResourcePool{}
		ovfManager = &importerfakes.FakeOVFManager{}
		lease = &importerfakes.FakeLease{}
		importr = importer.NewImporter(ovfManager, resourcePool, datastore)

		envelope := &ovf.Envelope{
			Network: &ovf.NetworkSection{
				Networks: []ovf.Network{
					ovf.Network{
						Name:        "network-a",
						Description: "the first network",
					},
					ovf.Network{
						Name:        "network-b",
						Description: "the second network",
					},
					ovf.Network{
						Name:        "network-3",
						Description: "the third network",
					},
				},
			},
		}

		buf := bytes.NewBuffer([]byte{})
		err := xml.NewEncoder(buf).Encode(envelope)
		Expect(err).NotTo(HaveOccurred())
		ovfContents = buf.String()

		networkRemap = map[string]object.NetworkReference{
			"network-a": object.Network{
				object.NewCommon(nil, types.ManagedObjectReference{
					Value: "reference #1",
				}),
			},
			"network-b": object.Network{
				object.NewCommon(nil, types.ManagedObjectReference{
					Value: "reference #2",
				}),
			},
		}

		ovfManager.CreateImportSpecReturns(&types.OvfCreateImportSpecResult{}, nil)
		resourcePool.ReferenceReturns(resourcePoolRef)
		resourcePool.ImportVAppReturns(lease, nil)
	})

	Describe("CreateImportSpec", func() {
		It("creates an import spec", func() {
			_, err := importr.CreateImportSpec(ovfContents, networkRemap)
			Expect(err).NotTo(HaveOccurred())

			_, actualContents, actualResourcePool, actualDatastore, params := ovfManager.CreateImportSpecArgsForCall(0)
			Expect(actualContents).To(Equal(ovfContents))

			Expect(params.NetworkMapping).To(ConsistOf([]types.OvfNetworkMapping{
				{
					Name: "network-a",
					Network: types.ManagedObjectReference{
						Value: "reference #1",
					},
				},
				{
					Name: "network-b",
					Network: types.ManagedObjectReference{
						Value: "reference #2",
					},
				},
				// TODO: don't think we need this but we might
				// {
				// 	Name: "network-3",
				// 	Network: types.ManagedObjectReference{
				// 		Value: "reference #3",
				// 	},
				// },
			}))
			Expect(actualResourcePool).To(Equal(resourcePoolRef))
			Expect(actualDatastore).To(Equal(datastore))
		})

		Context("when the contents fail to unmarshal", func() {
			BeforeEach(func() {
				ovfContents = `some gtrash @@##<<<< that isnt valid xml`
			})

			It("fails", func() {
				_, err := importr.CreateImportSpec(ovfContents, networkRemap)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the ovf manager returns a soap fault", func() {
			BeforeEach(func() {
				resultWithFault := &types.OvfCreateImportSpecResult{
					Error: []types.LocalizedMethodFault{
						types.LocalizedMethodFault{
							LocalizedMessage: "coconut",
						},
					},
				}
				ovfManager.CreateImportSpecReturns(resultWithFault, nil)
			})

			It("returns an error", func() {
				_, err := importr.CreateImportSpec(ovfContents, networkRemap)
				Expect(err).To(MatchError("SOAP API fault: coconut"))
			})
		})
	})

	Describe("Import", func() {
		var (
			createImportSpecResult *types.OvfCreateImportSpecResult
			importSpec             types.BaseImportSpec
			folder                 *object.Folder
		)

		BeforeEach(func() {
			importSpec = &types.ImportSpec{
				InstantiationOst: &types.OvfConsumerOstNode{
					Id: "some-id",
				},
			}
			createImportSpecResult = &types.OvfCreateImportSpecResult{
				FileItem: []types.OvfFileItem{
					{
						Path: "some-ovf-path",
						Size: 123,
					},
					{
						Path: "some-vmdk-path",
						Size: 456,
					},
				},
				ImportSpec: importSpec,
			}
			folder = &object.Folder{
				Common: object.Common{
					InventoryPath: "some-folder-path",
				},
			}
		})

		It("imports the OVF file to vSphere", func() {
			err := importr.Import(createImportSpecResult, folder, "some-local-dir")
			Expect(err).NotTo(HaveOccurred())

			By("using the resource pool to get a lease", func() {
				Expect(resourcePool.ImportVAppCallCount()).To(Equal(1))
				_, actualImportSpec, actualFolder, _ := resourcePool.ImportVAppArgsForCall(0)
				Expect(actualImportSpec).To(Equal(importSpec))
				Expect(actualFolder).To(Equal(folder))
			})

			By("using the lease to upload the files", func() {
				Expect(lease.UploadAllCallCount()).To(Equal(1))
				actualFileItem, actualLocalDir := lease.UploadAllArgsForCall(0)
				Expect(actualFileItem).To(Equal(createImportSpecResult.FileItem))
				Expect(actualLocalDir).To(Equal("some-local-dir"))
			})
		})

		Context("when importing the vApp into the resource pool errors", func() {
			BeforeEach(func() {
				resourcePool.ImportVAppReturns(nil, errors.New("mango"))
			})

			It("returns the error", func() {
				err := importr.Import(createImportSpecResult, folder, "some-local-dir")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when uploading the files errors", func() {
			BeforeEach(func() {
				lease.UploadAllReturns(errors.New("coconut"))
			})

			It("returns the error", func() {
				err := importr.Import(createImportSpecResult, folder, "some-local-dir")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

/*
	// together, the following lines are a sort of blocking ImportVApp
	// use our import spec to get an nfc lease from our resource pool
	nfcLease, err := resourcePool.ImportVApp(ctx, importSpec.ImportSpec, folder, nil)
	// use our import spec to get lease info (including list of item URLs) out of the nfc lease
	// info, err := nfcLease.Wait(ctx, importSpec.FileItem)

	// // use the nfc lease to upload the ovf file
	// lease := lease.NewLease(nfcLease)
	// err = lease.Upload(info.Items[0], ovfPath)

	// new way:
	lease := lease.NewLease(nfcLease)
	err = lease.Import(importSpec.FileItem, ovfPath)
	// fileItem, err := lease.Wait(importSpec.FileItem)
	// err = lease.Upload(fileItem, ovfPath)
*/
