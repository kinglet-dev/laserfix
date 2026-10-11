package app_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Suite")
}

// doc is a 100 mm square document holding the given markup.
func doc(body string) []byte {
	return []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100mm" height="100mm" viewBox="0 0 100 100">` +
		body + `</svg>`)
}

// sq is a 10 mm black hairline square with its corner at x, 10.
func sq(x string) string {
	return `<rect x="` + x + `" y="10" width="10" height="10" fill="none" stroke="#000"/>`
}
