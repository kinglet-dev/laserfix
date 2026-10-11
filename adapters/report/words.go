package report

import (
	"fmt"

	"github.com/kinglet-dev/laserfix/app"
)

// Words is how every report words a finding: a short label, what was
// found, and how to fix it.
type Words struct{ Label, Message, Fix string }

// Describe words a finding; a new kind of finding adds a case here.
func Describe(f app.Finding, s app.Settings) Words {
	n := len(f.Locations)
	switch f.Kind {
	case app.Junction:
		return Words{"Junction", fmt.Sprintf("%d %s where three or more ends meet, so which to join is unclear", n, Plural(n, "place", "places")),
			"join them by hand"}
	case app.Gap:
		return Words{"Gap", fmt.Sprintf("%d %s up to %.1f mm, wider than %g mm (--join)", n, Plural(n, "gap", "gaps"), f.Size, s.Join),
			"close them by hand, or raise --join"}
	}
	return Words{"Overlap", fmt.Sprintf("%.1f mm of line still burns twice, shared by different colors or filled shapes", f.Size),
		"combine them by hand if they should burn once"}
}

// Plural picks the word for n things.
func Plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
