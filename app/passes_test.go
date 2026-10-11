package app_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

// needsTwoPasses was found by a random search: trimming leaves pieces to
// join, and the joined path runs back along itself near a corner, which a
// second pass trims.
const needsTwoPasses = `<path d="M10 10 L20.06 37.673727663982255 L9.97 0.02" fill="none" stroke="#000"/>` +
	`<path d="M10 29.94 L24.740001574181882 30.06 L10 0 L33.85178298635877 22.866289168883092" fill="none" stroke="red"/>` +
	`<path d="M0.06 29.94 L0 20 L20.03 29.94" fill="none" stroke="#000"/>` +
	`<path d="M31.336375279450966 30.02 L23.952488206143187 0.03 L26.550940729816027 30.03 L9.97 0.03" fill="none" stroke="#000"/>`

var _ = Describe("Fix repeats its passes until nothing changes", func() {
	var fixer *app.Fixer

	BeforeEach(func() {
		// Arrange (shared)
		fixer = app.NewFixer(svg.Document{}, svg.Document{})
	})

	It("fixes in one run what one pass leaves for another", func() {
		// Act
		first, err := fixer.Fix(doc(needsTwoPasses), app.DefaultSettings)
		Expect(err).NotTo(HaveOccurred())
		second, err := fixer.Fix(first.Output, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(second.Changes).To(BeEmpty())
	})

	It("names every element changed by its number in the file given", func() {
		// Arrange: the input's shapes are elements 1 to 4.
		src := doc(needsTwoPasses)

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		for _, c := range r.Changes {
			Expect(c.Element).To(BeElementOf(1, 2, 3, 4))
		}
	})

	It("names changes in later passes by their elements' numbers in the file given", func() {
		// Arrange: a square and its copy come first, so once the copy (element 2)
		// is gone, every later element's number in the changed document is one less.
		src := doc(sq("50") + sq("50") + needsTwoPasses)

		// Act
		r, err := fixer.Fix(src, app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		for _, c := range r.Changes {
			Expect(c.Element).To(BeElementOf(2, 3, 4, 5, 6))
			if c.Element == 2 {
				Expect(c.Kind).To(Equal(app.Copy))
				Expect(c.Of).To(Equal(1))
			}
		}
	})

	It("reports what it leaves as it is once, as it is in the end", func() {
		// Arrange: a gap too wide to close, far from the rest.
		gap := hair("M60 60 H70") + hair("M70.5 60 H80")

		// Act
		r, err := fixer.Fix(doc(needsTwoPasses+gap), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		var gaps int
		for _, f := range r.Findings {
			if f.Kind == app.Gap {
				gaps++
			}
		}
		Expect(gaps).To(Equal(1))
	})

	It("doesn't touch a document that needs no changes", func() {
		// Arrange: an editor that fails shows it was never asked.
		fixer := app.NewFixer(svg.Document{}, editor{err: errors.New("asked")})

		// Act
		_, err := fixer.Fix(doc(sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
	})
})
