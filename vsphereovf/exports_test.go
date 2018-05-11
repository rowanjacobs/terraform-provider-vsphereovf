package vsphereovf

import (
	"context"
	"net/url"

	"github.com/vmware/govmomi"
)

func SetNewGovmomiClient(newClientFunc func(context.Context, *url.URL, bool) (*govmomi.Client, error)) {
	newGovmomiClient = newClientFunc
}

func ResetNewGovmomiClient() {
	newGovmomiClient = govmomi.NewClient
}
