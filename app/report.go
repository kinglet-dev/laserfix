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
)

// Change is one change made, or that would be made, to the document.
type Change struct {
	Kind    Kind
	Element int        // the element changed, numbered as the reader numbers them
	Of      int        // for a copy: the element it copies, which stays; otherwise 0
	At      geom.Point // where the change is: the start of the shape's line, mm
	Length  float64    // line no longer burned twice, mm
}

// Measure is how much of a design burns twice.
type Measure struct {
	Stacked   float64 // line lying on other line, mm
	Crossings int     // places where lines cross or one ends on another
}

// Report is the outcome of fixing one document.
type Report struct {
	Width, Height float64 // mm
	Settings      Settings
	Changes       []Change
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
