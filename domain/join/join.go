// Package join joins open pieces of line whose ends meet into one path, as
// an outline split into pieces should be, closing it when the chain comes
// back to where it started (the line merging of JTS's LineMerger). Only
// hairlines of the same stroke are joined. Ends are paired only when each
// is the other's one end within Join; where three or more ends meet, it is
// unclear which to join, so they are reported instead, as are gaps wider
// than Join and up to ReportGaps.
package join

import "github.com/kinglet-dev/core/domain/geom"

// Params sets which ends meet.
type Params struct {
	Join       float64 // ends nearer than this, mm, are joined
	ReportGaps float64 // ends further apart than Join and up to this, mm, are reported
}

// Piece is an open subpath: the shape it belongs to and its place there.
type Piece struct{ Shape, Sub int }

// Link is a piece in a chain, drawn the other way when Reversed.
type Link struct {
	Piece    Piece
	Reversed bool
}

// Chain is pieces that join end to end, in order.
type Chain struct {
	Links  []Link
	Closed bool
}

// Junction is a place where three or more ends meet.
type Junction struct {
	At   geom.Point
	Ends int
}

// Gap is two ends too far apart to join, and near enough to report.
type Gap struct {
	At   geom.Point // the end that comes first
	Size float64    // mm
}

// Result is what Find found.
type Result struct {
	Chains    []Chain
	Junctions []Junction
	Gaps      []Gap
}
