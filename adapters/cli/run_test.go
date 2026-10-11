package cli_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserfix/adapters/cli"
	"github.com/kinglet-dev/laserfix/app"
)

var _ = Describe("Run", func() {
	var w *world

	BeforeEach(func() {
		// Arrange (shared)
		w = newWorld()
		w.files["grid.svg"] = doc(sq + sq)
		w.files["clean.svg"] = doc(sq)
	})

	It("previews what would change, writes nothing, and exits 1", func() {
		// Act
		code := w.run("grid.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitChanges))
		Expect(w.stdout.String()).To(HavePrefix("laserfix 0.1.0 · grid.svg · 100.0 × 100.0 mm\n\nWould change\n"))
		Expect(w.written).To(BeEmpty())
	})

	It("exits 0 when there is nothing to fix", func() {
		// Act
		code := w.run("clean.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitClean))
		Expect(w.stdout.String()).To(ContainSubstring("Nothing to fix."))
	})

	It("writes the repaired copy with -o after the file, and exits 0", func() {
		// Act
		code := w.run("grid.svg", "-o", "fixed.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitClean))
		Expect(w.written).To(Equal(map[string]string{"fixed.svg": doc(sq)}))
		Expect(w.stdout.String()).To(HavePrefix("laserfix 0.1.0 · grid.svg → fixed.svg · "))
		Expect(w.files["grid.svg"]).To(Equal(doc(sq + sq)))
		Expect(w.opened).To(BeZero(), "every file opened is closed")
	})

	It("takes -o before the file too", func() {
		// Act
		code := w.run("-o", "fixed.svg", "grid.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitClean))
		Expect(w.written).To(HaveKey("fixed.svg"))
	})

	It("reads standard input and writes standard output, with the report on standard error", func() {
		// Arrange
		w.stdin = doc(sq + sq)

		// Act
		code := w.run("-", "-o", "-")

		// Assert
		Expect(code).To(Equal(cli.ExitClean))
		Expect(w.stdout.String()).To(Equal(doc(sq)))
		Expect(w.stderr.String()).To(HavePrefix("laserfix 0.1.0 · standard input → standard output · "))
	})

	It("writes the report as JSON with --json", func() {
		// Act
		code := w.run("--json", "grid.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitChanges))
		Expect(w.stdout.String()).To(ContainSubstring(`"schema": "kinglet.report/v1"`))
	})

	It("uses --same as the same-line distance", func() {
		// Arrange: the second square is 0.03 mm off, beyond the default 0.02 mm.
		w.files["near.svg"] = doc(sq + `<rect x="10.03" y="10" width="10" height="10" fill="none" stroke="#000"/>`)

		// Act
		code := w.run("--same", "0.05mm", "--json", "near.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitChanges))
		Expect(w.stdout.String()).To(ContainSubstring(`"same_mm": 0.05`))
	})

	It("exits 1 when the written copy came out worse", func() {
		// Arrange
		worse := app.Report{Changes: []app.Change{{Kind: app.Copy}}, After: app.Measure{Crossings: 1}}

		// Act
		code := w.runWith(func(d *cli.Deps) { d.Fixer = fixed(worse) }, "grid.svg", "-o", "fixed.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitChanges))
		Expect(w.written).To(HaveKey("fixed.svg"))
	})

	It("uses --join and --report-gaps for split outlines", func() {
		// Arrange: two lines 0.1 mm apart end to end.
		w.files["gap.svg"] = doc(`<path d="M10 10 H20" fill="none" stroke="#000"/><path d="M20.1 10 H30" fill="none" stroke="#000"/>`)

		// Act
		joined := w.run("--join", "0.2mm", "--json", "gap.svg")
		joinedOut := w.stdout.String()
		w.stdout.Reset()
		quiet := w.run("--report-gaps", "0.05mm", "gap.svg")

		// Assert
		Expect(joined).To(Equal(cli.ExitChanges))
		Expect(joinedOut).To(ContainSubstring(`"kind": "join"`))
		Expect(quiet).To(Equal(cli.ExitClean))
	})

	It("exits 1 when it leaves something it reports, though nothing needs fixing", func() {
		// Arrange
		w.files["gap.svg"] = doc(`<path d="M10 10 H20" fill="none" stroke="#000"/><path d="M20.5 10 H30" fill="none" stroke="#000"/>`)

		// Act
		code := w.run("gap.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitChanges))
		Expect(w.stdout.String()).To(ContainSubstring("Left as they are"))
	})

	It("shows the version", func() {
		// Act
		code := w.run("--version")

		// Assert
		Expect(code).To(Equal(cli.ExitClean))
		Expect(w.stdout.String()).To(Equal("laserfix 0.1.0\n"))
	})

	It("shows help with examples", func() {
		// Act
		code := w.run("--help")

		// Assert
		Expect(code).To(Equal(cli.ExitClean))
		Expect(w.stdout.String()).To(ContainSubstring("Usage:"))
		Expect(w.stdout.String()).To(ContainSubstring("Examples:"))
	})
})

// fixed is a test double for the Fixer port that always reports r.
type fixed app.Report

func (f fixed) Fix([]byte, app.Settings) (app.Report, error) { return app.Report(f), nil }
