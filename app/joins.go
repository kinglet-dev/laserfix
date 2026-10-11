package app

import (
	"maps"
	"slices"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/join"
)

// joins joins split outlines of the shapes not already changed, adds the
// changes and what it won't join to the report, and returns the edits.
// The joined path takes the place of its earliest piece; the other pieces
// leave their elements, and an element left with nothing is removed.
func joins(d design.Design, changed []bool, s Settings, r *Report) []Edit {
	found := join.Find(d.Shapes, changed, join.Params{Join: s.Join, ReportGaps: s.Gaps})
	subs := map[int][]*geom.Subpath{} // the new subpaths of each shape changed; nil drops one
	for _, c := range found.Chains {
		first := c.First()
		joined := join.Joined(d.Shapes, c)
		for _, l := range c.Links {
			if subs[l.Piece.Shape] == nil {
				subs[l.Piece.Shape] = pointers(d.Shapes[l.Piece.Shape].Local.Subpaths)
			}
			subs[l.Piece.Shape][l.Piece.Sub] = nil
		}
		subs[first.Shape][first.Sub] = &joined
		shape := d.Shapes[first.Shape]
		kind := Join
		if len(c.Links) == 1 {
			kind = Close
		}
		r.Changes = append(r.Changes, Change{Kind: kind, Element: shape.Element, At: shape.ToMM.Apply(joined.Start), Pieces: len(c.Links)})
	}
	var edits []Edit
	for _, i := range slices.Sorted(maps.Keys(subs)) {
		var p geom.Path
		for _, sub := range subs[i] {
			if sub != nil {
				p.Subpaths = append(p.Subpaths, *sub)
			}
		}
		edits = append(edits, Edit{Element: d.Shapes[i].Element, Path: p})
	}
	report(found, r)
	return edits
}

// report adds the junctions and gaps found to the report.
func report(found join.Result, r *Report) {
	if len(found.Junctions) > 0 {
		f := Finding{Kind: Junction}
		for _, j := range found.Junctions {
			f.Locations = append(f.Locations, j.At)
		}
		r.Findings = append(r.Findings, f)
	}
	if len(found.Gaps) > 0 {
		f := Finding{Kind: Gap}
		for _, g := range found.Gaps {
			f.Locations = append(f.Locations, g.At)
			f.Size = max(f.Size, g.Size)
		}
		r.Findings = append(r.Findings, f)
	}
}

// pointers lists pointers to copies of the subpaths.
func pointers(subs []geom.Subpath) []*geom.Subpath {
	out := make([]*geom.Subpath, len(subs))
	for i := range subs {
		s := subs[i]
		out[i] = &s
	}
	return out
}
