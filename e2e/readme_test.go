package e2e_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// readmeExample is the output shown in the README after a command, to the
// end of its code block.
func readmeExample(command string) []string {
	b, err := os.ReadFile(filepath.Join("..", "README.md"))
	Expect(err).NotTo(HaveOccurred())
	text := strings.ReplaceAll(string(b), "\r\n", "\n") // Windows checkouts may use CRLF
	_, after, found := strings.Cut(text, "$ "+command+"\n")
	Expect(found).To(BeTrue(), "the README's example command")
	block, _, found := strings.Cut(after, "\n```")
	Expect(found).To(BeTrue(), "the end of the README's example")
	return strings.Split(block, "\n")
}

// lines splits output into lines, with the paths given replaced by the
// names the README shows and the version by the README's.
func lines(out string, names map[string]string) []string {
	for path, name := range names {
		out = strings.ReplaceAll(out, path, name)
	}
	out = strings.ReplaceAll(out, "laserfix 0.0.0-test", "laserfix 0.1.0")
	return strings.Split(strings.TrimRight(out, "\n"), "\n")
}

var _ = Describe("The README", func() {
	example := filepath.Join("..", "examples", "stacked.svg")

	It("shows the real preview for the example design", func() {
		// Act
		out, _, code := laserfix("", example)

		// Assert
		Expect(code).To(Equal(1))
		Expect(lines(out, map[string]string{example: "stacked.svg"})).To(Equal(readmeExample("laserfix stacked.svg")))
	})

	It("shows the real output for writing the repaired copy", func() {
		// Arrange
		fixed := filepath.Join(GinkgoT().TempDir(), "stacked-fixed.svg")

		// Act
		out, _, code := laserfix("", example, "-o", fixed)

		// Assert
		Expect(code).To(Equal(0))
		got := lines(out, map[string]string{example: "stacked.svg", fixed: "stacked-fixed.svg"})
		Expect(got).To(Equal(readmeExample("laserfix stacked.svg -o stacked-fixed.svg")))
	})
})
