package lease_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lease Suite")
}
