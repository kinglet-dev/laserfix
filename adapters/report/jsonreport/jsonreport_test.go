package jsonreport_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/adapters/report"
	"github.com/kinglet-dev/laserfix/adapters/report/jsonreport"
	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("The JSON report", func() {
	var r app.Report

	BeforeEach(func() {
		// Arrange (shared)
		r = app.Report{
			Width: 100, Height: 80.00004, Settings: app.DefaultSettings,
			Changes: []app.Change{{Kind: app.Copy, Element: 3, Of: 1, At: geom.Point{X: 10.00004, Y: 10}, Length: 40}},
			Before:  app.Measure{Stacked: 40, Crossings: 1}, After: app.Measure{Crossings: 1},
		}
	})

	It("writes kinglet.report/v1 with the changes and the measurements before and after", func() {
		// Arrange
		var b strings.Builder

		// Act
		err := jsonreport.New().Write(&b, report.About{Version: "0.1.0", Input: "grid.svg", Output: "fixed.svg"}, r)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(b.String()).To(MatchJSON(`{
			"schema": "kinglet.report/v1",
			"tool": {"name": "laserfix", "version": "0.1.0"},
			"input": {"name": "grid.svg", "width_mm": 100, "height_mm": 80},
			"output": {"name": "fixed.svg"},
			"settings": {"same_mm": 0.02},
			"changes": [{"kind": "copy", "element": 3, "of": 1, "x_mm": 10, "y_mm": 10, "length_mm": 40}],
			"before": {"stacked_mm": 40, "crossings": 1},
			"after": {"stacked_mm": 0, "crossings": 1},
			"findings": [],
			"summary": {"changes": 1, "written": true, "worse": false}
		}`))
	})

	It("writes no output and empty lists for a preview with nothing to fix", func() {
		// Arrange
		var b strings.Builder
		r.Changes = nil

		// Act
		err := jsonreport.New().Write(&b, report.About{Version: "0.1.0", Input: "grid.svg"}, r)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(b.String()).To(ContainSubstring(`"output": null`))
		Expect(b.String()).To(ContainSubstring(`"changes": []`))
		Expect(b.String()).To(ContainSubstring(`"summary": {
    "changes": 0,
    "written": false,`))
	})

	It("says when the changes made the design worse", func() {
		// Arrange
		var b strings.Builder
		r.After.Crossings = 2

		// Act
		Expect(jsonreport.New().Write(&b, report.About{}, r)).To(Succeed())

		// Assert
		Expect(b.String()).To(ContainSubstring(`"worse": true`))
	})

	It("passes on an error writing the report", func() {
		// Act
		err := jsonreport.New().Write(failingWriter{}, report.About{}, r)

		// Assert
		Expect(err).To(MatchError(errBroken))
	})
})

var errBroken = errors.New("broken pipe")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errBroken }
