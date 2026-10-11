package overlap

import (
	"math"

	"github.com/kinglet-dev/core/domain/geom"
)

// piece is one segment of a subpath with the point it starts from.
type piece struct {
	p0 geom.Point
	s  geom.Segment
}

// pieces lists a subpath's pieces: its segments, then the side that
// closes it when it is closed and doesn't end where it starts.
func pieces(sub geom.Subpath) []piece {
	var out []piece
	at := sub.Start
	for _, s := range sub.Segments {
		out = append(out, piece{at, s})
		at = s.To
	}
	if sub.Closed && at != sub.Start {
		out = append(out, piece{at, geom.LineTo(sub.Start)})
	}
	return out
}

// transform maps a piece's points with m.
func (p piece) transform(m geom.Matrix) piece {
	return piece{m.Apply(p.p0), geom.Segment{C1: m.Apply(p.s.C1), C2: m.Apply(p.s.C2), To: m.Apply(p.s.To), Curved: p.s.Curved}}
}

// at is the point a fraction t along the piece.
func (p piece) at(t float64) geom.Point {
	start, _ := geom.Part(p.p0, p.s, t, t)
	return start
}

// direction is the way the piece runs at t, measured across a short step
// so it is defined even where a curve's control point sits on its end.
func (p piece) direction(t float64) geom.Point {
	const h = 1e-4
	a, b := p.at(math.Max(0, t-h)), p.at(math.Min(1, t+h))
	return geom.Point{X: b.X - a.X, Y: b.Y - a.Y}
}

func dist(a, b geom.Point) float64 { return math.Hypot(b.X-a.X, b.Y-a.Y) }
