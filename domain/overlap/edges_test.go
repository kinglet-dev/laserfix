package overlap_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

var _ = Describe("Find draws the edges of each trim precisely", func() {
	It("stops a trim exactly at a corner that turns off the shared line at a shallow angle", func() {
		// Arrange: the later line follows the earlier one back to x = 0, then turns 60° away.
		shapes := []design.Shape{line(pt(-10, 0), pt(10, 0)), line(pt(10, 0), pt(0, 0), pt(-2.5, 4.330127))}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(0, 0), pt(-2.5, 4.330127)}}))
	})

	It("leaves line a little over the same-line distance away", func() {
		// Act
		trims := find([]design.Shape{line(pt(0, 0), pt(100, 0)), line(pt(0, 0.03), pt(100, 0.03))}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("finds earlier line drawn right to left, just across a grid cell boundary", func() {
		// Act
		trims := find([]design.Shape{line(pt(100, -0.005), pt(0, -0.005)), line(pt(0, 0.005), pt(100, 0.005))}, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Keep).To(BeEmpty())
	})

	It("adds no extra side to a closed shape drawn back to its start", func() {
		// Arrange
		later := shape(poly(true, pt(20, 0), pt(20, 10), pt(10, 10), pt(10, 0), pt(20, 0)))
		shapes := []design.Shape{square(0, 0), later}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(10, 0), pt(20, 0), pt(20, 10), pt(10, 10)}}))
	})

	It("leaves a stretch shorter than five times the same-line distance", func() {
		// Act: the later line touches the earlier one along 0.05 mm.
		trims := find([]design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(5, 5), pt(5, 0), pt(5.05, 0), pt(5.05, 5))}, nil)

		// Assert
		Expect(trims).To(BeEmpty())
	})

	It("keeps the middle of an open line whose two ends lie on earlier line", func() {
		// Arrange
		shapes := []design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(20, 0), pt(30, 0)), line(pt(0, 0), pt(30, 0))}

		// Act
		trims := find(shapes, nil)

		// Assert: it starts where the first stretch cut away starts.
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].At).To(Equal(pt(0, 0)))
		Expect(lines(kept(shapes, trims, 2))).To(approx([][]geom.Point{{pt(10, 0), pt(20, 0)}}))
	})

	It("compares later shapes with what is left of earlier ones, not with what was cut", func() {
		// Arrange: the second line is cut where it runs 0.015 mm from the first; the third
		// runs 0.015 mm from that cut stretch but 0.03 mm from the first line.
		shapes := []design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(0, 0.015), pt(20, 0.015)), line(pt(0, 0.03), pt(10, 0.03))}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Shape).To(Equal(1))
	})

	It("keeps the rest of a closed shape when only the middle of one side is shared", func() {
		// Arrange
		shapes := []design.Shape{line(pt(3, 0), pt(7, 0)), square(0, 0)}

		// Act
		trims := find(shapes, nil)

		// Assert
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(7, 0), pt(10, 0), pt(10, 10), pt(0, 10), pt(0, 0), pt(3, 0)}}))
	})
})
