package svg_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSVG(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SVG Suite")
}

// doc is a 100 mm square document holding the given markup.
func doc(body string) []byte {
	return []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100mm" height="100mm" viewBox="0 0 100 100">` +
		body + `</svg>`)
}
