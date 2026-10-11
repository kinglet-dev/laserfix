package join_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/domain/join"
)

var _ = Describe("Find joins pieces whose ends meet", func() {
	It("joins a square drawn as four separate sides into one closed outline", func() {
		// Arrange: the sides come in any order and direction.
		shapes := []design.Shape{
			line(pt(0, 0), pt(10, 0)), line(pt(10, 10), pt(0, 10)),
			line(pt(10, 0), pt(10, 10)), line(pt(0, 0), pt(0, 10)),
		}

		// Act
		r := join.Find(shapes, nil, params)

		// Assert
		Expect(r.Chains).To(HaveLen(1))
		c := r.Chains[0]
		Expect(c.Closed).To(BeTrue())
		Expect(c.Links).To(HaveLen(4))
		Expect(c.Links[0].Piece).To(Equal(join.Piece{Shape: 0, Sub: 0}))
		got := join.Joined(shapes, c)
		Expect(got.Closed).To(BeTrue())
		Expect(points(got)).To(Equal([]geom.Point{pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10), pt(0, 0)}))
	})

	It("doesn't bridge a gap too small to write down", func() {
		// Arrange
		shapes := []design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(10+1e-9, 0), pt(20, 0))}

		// Act
		r := join.Find(shapes, nil, params)

		// Assert
		Expect(join.Joined(shapes, r.Chains[0]).Segments).To(HaveLen(2))
	})

	It("joins pieces into an open chain, keeping the earliest piece's direction", func() {
		// Arrange
		shapes := []design.Shape{line(pt(10, 0), pt(20, 0)), line(pt(0, 0), pt(10, 0.03))}

		// Act
		r := join.Find(shapes, nil, params)

		// Assert: the 0.03 mm gap is closed with a straight line.
		Expect(r.Chains).To(HaveLen(1))
		Expect(r.Chains[0].Closed).To(BeFalse())
		got := join.Joined(shapes, r.Chains[0])
		Expect(points(got)).To(Equal([]geom.Point{pt(0, 0), pt(10, 0.03), pt(10, 0), pt(20, 0)}))
	})

	It("closes a single outline whose ends are a small gap apart", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10), pt(0, 0.02))}, nil, params)

		// Assert
		Expect(r.Chains).To(HaveLen(1))
		Expect(r.Chains[0].Closed).To(BeTrue())
	})

	It("closes an outline that comes after other pieces", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(50, 50), pt(60, 50)), line(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10), pt(0, 0.02))}, nil, params)

		// Assert
		Expect(r.Chains).To(HaveLen(1))
		Expect(r.Chains[0].Links[0].Piece).To(Equal(join.Piece{Shape: 1, Sub: 0}))
	})

	It("leaves an outline whose ends already meet exactly, since nothing would change", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 0))}, nil, params)

		// Assert
		Expect(r.Chains).To(BeEmpty())
	})

	It("joins pieces from separate subpaths of one element, in place of the earliest", func() {
		// Act: the chain starts with the second subpath.
		r := join.Find([]design.Shape{shape(sub(false, pt(10, 0), pt(10, 10)), sub(false, pt(0, 0), pt(10, 0)))}, nil, params)

		// Assert
		Expect(r.Chains).To(HaveLen(1))
		Expect(r.Chains[0].Links[0].Piece).To(Equal(join.Piece{Shape: 0, Sub: 1}))
		Expect(r.Chains[0].First()).To(Equal(join.Piece{Shape: 0, Sub: 0}))
	})

	It("draws the joined path in the units of the earliest piece's element", func() {
		// Arrange: the first element is drawn in half millimeters.
		half := line(pt(20, 0), pt(40, 0))
		half.ToMM = geom.Matrix{A: 0.5, D: 0.5}
		half.Path = half.Local.Transform(half.ToMM)
		shapes := []design.Shape{half, line(pt(0, 0), pt(10, 0))}

		// Act
		r := join.Find(shapes, nil, params)

		// Assert
		Expect(r.Chains[0].First()).To(Equal(join.Piece{Shape: 0, Sub: 0}))
		Expect(points(join.Joined(shapes, r.Chains[0]))).To(Equal([]geom.Point{pt(0, 0), pt(20, 0), pt(40, 0)}))
	})

	It("ignores a subpath that is only a point", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0)), shape(geom.Subpath{Start: pt(10, 0)})}, nil, params)

		// Assert
		Expect(r.Chains).To(BeEmpty())
	})

	DescribeTable("joins ends either side of a grid cell boundary",
		func(a, b geom.Point) {
			// Act
			r := join.Find([]design.Shape{line(pt(-5, -5), a), line(b, pt(15, 15))}, nil, params)

			// Assert
			Expect(r.Chains).To(HaveLen(1))
		},
		Entry("across x", pt(9.99, 0), pt(10.01, 0)),
		Entry("across y", pt(0, 9.99), pt(0, 10.01)),
	)
})

var _ = Describe("Find reports what it won't join", func() {
	It("reports three ends meeting at one point and joins none of them", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(10, 0), pt(20, 0)), line(pt(10, 0), pt(10, 10))}, nil, params)

		// Assert
		Expect(r.Chains).To(BeEmpty())
		Expect(r.Junctions).To(Equal([]join.Junction{{At: pt(10, 0), Ends: 3}}))
	})

	It("reports a gap too wide to close, once", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(10.6, 0), pt(20, 0))}, nil, params)

		// Assert
		Expect(r.Chains).To(BeEmpty())
		Expect(r.Gaps).To(HaveLen(1))
		Expect(r.Gaps[0].At).To(Equal(pt(10, 0)))
		Expect(r.Gaps[0].Size).To(BeNumerically("~", 0.6, 1e-9))
	})

	It("doesn't join ends just over the join distance apart, in the next grid cell", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(10.07, 0), pt(20, 0))}, nil, params)

		// Assert
		Expect(r.Chains).To(BeEmpty())
		Expect(r.Gaps).To(HaveLen(1))
	})

	It("joins none of three ends in a row where only the middle one is near both", func() {
		// Act
		r := join.Find([]design.Shape{
			line(pt(0, 0), pt(10, 0)), line(pt(10.04, 0), pt(10.04, 10)), line(pt(10.08, 0), pt(20, 0)),
		}, nil, params)

		// Assert
		Expect(r.Chains).To(BeEmpty())
		Expect(r.Junctions).To(HaveLen(1))
	})

	It("reports gaps only between two ends that are both left open", func() {
		// Arrange: the first two lines join; the third ends 0.5 mm from that joint.
		shapes := []design.Shape{line(pt(10, 0), pt(20, 0)), line(pt(10.02, 0), pt(10.02, 10)), line(pt(0, 0), pt(9.5, 0))}

		// Act
		r := join.Find(shapes, nil, params)

		// Assert
		Expect(r.Chains).To(HaveLen(1))
		Expect(r.Gaps).To(BeEmpty())
	})

	It("doesn't report the ends of parallel lines as a gap, since neither points at the other", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(10, 25), pt(110, 25)), line(pt(10, 25.3), pt(110, 25.3))}, nil, params)

		// Assert
		Expect(r.Gaps).To(BeEmpty())
	})

	It("reports a gap between curves whose ends point at each other", func() {
		// Arrange: two arcs of a circle, each stopping 0.2 mm short of the other.
		// The first control point alone would make it run sideways out of its end.
		a := shape(geom.Subpath{Start: pt(0, 10), Segments: []geom.Segment{geom.CubicTo(pt(30, 0.1), pt(9.95, 5.5), pt(9.95, 0.1))}})
		b := shape(geom.Subpath{Start: pt(9.95, -0.1), Segments: []geom.Segment{geom.CubicTo(pt(9.95, -5.5), pt(5.5, -10), pt(0, -10))}})

		// Act
		r := join.Find([]design.Shape{a, b}, nil, params)

		// Assert
		Expect(r.Gaps).To(HaveLen(1))
	})

	It("finds which way a curve runs out of its end when a control point sits on that end", func() {
		// Arrange
		a := shape(geom.Subpath{Start: pt(0, 10), Segments: []geom.Segment{geom.CubicTo(pt(9.95, 10), pt(9.95, 0.1), pt(9.95, 0.1))}})
		b := shape(geom.Subpath{Start: pt(9.95, -0.1), Segments: []geom.Segment{geom.CubicTo(pt(9.95, -0.1), pt(9.95, -10), pt(0, -10))}})

		// Act
		r := join.Find([]design.Shape{a, b}, nil, params)

		// Assert
		Expect(r.Gaps).To(HaveLen(1))
	})

	It("reports a gap at the end of a line of several pieces, from the way its last piece runs", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(-100, 0.3), pt(10, 20), pt(10, 0.3)), line(pt(10, -0.3), pt(10, -10))}, nil, params)

		// Assert
		Expect(r.Gaps).To(HaveLen(1))
	})

	DescribeTable("doesn't report a line ending beside the end of a line running across it",
		func(across, along design.Shape, acrossFirst bool) {
			// Arrange
			shapes := []design.Shape{along, across}
			if acrossFirst {
				shapes = []design.Shape{across, along}
			}

			// Act
			r := join.Find(shapes, nil, params)

			// Assert
			Expect(r.Gaps).To(BeEmpty())
		},
		Entry("listed second", line(pt(10.5, -5), pt(10.5, 0)), line(pt(0, 0), pt(10, 0)), false),
		Entry("listed first", line(pt(10.5, -5), pt(10.5, 0)), line(pt(0, 0), pt(10, 0)), true),
	)

	It("doesn't report ends further apart than the gap limit", func() {
		// Act
		r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0)), line(pt(12, 0), pt(20, 0))}, nil, params)

		// Assert
		Expect(r.Gaps).To(BeEmpty())
	})
})

var _ = Describe("Find leaves alone what it mustn't join", func() {
	DescribeTable("leaves pieces that",
		func(change func(*design.Shape), skip []bool) {
			// Arrange
			other := line(pt(10, 0), pt(20, 0))
			change(&other)

			// Act
			r := join.Find([]design.Shape{line(pt(0, 0), pt(10, 0)), other}, skip, params)

			// Assert
			Expect(r.Chains).To(BeEmpty())
		},
		Entry("have a different stroke", func(s *design.Shape) { s.Stroke = "#ff0000" }, nil),
		Entry("are filled", func(s *design.Shape) { s.Fill = "#000000" }, nil),
		Entry("are skipped (copies removed, shapes trimmed)", func(*design.Shape) {}, []bool{false, true}),
		Entry("are closed already", func(s *design.Shape) { s.Local.Subpaths[0].Closed = true; s.Path = s.Local }, nil),
		Entry("can't be drawn in another's units", func(s *design.Shape) { s.ToMM = geom.Matrix{} }, nil),
	)
})
