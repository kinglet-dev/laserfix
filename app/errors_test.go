package app_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/core/domain/crossing"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
	"github.com/kinglet-dev/laserfix/domain/overlap"
)

// editor is a test double for the Editor port that returns what it is given.
type editor struct {
	out []byte
	err error
}

func (e editor) Edit([]byte, []app.Edit) ([]byte, error) { return e.out, e.err }

var _ = Describe("Fix stops with an error", func() {
	It("when the document can't be read", func() {
		// Arrange: an editor that fails shows nothing else was tried.
		fixer := app.NewFixer(svg.Document{}, editor{err: errors.New("failed")})

		// Act
		_, err := fixer.Fix([]byte("<html/>"), app.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(reader.ErrNotSVG))
	})

	It("when the design has too much line to measure", func() {
		// Arrange
		fixer := app.NewFixer(svg.Document{}, svg.Document{})
		s := app.DefaultSettings
		s.Crossings.MaxSegments = 5 // the copy brings 8 pieces down to 4

		// Act
		_, err := fixer.Fix(doc(sq("10")+sq("10")), s)

		// Assert
		Expect(err).To(MatchError(crossing.ErrTooComplex))
	})

	It("when the document can't be changed", func() {
		// Arrange
		failed := errors.New("failed")
		fixer := app.NewFixer(svg.Document{}, editor{err: failed})

		// Act
		_, err := fixer.Fix(doc(sq("10")+sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(failed))
		Expect(err.Error()).To(HavePrefix("changing the document: "))
	})

	It("when the changed document can't be read", func() {
		// Arrange
		fixer := app.NewFixer(svg.Document{}, editor{out: []byte("<html/>")})

		// Act: the copy needs a change.
		_, err := fixer.Fix(doc(sq("10")+sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(reader.ErrNotSVG))
		Expect(err.Error()).To(HavePrefix("reading the changed document: "))
	})

	It("when the changed document has too much line to measure", func() {
		// Arrange: an editor that, against its contract, adds line to the one shape it keeps.
		fixer := app.NewFixer(svg.Document{}, editor{out: doc(hair("M0 0 H10 V10 H0 Z M50 50 H60 V60 H50"))})
		s := app.DefaultSettings
		s.Crossings.MaxSegments = 5

		// Act: the second line copies the first.
		_, err := fixer.Fix(doc(hair("M0 0 H10")+hair("M0 0 H10")), s)

		// Assert
		Expect(err).To(MatchError(crossing.ErrTooComplex))
	})

	It("when the design has too much line to search for shared line", func() {
		// Arrange
		fixer := app.NewFixer(svg.Document{}, svg.Document{})
		s := app.DefaultSettings
		s.MaxShared = 5

		// Act
		_, err := fixer.Fix(doc(sq("10")), s)

		// Assert
		Expect(err).To(MatchError(overlap.ErrTooComplex))
	})

	It("when the changed document's shapes no longer match the ones kept", func() {
		// Arrange: an editor that, against its contract, keeps both shapes.
		fixer := app.NewFixer(svg.Document{}, editor{out: doc(sq("10") + sq("10"))})

		// Act: the copy should go, leaving one shape.
		_, err := fixer.Fix(doc(sq("10")+sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(HavePrefix("reading the changed document: its shapes no longer match")))
	})
})
