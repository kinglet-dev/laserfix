package jsonreport_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJSONReport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSON Report Suite")
}
