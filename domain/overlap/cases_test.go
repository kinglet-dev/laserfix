package overlap_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/overlap"
)

// find runs Find and expects no error.
func find(shapes []design.Shape, gone []bool) []overlap.Trim {
	trims, err := overlap.Find(shapes, gone, params)
	Expect(err).NotTo(HaveOccurred())
	return trims
}

// line is an open black hairline through the points.
func line(points ...geom.Point) design.Shape { return shape(poly(false, points...)) }

var _ = Describe("Find handles the cases from the spec", func() {
	It("trims the part of a line that overlaps an earlier line", func() {
		// Arrange
		shapes := []design.Shape{line(pt(10, 0), pt(110, 0)), line(pt(60, 0), pt(160, 0))}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Removed).To(BeNumerically("~", 50, 1e-6))
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(110, 0), pt(160, 0)}}))
	})

	It("removes all of a short line lying inside a longer one", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(100, 0)), line(pt(30, 0), pt(60, 0))}, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Keep).To(BeEmpty())
	})

	It("cuts a leftover bit shorter than the same-line distance too", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(0, 0), pt(10.01, 0))}, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Keep).To(BeEmpty())
	})

	It("trims line running beside earlier line within the same-line distance", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(100, 0)), line(pt(50, 0.01), pt(150, 0.01))}, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Removed).To(BeNumerically("~", 50, 1e-6))
	})

	It("leaves line further away than the same-line distance", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(100, 0)), line(pt(0, 0.3), pt(100, 0.3))}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("leaves a line that ends on another, as in a T", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(100, 0)), line(pt(50, 0), pt(50, 30))}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("leaves line crossing at a shallow angle over less than five times the same-line distance", func() {
		// Arrange: at 30°, the crossing line is within 0.02 mm of the other for 0.08 mm.
		far := 50 * math.Sqrt(3)

		// Act
		trims := find([]design.Shape{line(pt(-100, 0), pt(100, 0)), line(pt(-far, -50), pt(far, 50))}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("joins the stretch kept across the start of a closed shape", func() {
		// Arrange: the later square starts at its far corner, so the shared side is in the middle.
		later := shape(poly(true, pt(20, 0), pt(20, 10), pt(10, 10), pt(10, 0)))
		shapes := []design.Shape{square(0, 0), later}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(10, 0), pt(20, 0), pt(20, 10), pt(10, 10)}}))
	})
})

var _ = Describe("Find leaves alone what it mustn't change", func() {
	It("keeps a shape's other outlines whole when one is trimmed", func() {
		// Arrange: the later shape is two squares in one path; only the first touches.
		two := square(10, 0)
		two.Local.Subpaths = append(two.Local.Subpaths, square(50, 0).Local.Subpaths[0])
		two.Path = two.Local
		shapes := []design.Shape{square(0, 0), two}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{
			{pt(10, 0), pt(20, 0), pt(20, 10), pt(10, 10)},
			{pt(50, 0), pt(60, 0), pt(60, 10), pt(50, 10), pt(50, 0)},
		}))
	})

	It("doesn't trim a line passing near an earlier dot", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(0, 0)), line(pt(-5, 0.5), pt(5, 0.5))}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("leaves shapes of different strokes, which get different laser settings", func() {
		// Arrange
		red := square(10, 0)
		red.Stroke = "#ff0000"

		// Act
		trims := find([]design.Shape{square(0, 0), red}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("leaves filled shapes, and doesn't trim hairlines against them", func() {
		// Arrange
		filled := square(0, 0)
		filled.Fill = "#000000"
		otherFilled := square(10, 0)
		otherFilled.Fill = "#000000"

		// Act
		trims := find([]design.Shape{filled, otherFilled, square(10, 0)}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("skips shapes already gone, which neither trim nor are trimmed", func() {
		// Act
		trims := find([]design.Shape{square(0, 0), square(10, 0), square(0, 0)}, []bool{true, false, false})

		// Assert: with the first square gone, the second and third only share an edge,
		// which the third loses.
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Shape).To(Equal(2))
	})

	It("stops when the design has too much line to search", func() {
		// Arrange
		p := params
		p.MaxEntries = 5

		// Act
		_, err := overlap.Find([]design.Shape{line(pt(0, 0), pt(100, 0))}, nil, p)

		// Assert
		Expect(err).To(MatchError(overlap.ErrTooComplex))
	})

	It("stops when the line kept from several shapes goes over the limit", func() {
		// Arrange: each line passes through 3 grid cells, within the limit on its own.
		p := params
		p.MaxEntries = 5

		// Act
		_, err := overlap.Find([]design.Shape{line(pt(0.5, 0.5), pt(2.5, 0.5)), line(pt(0.5, 5.5), pt(2.5, 5.5))}, nil, p)

		// Assert
		Expect(err).To(MatchError(overlap.ErrTooComplex))
	})

	It("stops when one shape's own line goes over the limit, even if trimming would shorten it", func() {
		// Arrange: out through 3 grid cells and back beside the way out.
		p := params
		p.MaxEntries = 5

		// Act
		_, err := overlap.Find([]design.Shape{line(pt(0.5, 0.5), pt(2.5, 0.5), pt(0.5, 0.51))}, nil, p)

		// Assert
		Expect(err).To(MatchError(overlap.ErrTooComplex))
	})
})
