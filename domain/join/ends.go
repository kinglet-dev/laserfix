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
			ends = append(ends, end{n, false, sub.Start}, end{n, true, sub.Segments[len(sub.Segments)-1].To})
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

func dist(a, b geom.Point) float64 { return math.Hypot(b.X-a.X, b.Y-a.Y) }
