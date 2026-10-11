package text

import (
	"fmt"
	"strings"

	"github.com/kinglet-dev/laserfix/app"
)

// maxPlaces is how many places a line of the report lists.
const maxPlaces = 5

// kind is how the report words one kind of change.
type kind struct {
	label, one, many string // the label, and the thing changed, one or many
	would, done      string // what laserfix would do, and what it did
}

// kinds words each kind of change; a new kind adds an entry here.
var kinds = map[app.Kind]kind{
	app.Copy: {label: "Copies", one: "shape", many: "shapes", would: "remove the later copy", done: "removed"},
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
		line(b, kinds[k], byKind[k], done)
	}
}

// line writes one kind's changes.
func line(b *strings.Builder, k kind, cs []app.Change, done bool) {
	total := 0.0
	for _, c := range cs {
		total += c.Length
	}
	thing := k.many
	if len(cs) == 1 {
		thing = k.one
	}
	if done {
		fmt.Fprintf(b, "  %-8s %d %s %s, %.1f mm of line burned twice\n", k.label, len(cs), thing, k.done, total)
	} else {
		fmt.Fprintf(b, "  %-8s %d %s, %.1f mm of line burned twice → %s\n", k.label, len(cs), thing, total, k.would)
	}
	var places []string
	for _, c := range cs[:min(len(cs), maxPlaces)] {
		places = append(places, fmt.Sprintf("(%.1f, %.1f)", c.At.X, c.At.Y))
	}
	fmt.Fprintf(b, "           at %s", strings.Join(places, ", "))
	if more := len(cs) - maxPlaces; more > 0 {
		fmt.Fprintf(b, " and %d more", more)
	}
	b.WriteString("\n")
}
