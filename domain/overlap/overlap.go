// Package overlap finds line lying on earlier line of the same paint and
// trims it from the later shape, so the laser burns it once; the earlier
// shape stays whole. Only hairlines (stroke, no fill) are trimmed: cutting
// a filled shape's outline would break its fill.
//
// A stretch of a shape lies on other line where it runs within Same of
// it, beside it rather than off its end, and in the same or the opposite
// direction (within maxAngle). Lines that cross or meet at a corner are
// never trimmed. Points are sampled along each piece of line at most Same/2
// apart, and the ends of each stretch are found exactly by bisection.
package overlap

import "github.com/kinglet-dev/core/domain/geom"

// Params sets how near lines must be to count as the same.
type Params struct {
	Same       float64 // lines nearer than this, mm, are the same line
	Flatten    float64 // how closely curves are followed, mm
	MaxEntries int     // limit on the search's work, so huge designs fail clearly
}

// Pos is a place along a subpath: a piece of it and how far along that
// piece, from 0 to 1. The pieces are the subpath's segments, then the
// side that closes a closed subpath, when it has length.
type Pos struct {
	Piece int
	T     float64
}

// Stretch is a part of a subpath to keep: all of it as drawn, or the
// line from From forward to To, wrapping past the end of a closed subpath.
type Stretch struct {
	Sub      int
	Whole    bool
	From, To Pos
}

// Trim is what is left of one shape once line lying on earlier line is cut
// from it.
type Trim struct {
	Shape   int        // an index into the shapes
	Keep    []Stretch  // what stays; none means the shape goes
	Removed float64    // line cut away, mm
	At      geom.Point // where the first stretch cut away starts, mm
}
