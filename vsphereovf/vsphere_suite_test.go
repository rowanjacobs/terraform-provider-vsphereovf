package vsphereovf_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVSphereOVF(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vSphereOVF Suite")
}
