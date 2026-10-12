package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/kinglet-dev/core/domain/units"
	"github.com/kinglet-dev/laserfix/app"
)

// options are the parsed and validated arguments.
type options struct {
	file, output  string
	json, version bool
	settings      app.Settings
}

// parse reads the arguments, with flags before or after the file. Values
// are taken as text and validated here, so every message names the flag
// and says how to fix it.
func parse(args []string) (options, error) {
	var o options
	fs := flag.NewFlagSet("laserfix", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&o.output, "o", "", "")
	same := fs.String("same", "0.02mm", "")
	join := fs.String("join", "0.05mm", "")
	gaps := fs.String("report-gaps", "1mm", "")
	fs.BoolVar(&o.json, "json", false, "")
	fs.BoolVar(&o.version, "version", false, "")
	files, err := parseAll(fs, args)
	if err != nil || o.version {
		return o, err
	}
	o.settings = app.DefaultSettings
	if o.settings.Same, err = positive("--same", *same); err != nil {
		return o, err
	}
	if o.settings.Join, err = positive("--join", *join); err != nil {
		return o, err
	}
	if o.settings.Gaps, err = length("--report-gaps", *gaps); err != nil {
		return o, err
	}
	if o.settings.Gaps < o.settings.Join {
		return o, errors.New("--report-gaps must not be less than --join")
	}
	if len(files) != 1 {
		return o, errors.New("give one SVG file to fix, or - to read standard input")
	}
	o.file = files[0]
	return o, nil
}

// length reads a length that must have a unit, in millimeters.
func length(flag, value string) (float64, error) {
	l, err := units.ParseLength(value)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", flag, err)
	}
	return l.Millimeters(), nil
}

// positive reads a length that must also be more than zero.
func positive(flag, value string) (float64, error) {
	v, err := length(flag, value)
	if err == nil && v == 0 {
		err = fmt.Errorf("%s must be more than zero", flag)
	}
	return v, err
}

// parseAll parses flags wherever they come, and returns the other
// arguments: Go's flag package stops at the first one that isn't a flag.
func parseAll(fs *flag.FlagSet, args []string) ([]string, error) {
	var rest []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return rest, nil
		}
		rest = append(rest, fs.Arg(0))
		args = fs.Args()[1:]
	}
}
