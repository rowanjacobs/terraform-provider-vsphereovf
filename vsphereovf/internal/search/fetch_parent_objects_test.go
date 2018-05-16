package search_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/internal/search"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/internal/search/searchfakes"
	"github.com/vmware/govmomi/object"
)

var _ = Describe("FetchParentObjects", func() {
	var (
		finder     *searchfakes.FakeFinder
		networkMap map[string]interface{}

		datacenter   *object.Datacenter
		datastore    *object.Datastore
		folder       *object.Folder
		resourcePool *object.ResourcePool
		network1     object.NetworkReference
		network2     object.NetworkReference
	)

	BeforeEach(func() {
		finder = &searchfakes.FakeFinder{}

		networkMap = map[string]interface{}{
			"network-a": "network-1",
			"network-b": "network-2",
		}

		datacenter = &object.Datacenter{}
		datastore = &object.Datastore{}
		folder = &object.Folder{}
		resourcePool = &object.ResourcePool{}
		network1 = &object.Network{}
		network2 = &object.Network{}

		finder.DatacenterReturns(datacenter, nil)
		finder.DatastoreReturns(datastore, nil)
		finder.FolderReturns(folder, nil)
		finder.ResourcePoolReturns(resourcePool, nil)
		finder.NetworkReturnsOnCall(0, network1, nil)
		finder.NetworkReturnsOnCall(1, network2, nil)
	})

	It("delegates to finder", func() {
		tpo, err := search.FetchParentObjectsImpl(finder, "some-dc", "some-ds", "some-folder", "some-rp", networkMap)
		Expect(err).NotTo(HaveOccurred())

		Expect(finder.DatacenterArgsForCall(0)).To(Equal("some-dc"))
		Expect(finder.DatastoreArgsForCall(0)).To(Equal("some-ds"))
		Expect(finder.FolderArgsForCall(0)).To(Equal("some-folder"))
		Expect(finder.ResourcePoolArgsForCall(0)).To(Equal("some-rp"))

		Expect([]string{
			finder.NetworkArgsForCall(0),
			finder.NetworkArgsForCall(1),
		}).To(ConsistOf([]string{
			"network-1",
			"network-2",
		}))

		Expect(tpo).To(Equal(search.TemplateParentObjects{
			Datacenter:   datacenter,
			Datastore:    datastore,
			Folder:       folder,
			ResourcePool: resourcePool,
			Networks: map[string]object.NetworkReference{
				"network-a": network1,
				"network-b": network2,
			},
		}))
	})

	Describe("error cases", func() {
		Context("when finder fails to get datacenter", func() {
			BeforeEach(func() {
				finder.DatacenterReturns(nil, errors.New("papaya"))
			})

			It("returns an error", func() {
				_, err := search.FetchParentObjectsImpl(finder, "some-dc", "some-ds", "some-folder", "some-rp", networkMap)
				Expect(err).To(MatchError("error retrieving datacenter 'some-dc': papaya"))
			})
		})

		Context("when finder fails to get datastore", func() {
			BeforeEach(func() {
				finder.DatastoreReturns(nil, errors.New("mango"))
			})

			It("returns an error", func() {
				_, err := search.FetchParentObjectsImpl(finder, "some-dc", "some-ds", "some-folder", "some-rp", networkMap)
				Expect(err).To(MatchError("error retrieving datastore 'some-ds': mango"))
			})
		})

		Context("when finder fails to get folder", func() {
			BeforeEach(func() {
				finder.FolderReturns(nil, errors.New("lychee"))
			})

			It("returns an error", func() {
				_, err := search.FetchParentObjectsImpl(finder, "some-dc", "some-ds", "some-folder", "some-rp", networkMap)
				Expect(err).To(MatchError("error retrieving folder 'some-folder': lychee"))
			})
		})

		Context("when finder fails to get resource pool", func() {
			BeforeEach(func() {
				finder.ResourcePoolReturns(nil, errors.New("guava"))
			})

			It("returns an error", func() {
				_, err := search.FetchParentObjectsImpl(finder, "some-dc", "some-ds", "some-folder", "some-rp", networkMap)
				Expect(err).To(MatchError("error retrieving resource pool 'some-rp': guava"))
			})
		})

		Context("when finder fails to get network", func() {
			BeforeEach(func() {
				finder.NetworkReturnsOnCall(0, nil, errors.New("coconut"))
			})

			It("returns an error", func() {
				_, err := search.FetchParentObjectsImpl(finder, "some-dc", "some-ds", "some-folder", "some-rp", networkMap)
				Expect(err).To(MatchError(MatchRegexp(`error retrieving network 'network-\d': coconut`)))
			})
		})
	})
})
