package overlap

import (
	"math"

	"github.com/kinglet-dev/core/domain/geom"
)

// bisections finds where a stretch starts or ends to 2⁻⁴⁰ of a piece.
const bisections = 40

// interval is a stretch of a subpath that lies on kept line, or doesn't.
type interval struct {
	from, to Pos
	on       bool
	length   float64 // mm
}

// intervals splits a subpath's pieces (in mm) into stretches that lie on
// kept line and stretches that don't, in order.
func (f *finder) intervals(ps []piece) []interval {
	var out []interval
	for i, pc := range ps {
		n := max(1, int(math.Ceil(f.arc(pc, 0, 1)/(f.p.Same/2))))
		t0, on := 0.0, f.on(pc, 0)
		prev := 0.0
		for k := 1; k <= n; k++ {
			t := float64(k) / float64(n)
			if now := f.on(pc, t); now != on {
				edge := f.edge(pc, prev, t, on)
				out = appendInterval(out, f.interval(ps, i, t0, edge, on))
				t0, on = edge, now
			}
			prev = t
		}
		out = appendInterval(out, f.interval(ps, i, t0, 1, on))
	}
	return out
}

// on reports whether a piece lies on kept line at t.
func (f *finder) on(pc piece, t float64) bool {
	return f.kept.beside(pc.at(t), pc.direction(t), f.p.Same)
}

// edge finds, by bisection, where the piece stops being in state on
// between a and b.
func (f *finder) edge(pc piece, a, b float64, on bool) float64 {
	for range bisections {
		m := (a + b) / 2
		if f.on(pc, m) == on {
			a = m
		} else {
			b = m
		}
	}
	return b
}

// interval is the stretch of piece i from t0 to t1, with its length.
func (f *finder) interval(ps []piece, i int, t0, t1 float64, on bool) interval {
	return interval{from: Pos{i, t0}, to: Pos{i, t1}, on: on, length: f.arc(ps[i], t0, t1)}
}

// arc is the length of a piece between t0 and t1, following curves within
// Flatten.
func (f *finder) arc(pc piece, t0, t1 float64) float64 {
	start, s := geom.Part(pc.p0, pc.s, t0, t1)
	total := 0.0
	for _, l := range (geom.Path{Subpaths: []geom.Subpath{{Start: start, Segments: []geom.Segment{s}}}}).Flatten(f.p.Flatten) {
		for j := 1; j < len(l.Points); j++ {
			total += dist(l.Points[j-1], l.Points[j])
		}
	}
	return total
}

// appendInterval adds v, joining it to the last interval when both are in
// the same state.
func appendInterval(out []interval, v interval) []interval {
	if n := len(out); n > 0 && out[n-1].on == v.on {
		out[n-1].to = v.to
		out[n-1].length += v.length
		return out
	}
	return append(out, v)
}
