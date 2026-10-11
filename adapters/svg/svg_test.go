package svg_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/edit"
	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/laserfix/adapters/svg"
)

var _ = Describe("Document", func() {
	It("reads a document into a design in millimeters", func() {
		// Arrange
		src := doc(`<rect width="10" height="10" fill="none" stroke="#000"/>`)

		// Act
		d, err := svg.Document{}.Read(src)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Width).To(BeNumerically("~", 100, 1e-9))
		Expect(d.Shapes).To(HaveLen(1))
	})

	It("passes on the reader's error for a document it can't read", func() {
		// Act
		_, err := svg.Document{}.Read([]byte("<html/>"))

		// Assert
		Expect(err).To(MatchError(reader.ErrNotSVG))
	})

	It("removes the given elements and leaves every other byte", func() {
		// Arrange: elements are numbered root first, so the rects are 1, 2 and 3.
		src := doc(`<rect id="a"/><rect id="b"/><rect id="c"/>`)

		// Act
		out, err := svg.Document{}.Remove(src, []int{1, 3})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(string(doc(`<rect id="b"/>`))))
	})

	It("passes on the editor's error for an element the document doesn't have", func() {
		// Act
		_, err := svg.Document{}.Remove(doc(""), []int{7})

		// Assert
		Expect(err).To(MatchError(edit.ErrNoSuchElement))
	})
})
