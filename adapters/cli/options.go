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
	fs.BoolVar(&o.json, "json", false, "")
	fs.BoolVar(&o.version, "version", false, "")
	files, err := parseAll(fs, args)
	if err != nil || o.version {
		return o, err
	}
	o.settings = app.DefaultSettings
	l, err := units.ParseLength(*same)
	if err != nil {
		return o, fmt.Errorf("--same: %w", err)
	}
	if o.settings.Same = l.Millimetres(); o.settings.Same == 0 {
		return o, errors.New("--same must be more than zero")
	}
	if len(files) != 1 {
		return o, errors.New("give one SVG file to fix, or - to read standard input")
	}
	o.file = files[0]
	return o, nil
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
