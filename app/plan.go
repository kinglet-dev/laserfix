package app

import (
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/laserfix/domain/copies"
	"github.com/kinglet-dev/laserfix/domain/overlap"
)

// plan works out the changes, adds them to the report, and returns the
// edits that make them: copies first, then line lying on earlier line,
// with the copies already gone.
func plan(d design.Design, s Settings, r *Report) ([]Edit, error) {
	var edits []Edit
	gone := make([]bool, len(d.Shapes))
	for _, c := range copies.Find(d.Shapes, copies.Params{Same: s.Same, Flatten: flatten}) {
		r.Changes = append(r.Changes, copyChange(d.Shapes[c.Shape], d.Shapes[c.Of]))
		edits = append(edits, Edit{Element: d.Shapes[c.Shape].Element})
		gone[c.Shape] = true
	}
	trims, err := overlap.Find(d.Shapes, gone, overlap.Params{Same: s.Same, Flatten: flatten, MaxEntries: s.MaxShared})
	if err != nil {
		return nil, err
	}
	for _, t := range trims {
		shape := d.Shapes[t.Shape]
		r.Changes = append(r.Changes, Change{Kind: Trim, Element: shape.Element, At: t.At, Length: t.Removed})
		edits = append(edits, Edit{Element: shape.Element, Path: overlap.Cut(shape.Local, t.Keep)})
		gone[t.Shape] = true
	}
	return append(edits, joins(d, gone, s, r)...), nil
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
