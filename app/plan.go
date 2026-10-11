package app

import (
	"maps"
	"slices"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/laserfix/domain/copies"
	"github.com/kinglet-dev/laserfix/domain/overlap"
)

// plan works out the changes, adds them to the report, and returns the
// edits that make them, one per element: copies first, then line lying on
// earlier line with the copies gone, then outlines split into pieces as
// they are after trimming, so one run fixes everything it can.
func plan(d design.Design, s Settings, r *Report) ([]Edit, error) {
	edits := map[int]Edit{} // by shape; a later edit of a shape replaces an earlier one
	gone := make([]bool, len(d.Shapes))
	for _, c := range copies.Find(d.Shapes, copies.Params{Same: s.Same, Flatten: flatten}) {
		r.Changes = append(r.Changes, copyChange(d.Shapes[c.Shape], d.Shapes[c.Of]))
		edits[c.Shape] = Edit{Element: d.Shapes[c.Shape].Element}
		gone[c.Shape] = true
	}
	trims, err := overlap.Find(d.Shapes, gone, overlap.Params{Same: s.Same, Flatten: flatten, MaxEntries: s.MaxShared})
	if err != nil {
		return nil, err
	}
	trimmed := slices.Clone(d.Shapes)
	for _, t := range trims {
		shape := &trimmed[t.Shape]
		r.Changes = append(r.Changes, Change{Kind: Trim, Element: shape.Element, At: t.At, Length: t.Removed})
		shape.Local = overlap.Cut(shape.Local, t.Keep)
		shape.Path = shape.Local.Transform(shape.ToMM)
		edits[t.Shape] = Edit{Element: shape.Element, Path: shape.Local}
	}
	maps.Copy(edits, joins(trimmed, gone, s, r))
	var out []Edit
	for _, i := range slices.Sorted(maps.Keys(edits)) {
		out = append(out, edits[i])
	}
	return out, nil
}

// copyChange describes removing shape, a copy of of.
func copyChange(shape, of design.Shape) Change {
	lines := design.Design{Shapes: []design.Shape{shape}}.ScoredLines(flatten)
	c := Change{Kind: Copy, Element: shape.Element, Of: of.Element, At: lines[0].Points[0]}
	for _, l := range lines {
		c.Length += length(l)
	}
	return c
}
