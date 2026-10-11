package app

import (
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
	Join      float64         // ends nearer than this are joined, mm
	Gaps      float64         // gaps wider than Join and up to this are reported, mm
}

// DefaultSettings are laserfix's defaults; the measurements match
// laserlint's for its default 0.1 mm line.
var DefaultSettings = Settings{
	Same:      0.02,
	Crossings: crossing.Params{Merge: 0.1, Cell: 1, MaxSegments: 5_000_000, MaxEntries: 10_000_000},
	MaxShared: 10_000_000,
	Join:      0.05,
	Gaps:      1,
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
	r := Report{Width: d.Width, Height: d.Height, Settings: s, Output: src}
	if r.Before, err = measure(d, s); err != nil {
		return Report{}, err
	}
	if d, err = f.passes(d, s, &r); err != nil {
		return Report{}, err
	}
	r.After, err = measure(d, s)
	if shown(r.After.Stacked) > 0 {
		r.Findings = append(r.Findings, Finding{Kind: Stacked, Locations: r.After.Stacks, Size: r.After.Stacked})
	}
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
	m := Measure{Stacked: r.Stacked, Crossings: len(r.Crossings)}
	for _, st := range r.Stacks {
		m.Stacks = append(m.Stacks, st.At)
	}
	return m, err
}
