package overlap_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

// arcs is a circle of radius r around (cx, cy) as four quarter arcs,
// starting on the right; half is the first two only (the lower half on screen).
func arcs(cx, cy, r float64, quarters int) geom.Path {
	sub := geom.Subpath{Start: pt(cx+r, cy), Closed: quarters == 4}
	at := sub.Start
	for q := range quarters {
		a := float64(q+1) * math.Pi / 2
		to := pt(cx+r*math.Cos(a), cy+r*math.Sin(a))
		if q == 3 {
			to = sub.Start // end exactly where it starts, as SVG arcs do
		}
		sub.Segments = append(sub.Segments, geom.ArcTo(at, r, r, 0, false, true, to)...)
		at = to
	}
	return geom.Path{Subpaths: []geom.Subpath{sub}}
}

// length sums a path's flattened length.
func length(p geom.Path) float64 {
	total := 0.0
	for _, pts := range lines(p) {
		for i := 1; i < len(pts); i++ {
			total += math.Hypot(pts[i].X-pts[i-1].X, pts[i].Y-pts[i-1].Y)
		}
	}
	return total
}

var _ = Describe("Find trims curves without flattening them", func() {
	It("leaves the half of a circle that a half circle drawn earlier doesn't cover, as curves", func() {
		// Arrange
		shapes := []design.Shape{shape(arcs(50, 50, 20, 2)), shape(arcs(50, 50, 20, 4))}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		left := kept(shapes, trims, 1)
		Expect(left.Subpaths).To(HaveLen(1))
		for _, s := range left.Subpaths[0].Segments {
			Expect(s.Curved).To(BeTrue())
		}
		Expect(length(left)).To(BeNumerically("~", 20*math.Pi, 0.05))
		Expect(trims[0].Removed).To(BeNumerically("~", 20*math.Pi, 0.05))
	})

	It("cuts a shape in its own units", func() {
		// Arrange: the later line is drawn in half millimeters.
		later := line(pt(120, 0), pt(320, 0))
		later.ToMM = geom.Matrix{A: 0.5, D: 0.5}
		later.Path = later.Local.Transform(later.ToMM)
		shapes := []design.Shape{line(pt(0, 0), pt(100, 0)), later}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(200, 0), pt(320, 0)}}))
	})
})
