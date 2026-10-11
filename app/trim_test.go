package app_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("Fix trims line lying on earlier line", func() {
	var fixer *app.Fixer

	BeforeEach(func() {
		// Arrange (shared)
		fixer = app.NewFixer(svg.Document{}, svg.Document{})
	})

	It("trims the side two squares share from the later square", func() {
		// Arrange: sq puts squares at y = 10; the second rect is element 2.
		src := doc(sq("10") + sq("20"))

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		c := r.Changes[0]
		Expect(c.Kind).To(Equal(app.Trim))
		Expect(c.Element).To(Equal(2))
		Expect(c.Length).To(BeNumerically("~", 10, 1e-6))
		Expect(c.At).To(Equal(geom.Point{X: 20, Y: 20}))
		Expect(r.Before.Stacked).To(BeNumerically("~", 10, 1e-6))
		Expect(r.After.Stacked).To(BeZero())
		Expect(string(r.Output)).To(ContainSubstring(`<path fill="none" stroke="#000" d="M20 10 L30 10 L30 20 L20 20"/>`))
	})

	It("removes a shape left with no line", func() {
		// Arrange
		long := `<path d="M0 50 H100" fill="none" stroke="#000"/>`
		short := `<path d="M20 50 H60" fill="none" stroke="#000"/>`

		// Act
		r, err := fixer.Fix(doc(long+short), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect(r.Changes[0].Kind).To(Equal(app.Trim))
		Expect(string(r.Output)).To(Equal(string(doc(long))))
	})

	It("trims nothing from a copy it removes", func() {
		// Act
		r, err := fixer.Fix(doc(sq("10")+sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(HaveLen(1))
		Expect(r.Changes[0].Kind).To(Equal(app.Copy))
	})
})
