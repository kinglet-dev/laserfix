package app

import (
	"fmt"
	"math"

	"github.com/kinglet-dev/core/domain/crossing"
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/copies"
)

// flatten is how closely lines follow curves, mm: far below the width of
// a laser line.
const flatten = 0.005

// Reader turns a document into a design (a port).
type Reader interface {
	Read(src []byte) (design.Design, error)
}

// Remover removes elements from a document, numbered as the reader
// numbers them (a port).
type Remover interface {
	Remove(src []byte, elements []int) ([]byte, error)
}

// Settings set what counts as a fault and how the design is measured.
type Settings struct {
	Same      float64         // lines nearer than this are the same line, mm
	Crossings crossing.Params // how line burned twice is measured
}

// DefaultSettings are laserfix's defaults; the measurements match
// laserlint's for its default 0.1 mm line.
var DefaultSettings = Settings{
	Same:      0.02,
	Crossings: crossing.Params{Merge: 0.1, Cell: 1, MaxSegments: 5_000_000, MaxEntries: 10_000_000},
}

// Fixer fixes documents through its ports, given at construction.
type Fixer struct {
	read   Reader
	remove Remover
}

// NewFixer returns a fixer that reads and edits with the given ports.
func NewFixer(r Reader, rm Remover) *Fixer { return &Fixer{read: r, remove: rm} }

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
	var gone []int
	for _, c := range copies.Find(d.Shapes, copies.Params{Same: s.Same, Flatten: flatten}) {
		r.Changes = append(r.Changes, copyChange(d.Shapes[c.Shape], d.Shapes[c.Of]))
		gone = append(gone, d.Shapes[c.Shape].Element)
	}
	if r.Output, err = f.remove.Remove(src, gone); err != nil {
		return Report{}, fmt.Errorf("changing the document: %w", err)
	}
	after, err := f.read.Read(r.Output)
	if err != nil {
		return Report{}, fmt.Errorf("reading the changed document: %w", err)
	}
	r.After, err = measure(after, s)
	return r, err
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
