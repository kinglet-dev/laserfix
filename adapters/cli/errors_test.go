package cli_test

import (
	"strings"
	"testing/iotest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/laserfix/adapters/cli"
)

var _ = Describe("Run stops and exits 2", func() {
	var w *world

	BeforeEach(func() {
		// Arrange (shared)
		w = newWorld()
		w.files["grid.svg"] = doc(sq + sq)
	})

	DescribeTable("for arguments it can't use, saying how to fix them",
		func(want string, args ...string) {
			// Act
			code := w.run(args...)

			// Assert
			Expect(code).To(Equal(cli.ExitError))
			Expect(w.stderr.String()).To(Equal("laserfix: " + want + "\nRun laserfix --help for usage.\n"))
			Expect(w.stdout.String()).To(BeEmpty())
		},
		Entry("no file", "give one SVG file to fix, or - to read standard input"),
		Entry("two files", "give one SVG file to fix, or - to read standard input", "a.svg", "b.svg"),
		Entry("an unknown flag", "flag provided but not defined: -fast", "--fast", "grid.svg"),
		Entry("a length without a unit", "--same: length has no unit: add one of mm, cm, in, pt or px, for example 0.25mm", "--same", "0.02", "grid.svg"),
		Entry("a same-line distance of zero", "--same must be more than zero", "--same", "0mm", "grid.svg"),
	)

	It("when the file can't be opened", func() {
		// Act
		code := w.run("missing.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(Equal("laserfix: can't open missing.svg: file does not exist\n"))
	})

	It("when the file isn't an SVG", func() {
		// Arrange
		w.files["page.html"] = "<html/>"

		// Act
		code := w.run("page.html")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(HavePrefix("laserfix: page.html: the file is not an SVG"))
	})

	It("when the output file already exists, writing nothing", func() {
		// Act
		code := w.run("grid.svg", "-o", "grid.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(HavePrefix("laserfix: can't write grid.svg: it already exists"))
		Expect(w.stdout.String()).To(BeEmpty())
	})

	It("when the output can't be written", func() {
		// Act
		code := w.run("grid.svg", "-o", "broken.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(Equal("laserfix: can't write broken.svg: disk full\n"))
	})

	It("when standard input is larger than any SVG it reads", func() {
		// Arrange
		w.stdin = doc(sq) + strings.Repeat(" ", int(reader.DefaultLimits.Bytes))

		// Act
		code := w.run("-")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(HavePrefix("laserfix: standard input: the file is too large"))
	})

	It("when standard input can't be read", func() {
		// Act
		code := w.runWith(func(d *cli.Deps) { d.Stdin = iotest.ErrReader(errBroken) }, "-")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(Equal("laserfix: standard input: disk full\n"))
	})

	It("when standard output can't take the repaired copy", func() {
		// Act
		code := w.runWith(func(d *cli.Deps) { d.Stdout = failingWriter{} }, "grid.svg", "-o", "-")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(Equal("laserfix: can't write standard output: disk full\n"))
	})

	It("when the report can't be written", func() {
		// Act
		code := w.runWith(func(d *cli.Deps) { d.Stdout = failingWriter{} }, "grid.svg")

		// Assert
		Expect(code).To(Equal(cli.ExitError))
		Expect(w.stderr.String()).To(Equal("laserfix: can't write the report: disk full\n"))
	})
})

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errBroken }
