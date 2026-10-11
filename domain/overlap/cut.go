package overlap

import "github.com/kinglet-dev/core/domain/geom"

// Cut returns what is left of a path once only the stretches in keep stay.
// Curves stay curves: each is cut exactly where a stretch starts or ends.
func Cut(p geom.Path, keep []Stretch) geom.Path {
	var out geom.Path
	for _, k := range keep {
		sub := p.Subpaths[k.Sub]
		if k.Whole {
			out.Subpaths = append(out.Subpaths, sub)
			continue
		}
		ps := pieces(sub)
		var s geom.Subpath
		add := func(i int, t0, t1 float64) {
			start, seg := geom.Part(ps[i].p0, ps[i].s, t0, t1)
			if len(s.Segments) == 0 {
				s.Start = start
			}
			s.Segments = append(s.Segments, seg)
		}
		if k.From.Piece == k.To.Piece && k.From.T <= k.To.T {
			add(k.From.Piece, k.From.T, k.To.T)
		} else {
			// Round to the end of From's piece first: when the stretch wraps
			// all the way round, To is on that same piece.
			add(k.From.Piece, k.From.T, 1)
			for i := (k.From.Piece + 1) % len(ps); i != k.To.Piece; i = (i + 1) % len(ps) {
				add(i, 0, 1)
			}
			add(k.To.Piece, 0, k.To.T)
		}
		out.Subpaths = append(out.Subpaths, s)
	}
	return out
}
