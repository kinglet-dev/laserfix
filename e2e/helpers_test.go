package e2e_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type buffer = bytes.Buffer

func stringsReader(s string) io.Reader { return strings.NewReader(s) }

// svg is a 20 × 20 mm document where one user unit is one millimeter.
func svg(body string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="20mm" height="20mm" viewBox="0 0 20 20">` + body + `</svg>`
}

// square is a 10 mm hairline square.
const square = `<rect x="2" y="2" width="10" height="10" fill="none" stroke="#000"/>`

// file writes content to a new temporary folder and returns its path.
func file(name, content string) string {
	path := filepath.Join(GinkgoT().TempDir(), name)
	Expect(os.WriteFile(path, []byte(content), 0o600)).To(Succeed())
	return path
}

// read returns a file's content.
func read(path string) string {
	b, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
