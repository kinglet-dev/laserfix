// Package svg reads and edits SVG documents for laserfix through core's
// reader and editor (the Adapter pattern: it gives them the shape of the
// ports the use case owns).
package svg

import (
	"bytes"

	"github.com/kinglet-dev/core/adapters/svg/edit"
	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/core/adapters/svg/syntax"
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/laserfix/app"
)

// Document reads and edits whole SVG documents held in memory.
type Document struct{}

// Read turns a document into a design, within the reader's limits.
func (Document) Read(src []byte) (design.Design, error) {
	return reader.Read(bytes.NewReader(src))
}

// Edit returns the document with each element removed, or redrawn as a
// path in its own units. Every other byte stays as it was.
func (Document) Edit(src []byte, edits []app.Edit) ([]byte, error) {
	changes := make([]edit.Change, len(edits))
	for i, e := range edits {
		// An empty path writes no path data, and a change with neither
		// path data nor markup removes the element.
		changes[i] = edit.Change{Element: e.Element, Redraw: syntax.FormatPath(e.Path)}
	}
	return edit.Apply(src, changes)
}
