package cli

import (
	"fmt"
	"io"

	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/adapters/report/jsonreport"
	"github.com/kinglet-dev/laserfix/adapters/report/text"
	"github.com/kinglet-dev/laserfix/app"
)

// The names reports use for the standard streams.
const (
	stdinName  = "standard input"
	stdoutName = "standard output"
)

// fix reads the file (or standard input) and fixes it.
func fix(o options, d Deps) (report.About, app.Report, error) {
	a := report.About{Version: d.Version, Input: o.file}
	in := d.Stdin
	if o.file == "-" {
		a.Input = stdinName
	} else {
		f, err := d.Open(o.file)
		if err != nil {
			return a, app.Report{}, fmt.Errorf("can't open %s: %w", o.file, err)
		}
		defer f.Close()
		in = f
	}
	// One byte over the reader's limit lets it say the input is too large.
	src, err := io.ReadAll(io.LimitReader(in, reader.DefaultLimits.Bytes+1))
	if err != nil {
		return a, app.Report{}, fmt.Errorf("%s: %w", a.Input, err)
	}
	r, err := d.Fixer.Fix(src, o.settings)
	if err != nil {
		return a, app.Report{}, fmt.Errorf("%s: %w", a.Input, err)
	}
	return a, r, nil
}

// deliver writes the repaired copy, when asked, then the report: on
// standard output, or on standard error when the copy goes there.
func deliver(o options, d Deps, a report.About, r app.Report) error {
	to := d.Stdout
	switch o.output {
	case "":
	case "-":
		a.Output, to = stdoutName, d.Stderr
		if _, err := d.Stdout.Write(r.Output); err != nil {
			return fmt.Errorf("can't write %s: %w", stdoutName, err)
		}
	default:
		a.Output = o.output
		if err := d.Write(o.output, r.Output); err != nil {
			return fmt.Errorf("can't write %s: %w", o.output, err)
		}
	}
	var f Formatter = text.New()
	if o.json {
		f = jsonreport.New()
	}
	if err := f.Write(to, a, r); err != nil {
		return fmt.Errorf("can't write the report: %w", err)
	}
	return nil
}
