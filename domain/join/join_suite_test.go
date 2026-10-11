package join_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/join"
)

func TestJoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Join Suite")
}

var params = join.Params{Join: 0.05, ReportGaps: 1}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// sub is a subpath through the points.
func sub(closed bool, points ...geom.Point) geom.Subpath {
	s := geom.Subpath{Start: points[0], Closed: closed}
	for _, p := range points[1:] {
		s.Segments = append(s.Segments, geom.LineTo(p))
	}
	return s
}

// shape is a black hairline of the subpaths, in millimeters as written.
func shape(subs ...geom.Subpath) design.Shape {
	p := geom.Path{Subpaths: subs}
	return design.Shape{Path: p, Local: p, ToMM: geom.Identity(), Fill: "none", Stroke: "#000000"}
}

// line is a shape of one open subpath through the points.
func line(points ...geom.Point) design.Shape { return shape(sub(false, points...)) }

// points flattens a subpath to its points, ending at the start when it is
// closed and doesn't already.
func points(s geom.Subpath) []geom.Point {
	l := geom.Path{Subpaths: []geom.Subpath{s}}.Flatten(0.001)[0]
	if n := len(l.Points); l.Closed && l.Points[n-1] != l.Points[0] {
		return append(l.Points, l.Points[0])
	}
	return l.Points
}
