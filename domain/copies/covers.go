package copies

import (
	"math"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

// segment is one straight piece of a flattened shape.
type segment struct{ a, b geom.Point }

// segments flattens a shape into the straight pieces it draws, closing
// closed lines.
func segments(s design.Shape, flatten float64) []segment {
	var out []segment
	for _, l := range s.Path.Flatten(flatten) {
		pts := l.Points
		if l.Closed { // Flatten gives every line at least its start
			pts = append(pts[:len(pts):len(pts)], pts[0])
		}
		for i := 1; i < len(pts); i++ {
			out = append(out, segment{pts[i-1], pts[i]})
		}
	}
	return out
}

// covers reports whether every point of a lies within same of b: half of
// the discrete Hausdorff distance test. Points are sampled at each piece's
// ends and every same/2 between, so a gap wider than same can't slip
// between samples.
func covers(a, b []segment, same float64) bool {
	for _, s := range a {
		n := int(math.Ceil(length(s) / (same / 2)))
		for k := 0; k <= n; k++ {
			if !near(lerp(s.a, s.b, float64(k)/float64(max(n, 1))), b, same) {
				return false
			}
		}
	}
	return true
}

// near reports whether p lies within d of any piece of b.
func near(p geom.Point, b []segment, d float64) bool {
	for _, s := range b {
		if distance(p, s) <= d {
			return true
		}
	}
	return false
}

// distance is from p to the nearest point of s.
func distance(p geom.Point, s segment) float64 {
	vx, vy := s.b.X-s.a.X, s.b.Y-s.a.Y
	t := 0.0
	if l2 := vx*vx + vy*vy; l2 > 0 {
		t = math.Max(0, math.Min(1, ((p.X-s.a.X)*vx+(p.Y-s.a.Y)*vy)/l2))
	}
	return math.Hypot(p.X-(s.a.X+t*vx), p.Y-(s.a.Y+t*vy))
}

func length(s segment) float64 { return math.Hypot(s.b.X-s.a.X, s.b.Y-s.a.Y) }

func lerp(a, b geom.Point, t float64) geom.Point {
	return geom.Point{X: a.X + t*(b.X-a.X), Y: a.Y + t*(b.Y-a.Y)}
}
