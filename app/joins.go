package app

import (
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/join"
)

// joins joins split outlines of the shapes not gone, adds the changes and
// what it won't join to the report, and returns the edits by shape. The
// joined path takes the place of its earliest piece; the other pieces
// leave their elements, and an element left with nothing is removed.
func joins(shapes []design.Shape, gone []bool, s Settings, r *Report) map[int]Edit {
	found := join.Find(shapes, gone, join.Params{Join: s.Join, ReportGaps: s.Gaps})
	subs := map[int][]*geom.Subpath{} // the new subpaths of each shape changed; nil drops one
	for _, c := range found.Chains {
		first := c.First()
		joined := join.Joined(shapes, c)
		for _, l := range c.Links {
			if subs[l.Piece.Shape] == nil {
				subs[l.Piece.Shape] = pointers(shapes[l.Piece.Shape].Local.Subpaths)
			}
			subs[l.Piece.Shape][l.Piece.Sub] = nil
		}
		subs[first.Shape][first.Sub] = &joined
		shape := shapes[first.Shape]
		kind := Join
		if len(c.Links) == 1 {
			kind = Close
		}
		r.Changes = append(r.Changes, Change{Kind: kind, Element: shape.Element, At: shape.ToMM.Apply(joined.Start), Pieces: len(c.Links)})
	}
	edits := map[int]Edit{}
	for i, list := range subs {
		var p geom.Path
		for _, sub := range list {
			if sub != nil {
				p.Subpaths = append(p.Subpaths, *sub)
			}
		}
		edits[i] = Edit{Element: shapes[i].Element, Path: p}
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
