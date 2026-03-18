package nri_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNRI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NRI Suite")
}
