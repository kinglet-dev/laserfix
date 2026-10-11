package overlap

import (
	"errors"
	"math"

	"github.com/kinglet-dev/core/domain/design"
)

// ErrTooComplex means the design has too much line to search within the limits.
var ErrTooComplex = errors.New("the design has too much line to check for shared line")

// finder holds the line kept so far for one paint.
type finder struct {
	p    Params
	kept *index
}

// Find returns the trims for the shapes, in order. Shapes marked in gone
// (copies already removed, say) are skipped, and hairlines are only
// compared with earlier hairlines of the same stroke.
func Find(shapes []design.Shape, gone []bool, p Params) ([]Trim, error) {
	cell := math.Max(1, p.Same)
	byStroke := map[string]*finder{}
	var out []Trim
	for i, s := range shapes {
		if (i < len(gone) && gone[i]) || s.Fill != "none" {
			continue
		}
		f := byStroke[s.Stroke]
		if f == nil {
			f = &finder{p: p, kept: newIndex(cell, p.MaxEntries)}
			byStroke[s.Stroke] = f
		}
		t, changed := f.trim(s)
		if changed {
			t.Shape = i
			out = append(out, t)
		}
		if !f.keep(s, t, changed) {
			return nil, ErrTooComplex
		}
	}
	return out, nil
}
