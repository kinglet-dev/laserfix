package copies_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

func TestCopies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Copies Suite")
}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// outline is a path through the points, closed or open.
func outline(closed bool, points ...geom.Point) geom.Path {
	s := geom.Subpath{Start: points[0], Closed: closed}
	for _, p := range points[1:] {
		s.Segments = append(s.Segments, geom.LineTo(p))
	}
	return geom.Path{Subpaths: []geom.Subpath{s}}
}

// hairline is a stroked shape with no fill, in the given stroke color.
func hairline(stroke string, closed bool, points ...geom.Point) design.Shape {
	return design.Shape{Path: outline(closed, points...), Fill: "none", Stroke: stroke}
}

// square is a 10 mm hairline square with its corner at (x, y).
func square(x, y float64) design.Shape {
	return hairline("#000000", true, pt(x, y), pt(x+10, y), pt(x+10, y+10), pt(x, y+10))
}
