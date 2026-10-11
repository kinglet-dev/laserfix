package overlap_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/overlap"
)

func TestOverlap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Overlap Suite")
}

var params = overlap.Params{Same: 0.02, Flatten: 0.005, MaxEntries: 10_000_000}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// poly is a path through the points, closed or open.
func poly(closed bool, points ...geom.Point) geom.Path {
	s := geom.Subpath{Start: points[0], Closed: closed}
	for _, p := range points[1:] {
		s.Segments = append(s.Segments, geom.LineTo(p))
	}
	return geom.Path{Subpaths: []geom.Subpath{s}}
}

// shape is a black hairline drawing path, in millimeters as written.
func shape(path geom.Path) design.Shape {
	return design.Shape{Path: path, Local: path, ToMM: geom.Identity(), Fill: "none", Stroke: "#000000"}
}

// square is a 10 mm hairline square with its corner at (x, y), drawn clockwise from it.
func square(x, y float64) design.Shape {
	return shape(poly(true, pt(x, y), pt(x+10, y), pt(x+10, y+10), pt(x, y+10)))
}

// cut returns what remains of shape i after a trim, as path data points.
func kept(shapes []design.Shape, trims []overlap.Trim, i int) geom.Path {
	for _, t := range trims {
		if t.Shape == i {
			return overlap.Cut(shapes[i].Local, t.Keep)
		}
	}
	return shapes[i].Local
}

// lines flattens a path into its points per subpath, for comparing.
func lines(p geom.Path) [][]geom.Point {
	var out [][]geom.Point
	for _, l := range p.Flatten(0.001) {
		pts := l.Points
		if l.Closed {
			pts = append(pts, pts[0])
		}
		out = append(out, pts)
	}
	return out
}

// approx matches points within a millionth of a millimeter.
func approx(want [][]geom.Point) OmegaMatcher {
	var subs []any
	for _, w := range want {
		var ps []any
		for _, p := range w {
			ps = append(ps, And(
				WithTransform(func(q geom.Point) float64 { return q.X }, BeNumerically("~", p.X, 1e-6)),
				WithTransform(func(q geom.Point) float64 { return q.Y }, BeNumerically("~", p.Y, 1e-6)),
			))
		}
		subs = append(subs, HaveExactElements(ps...))
	}
	return HaveExactElements(subs...)
}
