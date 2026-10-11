package text_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/adapters/report/text"
	"github.com/kinglet-dev/laserfix/app"
)

// write formats a report and returns the text.
func write(a report.About, r app.Report) string {
	var b strings.Builder
	Expect(text.New().Write(&b, a, r)).To(Succeed())
	return b.String()
}

var _ = Describe("The text report", func() {
	var r app.Report

	BeforeEach(func() {
		// Arrange (shared)
		r = app.Report{Width: 100, Height: 80, Before: app.Measure{Stacked: 80, Crossings: 1}, After: app.Measure{Crossings: 1}}
	})

	It("shows what would change in a preview", func() {
		// Arrange
		r.Changes = []app.Change{copyAt(10, 10), copyAt(40, 10.04)}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg"}, r)

		// Assert
		Expect(got).To(Equal(`laserfix 0.1.0 · grid.svg · 100.0 × 80.0 mm

Would change
  Copies   2 shapes, 80.0 mm of line burned twice → remove the later copy
           at (10.0, 10.0), (40.0, 10.0)

Run again with -o FILE to write the repaired copy.
`))
	})

	It("lists each kind of change on its own line, in the order they first appear", func() {
		// Arrange
		trim := app.Change{Kind: app.Trim, Element: 3, At: geom.Point{X: 20, Y: 20}, Length: 10}
		r.Changes = []app.Change{copyAt(10, 10), trim}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg"}, r)

		// Assert
		Expect(got).To(ContainSubstring(`
  Copies   1 shape, 40.0 mm of line burned twice → remove the later copy
           at (10.0, 10.0)
  Shared   1 shape, 10.0 mm of line burned twice → trim it from the later shape
           at (20.0, 20.0)
`))
	})

	It("counts outlines and pieces for joins", func() {
		// Arrange
		r.Changes = []app.Change{
			{Kind: app.Join, At: geom.Point{X: 1, Y: 2}, Pieces: 4},
			{Kind: app.Join, At: geom.Point{X: 3, Y: 4}, Pieces: 8},
		}

		// Act
		preview := write(report.About{Version: "0.1.0", Input: "a.svg"}, r)
		written := write(report.About{Version: "0.1.0", Input: "a.svg", Output: "b.svg"}, r)

		// Assert
		Expect(preview).To(ContainSubstring("\n  Split    2 outlines in 12 pieces → join each into one path\n           at (1.0, 2.0), (3.0, 4.0)\n"))
		Expect(written).To(ContainSubstring("\n  Split    2 outlines joined from 12 pieces\n"))
	})

	It("lists what it leaves as it is, with a fix for each", func() {
		// Arrange
		r.Settings = app.DefaultSettings
		r.Findings = []app.Finding{
			{Kind: app.Junction, Locations: []geom.Point{{X: 50, Y: 50}}},
			{Kind: app.Gap, Locations: []geom.Point{{X: 1, Y: 1}, {X: 2, Y: 2}}, Size: 0.6},
			{Kind: app.Stacked, Locations: []geom.Point{{X: 5, Y: 5}}, Size: 40},
		}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "a.svg"}, r)

		// Assert
		Expect(got).To(HaveSuffix(`Nothing to fix.

Left as they are
  Junction 1 place where three or more ends meet, so which to join is unclear → join them by hand
           at (50.0, 50.0)
  Gap      2 gaps up to 0.6 mm, wider than 0.05 mm (--join) → close them by hand, or raise --join
           at (1.0, 1.0), (2.0, 2.0)
  Overlap  40.0 mm of line still burns twice, shared by different colors or filled shapes → combine them by hand if they should burn once
           at (5.0, 5.0)
`))
	})

	It("says when there is nothing to fix", func() {
		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg"}, app.Report{Width: 100, Height: 80})

		// Assert
		Expect(got).To(Equal("laserfix 0.1.0 · grid.svg · 100.0 × 80.0 mm\n\nNothing to fix.\n"))
	})

	It("shows what changed and measures the copy again when writing", func() {
		// Arrange
		r.Changes = []app.Change{copyAt(10, 10)}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg", Output: "fixed.svg"}, r)

		// Assert
		Expect(got).To(Equal(`laserfix 0.1.0 · grid.svg → fixed.svg · 100.0 × 80.0 mm

Changed
  Copies   1 shape removed, 40.0 mm of line burned twice
           at (10.0, 10.0)

Checked again
  Line burned twice   80.0 mm → 0.0 mm
  Crossings           1 → 1

Wrote fixed.svg. Check it with: laserlint fixed.svg
`))
	})

	It("says when the copy was written unchanged", func() {
		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg", Output: "fixed.svg"}, app.Report{Width: 100, Height: 80})

		// Assert
		Expect(got).To(HaveSuffix("Nothing to fix.\n\nWrote fixed.svg, unchanged.\n"))
	})

	It("doesn't suggest checking standard output with laserlint by name", func() {
		// Arrange
		r.Changes = []app.Change{copyAt(10, 10)}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "standard input", Output: "standard output"}, r)

		// Assert
		Expect(got).To(HaveSuffix("\nWrote the repaired SVG to standard output.\n"))
	})

	It("lists at most five places, then how many more", func() {
		// Arrange
		for i := range 7 {
			r.Changes = append(r.Changes, copyAt(float64(i), 0))
		}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg"}, r)

		// Assert
		Expect(got).To(ContainSubstring("at (0.0, 0.0), (1.0, 0.0), (2.0, 0.0), (3.0, 0.0), (4.0, 0.0) and 2 more\n"))
	})

	It("lists exactly five places without saying more", func() {
		// Arrange
		for i := range 5 {
			r.Changes = append(r.Changes, copyAt(float64(i), 0))
		}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg"}, r)

		// Assert
		Expect(got).To(ContainSubstring("(4.0, 0.0)\n"))
	})

	It("warns when the changes made the design worse", func() {
		// Arrange
		r.Changes = []app.Change{copyAt(10, 10)}
		r.After = app.Measure{Stacked: 90, Crossings: 1}

		// Act
		got := write(report.About{Version: "0.1.0", Input: "grid.svg", Output: "fixed.svg"}, r)

		// Assert
		Expect(got).To(ContainSubstring("\nWorse after the changes: check fixed.svg before you burn it, and please report this at https://github.com/kinglet-dev/laserfix/issues\n"))
	})

	It("passes on an error writing the report", func() {
		// Act
		err := text.New().Write(failingWriter{}, report.About{}, r)

		// Assert
		Expect(err).To(MatchError(errBroken))
	})
})

var errBroken = errors.New("broken pipe")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errBroken }
