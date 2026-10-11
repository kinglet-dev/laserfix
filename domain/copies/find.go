// Package copies finds shapes that copy an earlier shape: the same paint
// and the same line, wherever it starts and whichever way it runs. A laser
// burns such a line twice.
package copies

import "github.com/kinglet-dev/core/domain/design"

// Params sets how near lines must be to count as the same.
type Params struct {
	Same    float64 // lines nearer than this, mm, are the same line
	Flatten float64 // how closely curves are followed, mm
}

// Copy says one shape copies another, earlier one.
type Copy struct {
	Shape int // the copy, an index into the shapes
	Of    int // the earliest shape that stays and that it copies
}

// Find returns the shapes that copy an earlier shape, in order. A shape
// that draws no line copies nothing.
func Find(shapes []design.Shape, p Params) []Copy {
	lines := make([][]segment, len(shapes))
	for i, s := range shapes {
		lines[i] = segments(s, p.Flatten)
	}
	var out []Copy
	removed := make([]bool, len(shapes))
	for j := range shapes {
		if len(lines[j]) == 0 {
			continue // it draws no line, so nothing burns twice
		}
		for i := range j {
			if !removed[i] && samePaint(shapes[i], shapes[j]) && covers(lines[i], lines[j], p.Same) && covers(lines[j], lines[i], p.Same) {
				out = append(out, Copy{Shape: j, Of: i})
				removed[j] = true
				break
			}
		}
	}
	return out
}

// samePaint reports whether two shapes would get the same laser settings:
// laser software chooses settings by color.
func samePaint(a, b design.Shape) bool { return a.Fill == b.Fill && a.Stroke == b.Stroke }
