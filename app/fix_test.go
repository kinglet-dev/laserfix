package app_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("Fix removes shapes that copy an earlier shape", func() {
	var fixer *app.Fixer

	BeforeEach(func() {
		// Arrange (shared)
		fixer = app.NewFixer(svg.Document{}, svg.Document{})
	})

	It("removes a copy and lists the change", func() {
		// Arrange: the rects are elements 1, 2 and 3.
		src := doc(sq("10") + sq("40") + sq("10"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		c := r.Changes[0]
		Expect(c.Kind).To(Equal(app.Copy))
		Expect([]int{c.Element, c.Of}).To(Equal([]int{3, 1}))
		Expect(c.At).To(Equal(geom.Point{X: 10, Y: 10}))
		Expect(c.Length).To(BeNumerically("~", 40, 1e-9))
		Expect(string(r.Output)).To(Equal(string(doc(sq("10") + sq("40")))))
	})

	It("measures line burned twice before and after", func() {
		// Arrange
		src := doc(sq("10") + sq("10"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Before.Stacked).To(BeNumerically("~", 40, 1e-9))
		Expect(r.After.Stacked).To(BeZero())
	})

	It("counts crossings before and after", func() {
		// Arrange: a cross of two lines, drawn twice.
		cross := `<path d="M0 5H10M5 0V10" fill="none" stroke="#000"/>`
		src := doc(cross + cross)

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect([]int{r.Before.Crossings, r.After.Crossings}).To(Equal([]int{1, 1}))
	})

	It("returns the document unchanged when nothing needs fixing", func() {
		// Arrange
		src := doc(sq("10") + sq("40"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(BeEmpty())
		Expect(r.Output).To(Equal(src))
	})

	It("reports the design's size", func() {
		// Act
		r, err := fixer.Fix(doc(""), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect([]float64{r.Width, r.Height}).To(Equal([]float64{100, 100}))
	})

	It("reports the settings it used", func() {
		// Arrange
		s := app.DefaultSettings
		s.Same = 0.05

		// Act
		r, err := fixer.Fix(doc(""), s)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Settings).To(Equal(s))
	})

	It("finds a circle copied as a path of arcs, inside a group moved by nothing", func() {
		// Arrange
		circle := `<circle cx="50" cy="50" r="20" fill="none" stroke="#000"/>`
		arcs := `<g transform="translate(0 0)"><path d="M70 50 A20 20 0 0 1 30 50 A20 20 0 0 1 70 50Z" fill="none" stroke="#000"/></g>`

		// Act
		r, err := fixer.Fix(doc(circle+arcs), app.DefaultSettings)

		// Assert: the path is element 3, inside the group (2).
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect([]int{r.Changes[0].Element, r.Changes[0].Of}).To(Equal([]int{3, 1}))
	})
})
