// Package cli is laserfix's command line: it parses the arguments, reads
// the SVG, fixes it, writes the repaired copy when asked, and reports.
// Everything it touches is passed in (Deps), so the composition root does
// the wiring and tests run in memory.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/app"
)

// Exit codes shared by every Kinglet tool.
const (
	ExitClean   = 0 // nothing left to fix
	ExitChanges = 1 // changes are needed (preview), or the changes made things worse
	ExitError   = 2 // the file couldn't be read or written
)

// Fixer is the use case the command line drives.
type Fixer interface {
	Fix(src []byte, s app.Settings) (app.Report, error)
}

// Formatter writes a report in one format (Strategy pattern).
type Formatter interface {
	Write(w io.Writer, a report.About, r app.Report) error
}

// Deps is everything a run touches.
type Deps struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Open           func(name string) (io.ReadCloser, error)
	Write          func(name string, data []byte) error // a new file; never replaces one
	Fixer          Fixer
	Version        string
}

// Run runs laserfix with the given arguments (without the program name)
// and returns the exit code.
func Run(args []string, d Deps) int {
	o, err := parse(args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(d.Stdout, help)
		return ExitClean
	}
	if err != nil {
		fmt.Fprintf(d.Stderr, "laserfix: %v\nRun laserfix --help for usage.\n", err)
		return ExitError
	}
	if o.version {
		fmt.Fprintf(d.Stdout, "laserfix %s\n", d.Version)
		return ExitClean
	}
	a, r, err := fix(o, d)
	if err == nil {
		err = deliver(o, d, a, r)
	}
	if err != nil {
		fmt.Fprintf(d.Stderr, "laserfix: %v\n", err)
		return ExitError
	}
	if (o.output == "" && len(r.Changes) > 0) || r.Worse() {
		return ExitChanges
	}
	return ExitClean
}
