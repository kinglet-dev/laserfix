// Package svg reads and edits SVG documents for laserfix through core's
// reader and editor (the Adapter pattern: it gives them the shape of the
// ports the use case owns).
package svg

import (
	"bytes"

	"github.com/kinglet-dev/core/adapters/svg/edit"
	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/core/domain/design"
)

// Document reads and edits whole SVG documents held in memory.
type Document struct{}

// Read turns a document into a design, within the reader's limits.
func (Document) Read(src []byte) (design.Design, error) {
	return reader.Read(bytes.NewReader(src))
}

// Remove returns the document without the given elements, numbered as the
// reader numbers them. Every other byte stays as it was.
func (Document) Remove(src []byte, elements []int) ([]byte, error) {
	changes := make([]edit.Change, len(elements))
	for i, e := range elements {
		changes[i] = edit.Change{Element: e}
	}
	return edit.Apply(src, changes)
}
