package app_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/core/domain/crossing"
	"github.com/kinglet-dev/laserfix/adapters/svg"
	"github.com/kinglet-dev/laserfix/app"
)

// remover is a test double for the Remover port that returns what it is given.
type remover struct {
	out []byte
	err error
}

func (r remover) Remove([]byte, []int) ([]byte, error) { return r.out, r.err }

var _ = Describe("Fix stops with an error", func() {
	It("when the document can't be read", func() {
		// Arrange: a remover that fails shows nothing else was tried.
		fixer := app.NewFixer(svg.Document{}, remover{err: errors.New("failed")})

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
		fixer := app.NewFixer(svg.Document{}, remover{err: failed})

		// Act
		_, err := fixer.Fix(doc(sq("10")+sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(failed))
		Expect(err.Error()).To(HavePrefix("changing the document: "))
	})

	It("when the changed document can't be read", func() {
		// Arrange
		fixer := app.NewFixer(svg.Document{}, remover{out: []byte("<html/>")})

		// Act
		_, err := fixer.Fix(doc(sq("10")), app.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(reader.ErrNotSVG))
		Expect(err.Error()).To(HavePrefix("reading the changed document: "))
	})

	It("when the changed document has too much line to measure", func() {
		// Arrange: a remover that, against its contract, adds line.
		fixer := app.NewFixer(svg.Document{}, remover{out: doc(sq("10") + sq("40"))})
		s := app.DefaultSettings
		s.Crossings.MaxSegments = 5

		// Act
		_, err := fixer.Fix(doc(sq("10")), s)

		// Assert
		Expect(err).To(MatchError(crossing.ErrTooComplex))
	})
})
