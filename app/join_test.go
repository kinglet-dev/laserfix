package app_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

// hair is a black hairline path with path data d.
func hair(d string) string { return `<path d="` + d + `" fill="none" stroke="#000"/>` }

var _ = Describe("Fix joins outlines split into pieces", func() {
	var fixer *app.Fixer

	BeforeEach(func() {
		// Arrange (shared)
		fixer = app.NewFixer(svg.Document{}, svg.Document{})
	})

	It("joins pieces into the first and removes the others", func() {
		// Arrange: a square in four sides; the paths are elements 1 to 4.
		src := doc(hair("M10 10 H20") + hair("M20 10 V20") + hair("M20 20 H10") + hair("M10 20 V10"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(Equal([]app.Change{{Kind: app.Join, Element: 1, At: geom.Point{X: 10, Y: 10}, Pieces: 4}}))
		Expect(string(r.Output)).To(Equal(string(doc(hair("M10 10 L20 10 L20 20 L10 20 L10 10 Z")))))
	})

	It("joins pieces in another element's units and keeps that element's other outlines", func() {
		// Arrange: the second element is drawn at twice the scale and holds another square too.
		other := `<path d="M10 5 V10 M50 50 H60 V60 H50 Z" transform="scale(2)" fill="none" stroke="#000"/>`
		src := doc(hair("M0 0 H20 V10") + other)

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect(string(r.Output)).To(Equal(string(doc(hair("M0 0 L20 0 L20 10 L20 20") +
			`<path d="M50 50 L60 50 L60 60 L50 60 Z" transform="scale(2)" fill="none" stroke="#000"/>`))))
	})

	It("doesn't join pieces of shapes it trimmed", func() {
		// Arrange: the second square loses its shared side and is left open.
		src := doc(sq("10") + sq("20") + hair("M20 20 V30"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect(r.Changes[0].Kind).To(Equal(app.Trim))
	})

	It("joins two outlines of one element, each with another piece", func() {
		// Arrange
		src := doc(hair("M0 0 H10 M50 0 H60") + hair("M10 0 V10") + hair("M60 0 V10"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(2))
		Expect(string(r.Output)).To(Equal(string(doc(hair("M0 0 L10 0 L10 10 M50 0 L60 0 L60 10")))))
	})

	It("says where a join is in millimeters", func() {
		// Arrange
		scaled := `<path d="M5 5 H10" transform="scale(2)" fill="none" stroke="#000"/>`

		// Act
		r, err := fixer.Fix(doc(scaled+hair("M20 10 V30")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes[0].At).To(Equal(geom.Point{X: 10, Y: 10}))
	})

	It("doesn't join an open line it trimmed to another piece meeting its end", func() {
		// Act
		r, err := fixer.Fix(doc(hair("M0 0 H10")+hair("M5 0 H20")+hair("M20 0 V10")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect(r.Changes[0].Kind).To(Equal(app.Trim))
	})

	It("closes an outline left open by a small gap", func() {
		// Act
		r, err := fixer.Fix(doc(hair("M10 10 H20 V20 H10 V10.02")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect(r.Changes[0].Kind).To(Equal(app.Close))
		Expect(string(r.Output)).To(Equal(string(doc(hair("M10 10 L20 10 L20 20 L10 20 L10 10.02 Z")))))
	})
})
