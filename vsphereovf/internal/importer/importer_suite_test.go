package importer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestImportSpec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Import Spec Suite")
}
