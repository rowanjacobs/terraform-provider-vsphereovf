package vsphereovf_test

import (
	"context"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf"
	"github.com/vmware/govmomi"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var _ = Describe("Provider", func() {
	var (
		provider *schema.Provider
	)

	BeforeEach(func() {
		var ok bool
		provider, ok = vsphereovf.Provider().(*schema.Provider)
		Expect(ok).To(BeTrue())
	})

	AfterEach(func() {
		vsphereovf.ResetNewGovmomiClient()
	})

	It("returns a valid Provider", func() {
		err := provider.InternalValidate()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Provider().Configure", func() {
		var (
			resourceConfig        *terraform.ResourceConfig
			expectedGovmomiClient *govmomi.Client

			newGovmomiClientCallCount        int
			newGovmomiClientReceivedContext  context.Context
			newGovmomiClientReceivedURL      *url.URL
			newGovmomiClientReceivedInsecure bool
		)
		BeforeEach(func() {
			rawConfig, err := config.NewRawConfig(map[string]interface{}{
				"user":                 "some-user",
				"password":             "some-password",
				"vsphere_server":       "vsphere.example.com",
				"allow_unverified_ssl": false,
			})
			Expect(err).NotTo(HaveOccurred())
			resourceConfig = terraform.NewResourceConfig(rawConfig)

			expectedGovmomiClient = &govmomi.Client{}

			newGovmomiClientStub := func(ctx context.Context, u *url.URL, insecure bool) (*govmomi.Client, error) {
				newGovmomiClientCallCount++
				newGovmomiClientReceivedContext = ctx
				newGovmomiClientReceivedURL = u
				newGovmomiClientReceivedInsecure = insecure

				return expectedGovmomiClient, nil
			}
			vsphereovf.SetNewGovmomiClient(newGovmomiClientStub)
		})

		It("sets up a govmomi client to be returned by Provider().Meta()", func() {
			err := provider.Configure(resourceConfig)
			Expect(err).NotTo(HaveOccurred())

			client := provider.Meta()
			Expect(client).To(Equal(expectedGovmomiClient))

			Expect(newGovmomiClientCallCount).To(Equal(1))
			Expect(newGovmomiClientReceivedContext).To(Equal(context.Background()))
			Expect(newGovmomiClientReceivedInsecure).To(BeFalse())

			Expect(newGovmomiClientReceivedURL.Scheme).To(Equal("https"))
			Expect(newGovmomiClientReceivedURL.Host).To(Equal("vsphere.example.com"))
			Expect(newGovmomiClientReceivedURL.Path).To(Equal("/sdk"))

			userInfo := newGovmomiClientReceivedURL.User
			Expect(userInfo.Username()).To(Equal("some-user"))
			password, set := userInfo.Password()
			Expect(set).To(BeTrue())
			Expect(password).To(Equal("some-password"))
		})

		Context("when allow_unverified_ssl is set", func() {
			BeforeEach(func() {
				rawConfig, err := config.NewRawConfig(map[string]interface{}{
					// "user":                 "some-user",
					// "password":             "some-password",
					// "vsphere_server": "vsphere.example.com/sdk",
					"allow_unverified_ssl": true,
				})
				Expect(err).NotTo(HaveOccurred())
				resourceConfig = terraform.NewResourceConfig(rawConfig)
			})

			It("sets up a govmomi client that allows insecure connections", func() {
				err := provider.Configure(resourceConfig)
				Expect(err).NotTo(HaveOccurred())

				client := provider.Meta()
				Expect(client).To(Equal(expectedGovmomiClient))

				Expect(newGovmomiClientReceivedInsecure).To(BeTrue())
			})
		})

		Context("when host URL cannot be parsed", func() {
			BeforeEach(func() {
				rawConfig, err := config.NewRawConfig(map[string]interface{}{
					"vsphere_server": " ",
				})
				Expect(err).NotTo(HaveOccurred())
				resourceConfig = terraform.NewResourceConfig(rawConfig)
			})

			It("returns an error", func() {
				err := provider.Configure(resourceConfig)
				Expect(err).To(MatchError("error parsing vSphere server URL: parse https:// /sdk: invalid character \" \" in host name"))
			})
		})
	})
})
