package overlap

import "github.com/kinglet-dev/core/domain/design"

// minRun is how long, in multiples of Same, a stretch must lie on kept
// line to be cut: shorter ones are where lines meet at a shallow angle.
const minRun = 5

// trim works out what is left of a shape, in its own units, and whether
// anything was cut.
func (f *finder) trim(s design.Shape) (Trim, bool) {
	var t Trim
	changed := false
	for si, sub := range s.Local.Subpaths {
		ps := pieces(sub)
		for i := range ps {
			ps[i] = ps[i].transform(s.ToMM)
		}
		iv := f.tidy(f.intervals(ps), sub.Closed)
		if !anyOn(iv) {
			t.Keep = append(t.Keep, Stretch{Sub: si, Whole: true})
			continue
		}
		for _, v := range iv {
			if v.on {
				if !changed {
					t.At = ps[v.from.Piece].at(v.from.T)
				}
				changed = true
				t.Removed += v.length
			} else {
				t.Keep = append(t.Keep, Stretch{Sub: si, From: v.from, To: v.to})
			}
		}
	}
	return t, changed
}

// tidy joins a closed subpath's last stretch to its first when both are in
// the same state, keeps stretches too short to be shared line, and cuts
// leftover bits of line shorter than Same between cut stretches.
func (f *finder) tidy(iv []interval, closed bool) []interval {
	iv = merge(iv, closed)
	for i := range iv {
		if iv[i].on && iv[i].length < minRun*f.p.Same {
			iv[i].on = false
		}
	}
	iv = merge(iv, closed)
	if anyOn(iv) {
		for i := range iv {
			if !iv[i].on && iv[i].length < f.p.Same {
				iv[i].on = true
			}
		}
	}
	return merge(iv, closed)
}

// merge joins neighbors in the same state, and for a closed subpath the
// last to the first.
func merge(iv []interval, closed bool) []interval {
	var out []interval
	for _, v := range iv {
		out = appendInterval(out, v)
	}
	if n := len(out); closed && n > 1 && out[0].on == out[n-1].on {
		out[0].from = out[n-1].from
		out[0].length += out[n-1].length
		out = out[:n-1]
	}
	return out
}

func anyOn(iv []interval) bool {
	for _, v := range iv {
		if v.on {
			return true
		}
	}
	return false
}

// keep adds what is left of a shape to the kept line, and reports false
// when that goes over the search's limit.
func (f *finder) keep(s design.Shape, t Trim, changed bool) bool {
	path := s.Local
	if changed {
		path = Cut(s.Local, t.Keep)
	}
	for _, l := range path.Transform(s.ToMM).Flatten(f.p.Flatten) {
		pts := l.Points
		if l.Closed {
			pts = append(pts[:len(pts):len(pts)], pts[0])
		}
		for i := 1; i < len(pts); i++ {
			if !f.kept.add(pts[i-1], pts[i]) {
				return false
			}
		}
	}
	return true
}
