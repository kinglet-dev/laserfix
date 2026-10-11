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
	Of      int     `json:"of"`
	X       float64 `json:"x_mm"`
	Y       float64 `json:"y_mm"`
	Length  float64 `json:"length_mm"`
}

type measure struct {
	Stacked   float64 `json:"stacked_mm"`
	Crossings int     `json:"crossings"`
}

// finding is what laserfix reports but won't change; none yet.
type finding struct{}

type summary struct {
	Changes int  `json:"changes"`
	Written bool `json:"written"`
	Worse   bool `json:"worse"`
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
		Summary:  summary{Changes: len(r.Changes), Written: a.Output != "", Worse: r.Worse()},
	}
	if a.Output != "" {
		d.Output = &output{Name: a.Output}
	}
	for _, c := range r.Changes {
		d.Changes = append(d.Changes, change{Kind: string(c.Kind), Element: c.Element, Of: c.Of,
			X: mm(c.At.X), Y: mm(c.At.Y), Length: mm(c.Length)})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(d)
}

// mm rounds to a thousandth of a millimeter, so output is the same on
// every CPU.
func mm(v float64) float64 { return math.Round(v*1000) / 1000 }
