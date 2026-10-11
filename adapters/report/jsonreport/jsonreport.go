// Package jsonreport writes laserfix's report as JSON, in the
// kinglet.report/v1 schema that every Kinglet tool shares, with laserfix's
// own changes and measurements added. Fields are only ever added within v1.
package jsonreport

import (
	"encoding/json"
	"io"
	"math"

	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/app"
)

// Schema names the report format.
const Schema = "kinglet.report/v1"

type document struct {
	Schema   string    `json:"schema"`
	Tool     tool      `json:"tool"`
	Input    input     `json:"input"`
	Output   *output   `json:"output"`
	Settings settings  `json:"settings"`
	Changes  []change  `json:"changes"`
	Before   measure   `json:"before"`
	After    measure   `json:"after"`
	Findings []finding `json:"findings"`
	Summary  summary   `json:"summary"`
}

type tool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type input struct {
	Name   string  `json:"name"`
	Width  float64 `json:"width_mm"`
	Height float64 `json:"height_mm"`
}

type output struct {
	Name string `json:"name"`
}

type settings struct {
	Same float64 `json:"same_mm"`
}

type change struct {
	Kind    string  `json:"kind"`
	Element int     `json:"element"`
	Of      *int    `json:"of,omitempty"` // only for a copy
	X       float64 `json:"x_mm"`
	Y       float64 `json:"y_mm"`
	Length  float64 `json:"length_mm"`
	Pieces  int     `json:"pieces,omitempty"` // only for a join
}

type measure struct {
	Stacked   float64 `json:"stacked_mm"`
	Crossings int     `json:"crossings"`
}

// finding is what laserfix reports and leaves as it is, in the shape of
// laserlint's findings.
type finding struct {
	Check     string     `json:"check"`
	Severity  string     `json:"severity"`
	Message   string     `json:"message"`
	Fix       string     `json:"fix"`
	Size      float64    `json:"size_mm"`
	Locations []location `json:"locations"`
}

type location struct {
	X float64 `json:"x_mm"`
	Y float64 `json:"y_mm"`
}

type summary struct {
	Changes  int  `json:"changes"`
	Findings int  `json:"findings"`
	Written  bool `json:"written"`
	Worse    bool `json:"worse"`
}

// Format is the JSON report (a Strategy for the command line).
type Format struct{}

// New returns the JSON format.
func New() Format { return Format{} }

// Write writes the report as indented JSON.
func (Format) Write(w io.Writer, a report.About, r app.Report) error {
	d := document{
		Schema:   Schema,
		Tool:     tool{Name: "laserfix", Version: a.Version},
		Input:    input{Name: a.Input, Width: mm(r.Width), Height: mm(r.Height)},
		Settings: settings{Same: mm(r.Settings.Same)},
		Changes:  []change{},
		Before:   measure{Stacked: mm(r.Before.Stacked), Crossings: r.Before.Crossings},
		After:    measure{Stacked: mm(r.After.Stacked), Crossings: r.After.Crossings},
		Findings: []finding{},
		Summary:  summary{Changes: len(r.Changes), Findings: len(r.Findings), Written: a.Output != "", Worse: r.Worse()},
	}
	if a.Output != "" {
		d.Output = &output{Name: a.Output}
	}
	for _, c := range r.Changes {
		out := change{Kind: string(c.Kind), Element: c.Element, X: mm(c.At.X), Y: mm(c.At.Y), Length: mm(c.Length), Pieces: c.Pieces}
		if c.Kind == app.Copy {
			out.Of = &c.Of
		}
		d.Changes = append(d.Changes, out)
	}
	for _, f := range r.Findings {
		words := report.Describe(f, r.Settings)
		out := finding{Check: string(f.Kind), Severity: "warning", Message: words.Message, Fix: words.Fix,
			Size: mm(f.Size), Locations: []location{}}
		for _, p := range f.Locations {
			out.Locations = append(out.Locations, location{X: mm(p.X), Y: mm(p.Y)})
		}
		d.Findings = append(d.Findings, out)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(d)
}

// mm rounds to a thousandth of a millimeter, so output is the same on
// every CPU.
func mm(v float64) float64 { return math.Round(v*1000) / 1000 }
