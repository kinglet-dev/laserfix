package text

import (
	"fmt"
	"strings"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/app"
)

// maxPlaces is how many places a line of the report lists.
const maxPlaces = 5

// total is what a line of changes adds up: how many, line burned twice
// that they remove (mm), and pieces joined.
type total struct {
	n, pieces int
	length    float64
}

// kind is how the report words one kind of change, before and after.
type kind struct {
	label       string
	would, done func(t total) string
}

// kinds words each kind of change; a new kind adds an entry here.
var kinds = map[app.Kind]kind{
	app.Copy: {"Copies",
		func(t total) string {
			return fmt.Sprintf("%d %s, %.1f mm of line burned twice → remove the later copy", t.n, report.Plural(t.n, "shape", "shapes"), t.length)
		},
		func(t total) string {
			return fmt.Sprintf("%d %s removed, %.1f mm of line burned twice", t.n, report.Plural(t.n, "shape", "shapes"), t.length)
		}},
	app.Trim: {"Shared",
		func(t total) string {
			return fmt.Sprintf("%d %s, %.1f mm of line burned twice → trim it from the later shape", t.n, report.Plural(t.n, "shape", "shapes"), t.length)
		},
		func(t total) string {
			return fmt.Sprintf("%d %s trimmed, %.1f mm of line burned twice", t.n, report.Plural(t.n, "shape", "shapes"), t.length)
		}},
	app.Join: {"Split",
		func(t total) string {
			return fmt.Sprintf("%d %s in %d pieces → join each into one path", t.n, report.Plural(t.n, "outline", "outlines"), t.pieces)
		},
		func(t total) string {
			return fmt.Sprintf("%d %s joined from %d pieces", t.n, report.Plural(t.n, "outline", "outlines"), t.pieces)
		}},
}

// changes lists the changes, one line per kind, in the order kinds first
// appear, each with the places it happens.
func changes(b *strings.Builder, all []app.Change, done bool) {
	var order []app.Kind
	byKind := map[app.Kind][]app.Change{}
	for _, c := range all {
		if byKind[c.Kind] == nil {
			order = append(order, c.Kind)
		}
		byKind[c.Kind] = append(byKind[c.Kind], c)
	}
	for _, k := range order {
		cs := byKind[k]
		t := total{n: len(cs)}
		places := make([]geom.Point, len(cs))
		for i, c := range cs {
			t.length += c.Length
			t.pieces += c.Pieces
			places[i] = c.At
		}
		say := kinds[k].would
		if done {
			say = kinds[k].done
		}
		line(b, kinds[k].label, say(t), places)
	}
}

// line writes one line of the report and the places it is about.
func line(b *strings.Builder, label, text string, places []geom.Point) {
	fmt.Fprintf(b, "  %-8s %s\n", label, text)
	var at []string
	for _, p := range places[:min(len(places), maxPlaces)] {
		at = append(at, fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y))
	}
	fmt.Fprintf(b, "           at %s", strings.Join(at, ", "))
	if more := len(places) - maxPlaces; more > 0 {
		fmt.Fprintf(b, " and %d more", more)
	}
	b.WriteString("\n")
}
