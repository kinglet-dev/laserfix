package app

import (
	"errors"
	"fmt"

	"github.com/kinglet-dev/core/domain/design"
)

// maxPasses limits how often Fix plans and makes changes. Each fix can
// leave work for another (trimming leaves pieces to join, and a joined
// path can run back along itself), so Fix repeats until a pass changes
// nothing, or this many passes have run.
const maxPasses = 5

// errLostTrack means the changed document's shapes no longer match the
// shapes kept, so changes could no longer be named by their element.
var errLostTrack = errors.New("its shapes no longer match the ones kept")

// passes makes the changes in r.Output, pass after pass, and returns the
// design as it is in the end. Every change is named by the element's
// number in the file Fix was given; the findings are the last pass's.
func (f *Fixer) passes(d design.Design, s Settings, r *Report) (design.Design, error) {
	input := make([]int, len(d.Shapes)) // each shape's element in the file given
	for i, shape := range d.Shapes {
		input[i] = shape.Element
	}
	for range maxPasses {
		pass := Report{}
		byShape, err := plan(d, s, &pass)
		if err != nil {
			return d, err
		}
		r.Changes = append(r.Changes, renumber(pass.Changes, d, input)...)
		r.Findings = pass.Findings
		if len(byShape) == 0 {
			break
		}
		if d, input, err = f.apply(d, input, byShape, r); err != nil {
			return d, err
		}
	}
	return d, nil
}

// apply makes one pass's edits and reads the document again, keeping track
// of which element of the file given each remaining shape is.
func (f *Fixer) apply(d design.Design, input []int, byShape map[int]Edit, r *Report) (design.Design, []int, error) {
	var edits []Edit
	var kept []int
	for i := range d.Shapes {
		e, changed := byShape[i]
		if changed {
			edits = append(edits, e)
		}
		if !changed || len(e.Path.Subpaths) > 0 {
			kept = append(kept, input[i])
		}
	}
	out, err := f.edit.Edit(r.Output, edits)
	if err != nil {
		return d, nil, fmt.Errorf("changing the document: %w", err)
	}
	if d, err = f.read.Read(out); err != nil {
		return d, nil, fmt.Errorf("reading the changed document: %w", err)
	}
	if len(d.Shapes) != len(kept) {
		return d, nil, fmt.Errorf("reading the changed document: %w", errLostTrack)
	}
	r.Output = out
	return d, kept, nil
}

// renumber names the elements in a pass's changes by their numbers in the
// file Fix was given, rather than in the document the pass read. A copy's
// Of needs no renumbering: copies are only ever found in the first pass,
// since later passes see filled shapes unchanged and hairlines already
// trimmed of every copy.
func renumber(changes []Change, d design.Design, input []int) []Change {
	given := map[int]int{}
	for i, shape := range d.Shapes {
		given[shape.Element] = input[i]
	}
	for i := range changes {
		changes[i].Element = given[changes[i].Element]
	}
	return changes
}
