// Package app is laserfix's use case: find what wastes material in a
// design, change the document to fix it, and measure it again.
package app

import "github.com/kinglet-dev/core/domain/geom"

// Kind names a kind of change; the names are stable, for --json.
type Kind string

// Copy is the removal of a shape that copies an earlier one.
const Copy Kind = "copy"

// Change is one change made, or that would be made, to the document.
type Change struct {
	Kind    Kind
	Element int        // the element changed, numbered as the reader numbers them
	Of      int        // for a copy: the element it copies, which stays
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
	Changes       []Change
	Before, After Measure
	Output        []byte // the fixed document
}
