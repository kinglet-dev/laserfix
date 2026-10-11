package overlap

import (
	"math"

	"github.com/kinglet-dev/core/domain/geom"
)

// maxSin is the sine of the largest angle, 20°, between two lines that run
// the same way: wide enough for curves flattened at different points.
const maxSin = 0.342

// seg is a straight piece of line already kept.
type seg struct{ a, b geom.Point }

// index finds kept line near a point: a uniform grid whose cells list the
// pieces passing through them, found with Amanatides and Woo's voxel
// traversal ("A Fast Voxel Traversal Algorithm for Ray Tracing", 1987).
// Cells are at least Same wide, so line within Same of a point passes
// through the point's cell or one next to it.
type index struct {
	cell     float64
	cells    map[[2]int][]seg
	entries  int
	maxEntry int
}

func newIndex(cell float64, maxEntries int) *index {
	return &index{cell: cell, cells: map[[2]int][]seg{}, maxEntry: maxEntries}
}

func (x *index) key(v float64) int { return int(math.Floor(v / x.cell)) }

// add files a piece in every cell it passes through, and reports false
// once the index holds more entries than its limit.
func (x *index) add(a, b geom.Point) bool {
	ix, iy, jx, jy := x.key(a.X), x.key(a.Y), x.key(b.X), x.key(b.Y)
	stepX, tMaxX, tDeltaX := x.axis(a.X, b.X, ix)
	stepY, tMaxY, tDeltaY := x.axis(a.Y, b.Y, iy)
	steps := abs(jx-ix) + abs(jy-iy)
	x.entries += steps + 1
	if x.entries > x.maxEntry {
		return false
	}
	s := seg{a, b}
	x.cells[[2]int{ix, iy}] = append(x.cells[[2]int{ix, iy}], s)
	for range steps {
		if tMaxX < tMaxY {
			ix, tMaxX = ix+stepX, tMaxX+tDeltaX
		} else {
			iy, tMaxY = iy+stepY, tMaxY+tDeltaY
		}
		x.cells[[2]int{ix, iy}] = append(x.cells[[2]int{ix, iy}], s)
	}
	return true
}

// axis sets up the traversal along one axis: the step between cells, how
// far along the piece (0 to 1) the first cell boundary is, and how far
// apart boundaries are.
func (x *index) axis(from, to float64, i int) (step int, tMax, tDelta float64) {
	d := to - from
	switch {
	case d > 0:
		return 1, (float64(i+1)*x.cell - from) / d, x.cell / d
	case d < 0:
		return -1, (float64(i)*x.cell - from) / d, -x.cell / d
	}
	return 0, math.Inf(1), math.Inf(1)
}

// beside reports whether p, running in direction dir, lies on kept line:
// within same of it, beside it rather than past its end, and running the
// same or the opposite way.
func (x *index) beside(p, dir geom.Point, same float64) bool {
	cx, cy := x.key(p.X), x.key(p.Y)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for _, s := range x.cells[[2]int{cx + dx, cy + dy}] {
				if alongside(p, dir, s, same) {
					return true
				}
			}
		}
	}
	return false
}

// alongside tests one kept piece.
func alongside(p, dir geom.Point, s seg, same float64) bool {
	vx, vy := s.b.X-s.a.X, s.b.Y-s.a.Y
	l2 := vx*vx + vy*vy
	if l2 == 0 {
		return false
	}
	u := ((p.X-s.a.X)*vx + (p.Y-s.a.Y)*vy) / l2
	if u < 0 || u > 1 {
		return false
	}
	if math.Hypot(p.X-(s.a.X+u*vx), p.Y-(s.a.Y+u*vy)) > same {
		return false
	}
	cross := math.Abs(dir.X*vy - dir.Y*vx)
	return cross <= maxSin*math.Hypot(dir.X, dir.Y)*math.Sqrt(l2)
}

func abs(v int) int { return max(v, -v) }
