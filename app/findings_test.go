package app_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("Fix reports what it leaves as it is", func() {
	var fixer *app.Fixer

	BeforeEach(func() {
		// Arrange (shared)
		fixer = app.NewFixer(svg.Document{}, svg.Document{})
	})

	It("reports three ends meeting at a point", func() {
		// Act
		r, err := fixer.Fix(doc(hair("M10 50 H50")+hair("M50 50 H90")+hair("M50 50 V90")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(BeEmpty())
		Expect(r.Findings).To(Equal([]app.Finding{{Kind: app.Junction, Locations: []geom.Point{{X: 50, Y: 50}}}}))
	})

	It("reports gaps too wide to close, with the widest", func() {
		// Act
		r, err := fixer.Fix(doc(hair("M10 50 H50")+hair("M50.6 50 H90")+hair("M10 60 H50")+hair("M50.2 60 H90")), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Findings).To(HaveLen(1))
		f := r.Findings[0]
		Expect(f.Kind).To(Equal(app.Gap))
		Expect(f.Locations).To(Equal([]geom.Point{{X: 50, Y: 50}, {X: 50, Y: 60}}))
		Expect(f.Size).To(BeNumerically("~", 0.6, 1e-9))
	})

	It("reports line that still burns twice, shared by different colors", func() {
		// Arrange
		red := `<rect x="10" y="10" width="10" height="10" fill="none" stroke="red"/>`

		// Act
		r, err := fixer.Fix(doc(sq("10")+red), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Changes).To(BeEmpty())
		Expect(r.Findings).To(HaveLen(1))
		Expect(r.Findings[0].Kind).To(Equal(app.Stacked))
		Expect(r.Findings[0].Size).To(BeNumerically("~", 40, 1e-6))
	})

	It("doesn't report line burned twice too short to show, to 0.1 mm", func() {
		// Act
		r, err := fixer.Fix(doc(hair("M0 0 H10")+`<path d="M9.97 0 H20" fill="none" stroke="red"/>`), app.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Findings).To(BeEmpty())
	})
})
