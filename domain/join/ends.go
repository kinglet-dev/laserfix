package join

import (
	"math"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

// end is one end of a piece: its start, or its last point.
type end struct {
	piece int // an index into the pieces
	last  bool
	at    geom.Point // mm
	out   geom.Point // the way the line runs out of the piece there
}

// candidates lists the pieces that may be joined, and their ends in order:
// the open subpaths of hairlines not skipped, whose units can be converted.
func candidates(shapes []design.Shape, skip []bool) ([]Piece, []end) {
	var pieces []Piece
	var ends []end
	for i, s := range shapes {
		if _, ok := s.ToMM.Inverse(); !ok || s.Fill != "none" || (i < len(skip) && skip[i]) {
			continue
		}
		for j, sub := range s.Path.Subpaths {
			if sub.Closed || len(sub.Segments) == 0 {
				continue
			}
			n := len(pieces)
			pieces = append(pieces, Piece{i, j})
			back := sub.Reverse()
			ends = append(ends, end{n, false, sub.Start, outward(back)}, end{n, true, back.Start, outward(sub)})
		}
	}
	return pieces, ends
}

// near lists, for each end, the other ends of the same stroke within r,
// found with a uniform grid of cells r wide.
func near(ends []end, stroke func(end) string, r float64) [][]int {
	cell := func(p geom.Point) [2]int64 { return [2]int64{int64(math.Floor(p.X / r)), int64(math.Floor(p.Y / r))} }
	grid := map[[2]int64][]int{}
	for i, e := range ends {
		k := cell(e.at)
		grid[k] = append(grid[k], i)
	}
	out := make([][]int, len(ends))
	for i, e := range ends {
		k := cell(e.at)
		for dx := int64(-1); dx <= 1; dx++ {
			for dy := int64(-1); dy <= 1; dy++ {
				for _, o := range grid[[2]int64{k[0] + dx, k[1] + dy}] {
					if o != i && stroke(ends[o]) == stroke(e) && dist(ends[o].at, e.at) <= r {
						out[i] = append(out[i], o)
					}
				}
			}
		}
	}
	return out
}

// outward is the way a subpath runs out of its last point: from the
// nearest distinct point before it (a curve's control point, or the
// segment's start).
func outward(s geom.Subpath) geom.Point {
	n := len(s.Segments)
	last, from := s.Segments[n-1], s.Start
	if n > 1 {
		from = s.Segments[n-2].To
	}
	if last.Curved {
		switch {
		case last.C2 != last.To:
			from = last.C2
		case last.C1 != last.To:
			from = last.C1
		}
	}
	return geom.Point{X: last.To.X - from.X, Y: last.To.Y - from.Y}
}

// facing reports whether the line runs out of end e toward point p,
// within 45°: two ends a gap apart that face each other look like one
// line broken, not two lines that end near each other.
func facing(e end, p geom.Point) bool {
	vx, vy := p.X-e.at.X, p.Y-e.at.Y
	return e.out.X*vx+e.out.Y*vy >= math.Sqrt2/2*math.Hypot(e.out.X, e.out.Y)*math.Hypot(vx, vy)
}

func dist(a, b geom.Point) float64 { return math.Hypot(b.X-a.X, b.Y-a.Y) }
