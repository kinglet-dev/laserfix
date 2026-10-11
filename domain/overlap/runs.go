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
// kept line and stretches that don't, in order, adding each piece to the
// shape's own line once it has been looked at. It reports false when that
// goes over the search's limit.
func (f *finder) intervals(ps []piece) ([]interval, bool) {
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
		if !f.addOwn(pc) {
			return nil, false
		}
	}
	return out, true
}

// addOwn adds a piece, flattened, to the shape's own line.
func (f *finder) addOwn(pc piece) bool {
	l := (geom.Path{Subpaths: []geom.Subpath{{Start: pc.p0, Segments: []geom.Segment{pc.s}}}}).Flatten(f.p.Flatten)[0]
	for j := 1; j < len(l.Points); j++ {
		if !f.own.add(l.Points[j-1], l.Points[j]) {
			return false
		}
	}
	return true
}

// on reports whether a piece lies on kept line, or on its own shape's
// earlier line, at t.
func (f *finder) on(pc piece, t float64) bool {
	p, dir := pc.at(t), pc.direction(t)
	return f.kept.beside(p, dir, f.p.Same) || f.own.beside(p, dir, f.p.Same)
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
