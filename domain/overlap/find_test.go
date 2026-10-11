package overlap_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/overlap"
)

var _ = Describe("Find trims line that lies on earlier line from the later shape", func() {
	It("trims the edge two squares share from the later square", func() {
		// Arrange
		shapes := []design.Shape{square(0, 0), square(10, 0)}

		// Act
		trims, err := overlap.Find(shapes, nil, params)
		Expect(err).NotTo(HaveOccurred())

		// Assert
		Expect(trims).To(HaveLen(1))
		Expect(trims[0].Shape).To(Equal(1))
		Expect(trims[0].Removed).To(BeNumerically("~", 10, 1e-6))
		Expect(trims[0].At).To(Equal(geom.Point{X: 10, Y: 10}))
		Expect(lines(kept(shapes, trims, 1))).To(approx([][]geom.Point{{pt(10, 0), pt(20, 0), pt(20, 10), pt(10, 10)}}))
	})

	It("leaves shapes that only touch at a corner or cross alone", func() {
		// Arrange: a square, one touching its corner, and a line crossing it.
		shapes := []design.Shape{square(0, 0), square(10, 10), shape(poly(false, pt(-5, 5), pt(15, 5)))}

		// Act
		trims, err := overlap.Find(shapes, nil, params)
		Expect(err).NotTo(HaveOccurred())

		// Assert
		Expect(trims).To(BeEmpty())
	})
})
