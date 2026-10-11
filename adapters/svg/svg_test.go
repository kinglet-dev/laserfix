package svg_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/edit"
	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
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

	It("removes elements given no path, and leaves every other byte", func() {
		// Arrange: elements are numbered root first, so the rects are 1, 2 and 3.
		src := doc(`<rect id="a"/><rect id="b"/><rect id="c"/>`)

		// Act
		out, err := svg.Document{}.Edit(src, []app.Edit{{Element: 1}, {Element: 3}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(string(doc(`<rect id="b"/>`))))
	})

	It("passes on the editor's error for an element the document doesn't have", func() {
		// Act
		_, err := svg.Document{}.Edit(doc(""), []app.Edit{{Element: 7}})

		// Assert
		Expect(err).To(MatchError(edit.ErrNoSuchElement))
	})

	It("redraws an element with its new path, in its own units", func() {
		// Arrange
		src := doc(`<rect id="a" width="10" height="10" fill="none"/>`)
		path := geom.Path{Subpaths: []geom.Subpath{{Start: geom.Point{X: 0, Y: 0}, Segments: []geom.Segment{geom.LineTo(geom.Point{X: 10, Y: 0})}}}}

		// Act
		out, err := svg.Document{}.Edit(src, []app.Edit{{Element: 1, Path: path}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(string(doc(`<path id="a" fill="none" d="M0 0 L10 0"/>`))))
	})
})
