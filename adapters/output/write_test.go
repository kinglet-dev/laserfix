package output_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserfix/adapters/output"
)

var _ = Describe("Write", func() {
	var dir string

	BeforeEach(func() {
		// Arrange (shared)
		dir = GinkgoT().TempDir()
	})

	It("writes a new file", func() {
		// Arrange
		name := filepath.Join(dir, "fixed.svg")

		// Act
		err := output.Write(name, []byte("<svg/>"))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(os.ReadFile(name)).To(Equal([]byte("<svg/>")))
	})

	It("refuses a file that already exists and leaves it as it was", func() {
		// Arrange
		name := filepath.Join(dir, "fixed.svg")
		Expect(os.WriteFile(name, []byte("mine"), 0o600)).To(Succeed())

		// Act
		err := output.Write(name, []byte("<svg/>"))

		// Assert
		Expect(err).To(MatchError(output.ErrExists))
		Expect(os.ReadFile(name)).To(Equal([]byte("mine")))
	})

	It("leaves nothing else in the folder", func() {
		// Act
		Expect(output.Write(filepath.Join(dir, "fixed.svg"), []byte("<svg/>"))).To(Succeed())

		// Assert
		entries, err := os.ReadDir(dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(entries).To(HaveLen(1))
	})

	It("passes on the error for a folder that doesn't exist", func() {
		// Act
		err := output.Write(filepath.Join(dir, "missing", "fixed.svg"), []byte("<svg/>"))

		// Assert
		Expect(err).To(MatchError(os.ErrNotExist))
		Expect(err.Error()).NotTo(ContainSubstring("missing"), "the caller shows the name")
	})

	It("lets others read the new file, as for any new document", func() {
		// Arrange
		name := filepath.Join(dir, "fixed.svg")

		// Act
		Expect(output.Write(name, []byte("<svg/>"))).To(Succeed())

		// Assert: Windows reports every readable file as readable by all.
		info, err := os.Stat(name)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Mode().Perm() & 0o044).To(Equal(os.FileMode(0o044)))
	})
})
