package copies_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/laserfix/domain/copies"
)

var params = copies.Params{Same: 0.02, Flatten: 0.005}

var _ = Describe("Find finds shapes that copy an earlier shape", func() {
	It("finds nothing among different shapes", func() {
		// Act
		got := copies.Find([]design.Shape{square(0, 0), square(20, 0)}, params)

		// Assert
		Expect(got).To(BeEmpty())
	})

	It("finds the later of two identical shapes", func() {
		// Act
		got := copies.Find([]design.Shape{square(0, 0), square(0, 0)}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})
	It("finds a copy drawn from another start point, in the other direction", func() {
		// Arrange
		again := hairline("#000000", true, pt(10, 10), pt(10, 0), pt(0, 0), pt(0, 10))

		// Act
		got := copies.Find([]design.Shape{square(0, 0), again}, params)

		// Assert
		Expect(got).To(Equal([]copies.Copy{{Shape: 1, Of: 0}}))
	})
	DescribeTable("treats the same line in other paint as a different shape",
		func(other design.Shape) {
			// Act
			got := copies.Find([]design.Shape{square(0, 0), other}, params)

			// Assert
			Expect(got).To(BeEmpty())
		},
		Entry("another stroke color", hairline("#ff0000", true, pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10))),
		Entry("a fill", func() design.Shape { s := square(0, 0); s.Fill = "black"; return s }()),
	)
})
