package ovx_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOvx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ovx Suite")
}
