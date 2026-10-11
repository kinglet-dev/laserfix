package app

import (
	"fmt"
	"math"

	"github.com/kinglet-dev/core/domain/crossing"
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

// flatten is how closely lines follow curves, mm: far below the width of
// a laser line.
const flatten = 0.005

// Reader turns a document into a design (a port).
type Reader interface {
	Read(src []byte) (design.Design, error)
}

// Edit is one change to a document: an element, numbered as the reader
// numbers them, and its new path in its own units; an empty path removes it.
type Edit struct {
	Element int
	Path    geom.Path
}

// Editor makes edits to a document (a port).
type Editor interface {
	Edit(src []byte, edits []Edit) ([]byte, error)
}

// Settings set what counts as a fault and how the design is measured.
type Settings struct {
	Same      float64         // lines nearer than this are the same line, mm
	Crossings crossing.Params // how line burned twice is measured
	MaxShared int             // limit on the search for shared line, so huge designs fail clearly
}

// DefaultSettings are laserfix's defaults; the measurements match
// laserlint's for its default 0.1 mm line.
var DefaultSettings = Settings{
	Same:      0.02,
	Crossings: crossing.Params{Merge: 0.1, Cell: 1, MaxSegments: 5_000_000, MaxEntries: 10_000_000},
	MaxShared: 10_000_000,
}

// Fixer fixes documents through its ports, given at construction.
type Fixer struct {
	read Reader
	edit Editor
}

// NewFixer returns a fixer that reads and edits with the given ports.
func NewFixer(r Reader, e Editor) *Fixer { return &Fixer{read: r, edit: e} }

// Fix finds the changes, makes them in a copy of the document, and
// measures the design before and after. src itself is never changed.
func (f *Fixer) Fix(src []byte, s Settings) (Report, error) {
	d, err := f.read.Read(src)
	if err != nil {
		return Report{}, err
	}
	r := Report{Width: d.Width, Height: d.Height, Settings: s}
	if r.Before, err = measure(d, s); err != nil {
		return Report{}, err
	}
	edits, err := plan(d, s, &r)
	if err != nil {
		return Report{}, err
	}
	if r.Output, err = f.edit.Edit(src, edits); err != nil {
		return Report{}, fmt.Errorf("changing the document: %w", err)
	}
	after, err := f.read.Read(r.Output)
	if err != nil {
		return Report{}, fmt.Errorf("reading the changed document: %w", err)
	}
	r.After, err = measure(after, s)
	return r, err
}

// length is a line's length, mm, with the closing side of a closed line.
func length(l geom.Polyline) float64 {
	pts := l.Points
	if l.Closed {
		pts = append(pts[:len(pts):len(pts)], pts[0])
	}
	total := 0.0
	for i := 1; i < len(pts); i++ {
		total += math.Hypot(pts[i].X-pts[i-1].X, pts[i].Y-pts[i-1].Y)
	}
	return total
}

// measure finds how much of a design burns twice.
func measure(d design.Design, s Settings) (Measure, error) {
	r, err := crossing.Find(d.ScoredLines(flatten), s.Crossings)
	return Measure{Stacked: r.Stacked, Crossings: len(r.Crossings)}, err
}
