package copies_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/laserfix/domain/copies"
)

var _ = Describe("Find compares whole lines within the same-line distance", func() {
	It("finds a copy moved by less than the same-line distance", func() {
		// Act
		got := copies.Find([]design.Shape{square(0, 0), square(0.01, 0)}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})

	It("doesn't count a shape moved by more than the same-line distance", func() {
		// Act
		got := copies.Find([]design.Shape{square(0, 0), square(0.05, 0)}, params)

		// Assert
		Expect(got).To(BeEmpty())
	})

	It("doesn't count a shorter line lying on a longer one as a copy", func() {
		// Arrange
		long := hairline("#000000", false, pt(0, 0), pt(10, 0))
		start := hairline("#000000", false, pt(0, 0), pt(5, 0))
		end := hairline("#000000", false, pt(5, 0), pt(10, 0))

		// Act
		got := [][]copies.Copy{
			copies.Find([]design.Shape{long, start}, params), copies.Find([]design.Shape{start, long}, params),
			copies.Find([]design.Shape{long, end}, params), copies.Find([]design.Shape{end, long}, params),
		}

		// Assert
		Expect(got).To(HaveEach(BeEmpty()))
	})

	It("finds every later copy as a copy of the first", func() {
		// Act
		got := copies.Find([]design.Shape{square(0, 0), square(0, 0), square(0, 0)}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}, {Shape: 2, Of: 0}}))
	})

	It("finds a copy drawn with an extra point along one side", func() {
		// Arrange
		extra := hairline("#000000", true, pt(0, 0), pt(5, 0), pt(10, 0), pt(10, 10), pt(0, 10))

		// Act
		got := copies.Find([]design.Shape{square(0, 0), extra}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})

	It("doesn't count lines that only share their corners", func() {
		// Arrange: a bow-tie through the square's four corners.
		bowtie := hairline("#000000", true, pt(0, 0), pt(10, 10), pt(10, 0), pt(0, 10))

		// Act
		got := copies.Find([]design.Shape{square(0, 0), bowtie}, params)

		// Assert
		Expect(got).To(BeEmpty())
	})

	It("finds a copy drawn with a repeated point", func() {
		// Arrange
		repeated := hairline("#000000", true, pt(0, 0), pt(10, 0), pt(10, 0), pt(10, 10), pt(0, 10))

		// Act
		got := copies.Find([]design.Shape{square(0, 0), repeated}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})

	It("finds a copy of a dot, which burns as a spot", func() {
		// Arrange
		dot := hairline("#000000", false, pt(5, 5), pt(5, 5))

		// Act
		got := copies.Find([]design.Shape{dot, dot}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})

	It("samples between corners finely enough to see a difference just over the same-line distance", func() {
		// Arrange: a bow-tie through the corners of a 0.0415 mm square. The
		// middle of its crossing is 0.02075 mm from the square's sides, just
		// over the 0.02 mm same-line distance; everything else is within it.
		const side = 0.0415
		tiny := hairline("#000000", true, pt(0, 0), pt(side, 0), pt(side, side), pt(0, side))
		bowtie := hairline("#000000", true, pt(0, 0), pt(side, side), pt(side, 0), pt(0, side))

		// Act
		got := copies.Find([]design.Shape{tiny, bowtie}, params)

		// Assert
		Expect(got).To(BeEmpty())
	})

	It("compares each shape only with shapes that stay, not with copies already found", func() {
		// Arrange: each square is 0.015 mm from the one before. The third is
		// 0.03 mm from the first, which stays, so it is not a copy of anything
		// that stays.
		shapes := []design.Shape{square(0, 0), square(0.015, 0), square(0.03, 0)}

		// Act
		got := copies.Find(shapes, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})

	It("finds a copy of a filled shape whose outline is left open, because the fill closes it", func() {
		// Arrange
		closed := design.Shape{Path: outline(true, pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10)), Filled: true, Fill: "#000000", Stroke: "none"}
		open := closed
		open.Path = outline(false, pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10))

		// Act
		got := copies.Find([]design.Shape{closed, open}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})

	It("ignores shapes that draw no line, such as an empty path", func() {
		// Arrange
		empty := design.Shape{Fill: "none", Stroke: "#000000"}

		// Act
		got := copies.Find([]design.Shape{empty, empty}, params)

		// Assert
		Expect(got).To(BeEmpty())
	})
})
