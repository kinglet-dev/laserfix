package e2e_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("The laserfix binary", func() {
	It("previews the changes, writing nothing, and exits 1", func() {
		// Arrange
		path := file("twice.svg", svg(square+square))

		// Act
		out, errOut, code := laserfix("", path)

		// Assert
		Expect(code).To(Equal(1))
		Expect(out).To(ContainSubstring("laserfix 0.0.0-test · "))
		Expect(out).To(ContainSubstring("Copies   1 shape, 40.0 mm of line burned twice → remove the later copy\n"))
		Expect(errOut).To(BeEmpty())
		Expect(read(path)).To(Equal(svg(square + square)))
	})

	It("writes the repaired copy and leaves the original as it was", func() {
		// Arrange
		path := file("twice.svg", svg(square+square))
		fixed := filepath.Join(filepath.Dir(path), "fixed.svg")

		// Act
		out, _, code := laserfix("", path, "-o", fixed)

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("Line burned twice   40.0 mm → 0.0 mm\n"))
		Expect(read(fixed)).To(Equal(svg(square)))
		Expect(read(path)).To(Equal(svg(square + square)))
	})

	It("refuses to replace an existing file, exiting 2", func() {
		// Arrange
		path := file("twice.svg", svg(square+square))

		// Act
		_, errOut, code := laserfix("", path, "-o", path)

		// Assert
		Expect(code).To(Equal(2))
		Expect(errOut).To(ContainSubstring("it already exists"))
		Expect(read(path)).To(Equal(svg(square + square)))
	})

	It("works in a pipeline from standard input to standard output", func() {
		// Act
		out, errOut, code := laserfix(svg(square+square), "-", "-o", "-")

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(Equal(svg(square)))
		Expect(errOut).To(ContainSubstring("standard input → standard output"))
	})

	It("shows its version", func() {
		// Act
		out, _, code := laserfix("", "--version")

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(Equal("laserfix 0.0.0-test\n"))
	})
})
