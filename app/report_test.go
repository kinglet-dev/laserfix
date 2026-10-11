package app_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("Report.Worse says whether the changes made anything worse", func() {
	DescribeTable("compares after with before",
		func(before, after app.Measure, worse bool) {
			// Arrange
			r := app.Report{Before: before, After: after}

			// Act
			got := r.Worse()

			// Assert
			Expect(got).To(Equal(worse))
		},
		Entry("not when both went down", app.Measure{Stacked: 40, Crossings: 2}, app.Measure{Stacked: 0, Crossings: 1}, false),
		Entry("not when nothing changed", app.Measure{Stacked: 40, Crossings: 2}, app.Measure{Stacked: 40, Crossings: 2}, false),
		Entry("when more line burns twice", app.Measure{Stacked: 40}, app.Measure{Stacked: 40.1}, true),
		Entry("when there are more crossings", app.Measure{Crossings: 2}, app.Measure{Crossings: 3}, true),
		Entry("not for a rise as shown, to 0.1 mm, of nothing", app.Measure{Stacked: 40}, app.Measure{Stacked: 40.04}, false),
	)
})
