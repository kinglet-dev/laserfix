// Package app is laserfix's use case: find what wastes material in a
// design, change the document to fix it, and measure it again.
package app

import (
	"math"

	"github.com/kinglet-dev/core/domain/geom"
)

// Kind names a kind of change; the names are stable, for --json.
type Kind string

// The kinds of change.
const (
	Copy Kind = "copy" // a shape that copies an earlier one is removed
	Trim Kind = "trim" // line lying on earlier line is cut from the later shape
	Join Kind = "join" // pieces of one outline are joined into one path
)

// Change is one change made, or that would be made, to the document.
type Change struct {
	Kind    Kind
	Element int        // the element changed, numbered as the reader numbers them
	Of      int        // for a copy: the element it copies, which stays; otherwise 0
	At      geom.Point // where the change is: the start of the shape's line, mm
	Length  float64    // line no longer burned twice, mm
	Pieces  int        // for a join: how many pieces were joined
}

// FindingKind names a kind of finding; the names are stable, for --json.
type FindingKind string

// The kinds of finding: what laserfix reports and leaves as it is.
const (
	Junction FindingKind = "junction" // three or more ends meet, so which to join is unclear
	Gap      FindingKind = "gap"      // ends too far apart to join safely
	Stacked  FindingKind = "stacked"  // line still lying on other line: different colors, or filled shapes
)

// Finding is one kind of thing left as it is, and where.
type Finding struct {
	Kind      FindingKind
	Locations []geom.Point // mm
	Size      float64      // the widest gap, or the length of line burned twice, mm
}

// Measure is how much of a design burns twice.
type Measure struct {
	Stacked   float64      // line lying on other line, mm
	Crossings int          // places where lines cross or one ends on another
	Stacks    []geom.Point // the middle of each stretch of line lying on other line
}

// Report is the outcome of fixing one document.
type Report struct {
	Width, Height float64 // mm
	Settings      Settings
	Changes       []Change
	Findings      []Finding // what was left as it is
	Before, After Measure
	Output        []byte // the fixed document
}

// Worse reports whether the changes left more line burning twice, as
// shown to 0.1 mm, or more crossings than before.
func (r Report) Worse() bool {
	return shown(r.After.Stacked) > shown(r.Before.Stacked) || r.After.Crossings > r.Before.Crossings
}

// shown is a length as reports show it, to 0.1 mm.
func shown(mm float64) float64 { return math.Round(10*mm) / 10 }
