package report_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("Describe words each finding", func() {
	one := []geom.Point{{X: 1, Y: 1}}
	two := []geom.Point{{X: 1, Y: 1}, {X: 2, Y: 2}}

	DescribeTable("with a label, what was found and how to fix it",
		func(f app.Finding, want report.Words) {
			// Act
			got := report.Describe(f, app.DefaultSettings)

			// Assert
			Expect(got).To(Equal(want))
		},
		Entry("one junction", app.Finding{Kind: app.Junction, Locations: one},
			report.Words{"Junction", "1 place where three or more ends meet, so which to join is unclear", "join them by hand"}),
		Entry("two junctions", app.Finding{Kind: app.Junction, Locations: two},
			report.Words{"Junction", "2 places where three or more ends meet, so which to join is unclear", "join them by hand"}),
		Entry("one gap", app.Finding{Kind: app.Gap, Locations: one, Size: 0.1},
			report.Words{"Gap", "1 gap up to 0.1 mm, wider than 0.05 mm (--join)", "close them by hand, or raise --join"}),
		Entry("line burned twice", app.Finding{Kind: app.Stacked, Locations: one, Size: 30},
			report.Words{"Overlap", "30.0 mm of line still burns twice, shared by different colors or filled shapes",
				"combine them by hand if they should burn once"}),
	)

	DescribeTable("Plural picks the word for how many there are",
		func(n int, want string) {
			// Act
			got := report.Plural(n, "gap", "gaps")

			// Assert
			Expect(got).To(Equal(want))
		},
		Entry("none", 0, "gaps"),
		Entry("one", 1, "gap"),
		Entry("two", 2, "gaps"),
	)
})
