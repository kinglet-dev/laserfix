package join

import "github.com/kinglet-dev/core/domain/design"

// Find finds the chains to join and what to report. Shapes marked in skip
// are left out.
func Find(shapes []design.Shape, skip []bool, p Params) Result {
	pieces, ends := candidates(shapes, skip)
	stroke := func(e end) string { return shapes[pieces[e.piece].Shape].Stroke }
	joinable := near(ends, stroke, p.Join)
	var r Result
	partner := map[int]int{}
	counted := make([]bool, len(ends))
	for i, ns := range joinable {
		switch {
		case len(ns) == 1 && len(joinable[ns[0]]) == 1:
			partner[i] = ns[0]
		case len(ns) > 1 && !counted[i]:
			r.Junctions = append(r.Junctions, Junction{At: ends[i].at, Ends: len(ns) + 1})
			for _, o := range ns {
				counted[o] = true
			}
		}
	}
	r.Gaps = gaps(ends, joinable, near(ends, stroke, p.ReportGaps))
	for _, c := range chains(len(pieces), partner) {
		n := c.links[0] / 2 // the first piece; its ends are 2n and 2n+1
		single := len(c.links) == 1 && ends[2*n].at == ends[2*n+1].at
		if !single {
			r.Chains = append(r.Chains, c.chain(pieces, ends))
		}
	}
	return r
}

// gaps reports each pair of ends, neither joinable, within the report
// distance, once, at the end that comes first.
func gaps(ends []end, joinable, reportable [][]int) []Gap {
	var out []Gap
	for i, ns := range reportable {
		if len(joinable[i]) > 0 {
			continue
		}
		for _, o := range ns {
			if o > i && len(joinable[o]) == 0 && (facing(ends[i], ends[o].at) || facing(ends[o], ends[i].at)) {
				out = append(out, Gap{At: ends[i].at, Size: dist(ends[i].at, ends[o].at)})
			}
		}
	}
	return out
}
