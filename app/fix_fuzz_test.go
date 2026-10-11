package app_test

import (
	"testing"

	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

// FuzzFix feeds Fix arbitrary documents. It must never panic, and fixing
// is idempotent: a fixed document has nothing left to fix.
func FuzzFix(f *testing.F) {
	// Arrange: seeds from the unit tests.
	f.Add(string(doc(sq("10") + sq("10"))))
	f.Add(string(doc(sq("10") + sq("40") + sq("10.01"))))
	f.Add(string(doc(`<g transform="scale(2)"><path d="M0 5H10M5 0V10" fill="none" stroke="#000"/></g><path d="M0 10H20M10 0V20" fill="none" stroke="#000"/>`)))
	f.Add(string(doc(`<path d=""/><path d="" stroke="#000"/>`)))
	f.Add(string(doc(sq("10") + sq("20") + hair("M20 20 V30") + hair("M0 0 H10 M50 0 H60") + hair("M10 0 V10"))))
	f.Add(string(doc(hair("M10 10 H20 V20 H10 V10.02") + hair("M10 50 H50") + hair("M50 50 H90") + hair("M50 50 V90"))))
	f.Add(string(doc(`<circle cx="50" cy="50" r="20" fill="none" stroke="#000"/><path d="M70 50 A20 20 0 0 1 30 50" fill="none" stroke="#000"/>`)))
	f.Add("<html/>")
	fixer := app.NewFixer(svg.Document{}, svg.Document{})
	f.Fuzz(func(t *testing.T, src string) {
		// Act
		r, err := fixer.Fix([]byte(src), app.DefaultSettings)
		if err != nil {
			return
		}
		again, err := fixer.Fix(r.Output, app.DefaultSettings)

		// Assert
		if err != nil {
			t.Fatalf("the fixed document can't be fixed again: %v", err)
		}
		if len(again.Changes) > 0 {
			t.Fatalf("the fixed document still needs %d changes", len(again.Changes))
		}
	})
}
