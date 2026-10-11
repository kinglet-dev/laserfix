// Command laserfix repairs an SVG so the laser burns each line once. This
// is the composition root: it only wires the parts together.
package main

import (
	"os"
	"runtime/debug"

	"github.com/kinglet-dev/core/adapters/buildinfo"
	"github.com/kinglet-dev/core/adapters/files"
	"github.com/kinglet-dev/laserfix/adapters/cli"
	"github.com/kinglet-dev/laserfix/adapters/output"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

// version is set at release time with -ldflags "-X main.version=…"; a
// copy built with go install reports the module version Go recorded.
var version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.Deps{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Open:    files.Open,
		Write:   output.Write,
		Fixer:   app.NewFixer(svg.Document{}, svg.Document{}),
		Version: buildinfo.Version(version, debug.ReadBuildInfo),
	}))
}
