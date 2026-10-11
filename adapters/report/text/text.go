// Package text writes laserfix's report for people to read.
package text

import (
	"fmt"
	"io"
	"strings"

	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/app"
)

// stdout is the name About uses for standard output.
const stdout = "standard output"

// issues is where to report a fix that went wrong.
const issues = "https://github.com/kinglet-dev/laserfix/issues"

// Format is the text report (a Strategy for the command line).
type Format struct{}

// New returns the text format.
func New() Format { return Format{} }

// Write writes the report in one piece.
func (Format) Write(w io.Writer, a report.About, r app.Report) error {
	var b strings.Builder
	fmt.Fprintf(&b, "laserfix %s · %s", a.Version, a.Input)
	if a.Output != "" {
		fmt.Fprintf(&b, " → %s", a.Output)
	}
	fmt.Fprintf(&b, " · %.1f × %.1f mm\n\n", r.Width, r.Height)
	switch {
	case len(r.Changes) == 0:
		b.WriteString("Nothing to fix.\n")
	case a.Output == "":
		b.WriteString("Would change\n")
		changes(&b, r.Changes, false)
		b.WriteString("\nRun again with -o FILE to write the repaired copy.\n")
	default:
		b.WriteString("Changed\n")
		changes(&b, r.Changes, true)
		checked(&b, a, r)
	}
	if a.Output != "" {
		wrote(&b, a.Output, len(r.Changes) > 0)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// checked shows the measurements before and after, and warns if worse.
func checked(b *strings.Builder, a report.About, r app.Report) {
	fmt.Fprintf(b, "\nChecked again\n  Line burned twice   %.1f mm → %.1f mm\n  Crossings           %d → %d\n",
		r.Before.Stacked, r.After.Stacked, r.Before.Crossings, r.After.Crossings)
	if r.Worse() {
		fmt.Fprintf(b, "\nWorse after the changes: check %s before you burn it, and please report this at %s\n", a.Output, issues)
	}
}

// wrote says where the copy went.
func wrote(b *strings.Builder, out string, changed bool) {
	switch {
	case out == stdout:
		b.WriteString("\nWrote the repaired SVG to standard output.\n")
	case changed:
		fmt.Fprintf(b, "\nWrote %s. Check it with: laserlint %s\n", out, out)
	default:
		fmt.Fprintf(b, "\nWrote %s, unchanged.\n", out)
	}
}
