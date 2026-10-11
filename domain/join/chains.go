package join

import (
	"github.com/kinglet-dev/core/domain/design"
	"github.com/kinglet-dev/core/domain/geom"
)

// walked is a chain as pieces entered by an end each: piece n's start is
// end 2n and its last point end 2n+1.
type walked struct {
	links  []int // the end each piece is entered by
	closed bool
}

// chains walks every chain of paired ends, in the order of each chain's
// earliest piece.
func chains(pieces int, partner map[int]int) []walked {
	seen := make([]bool, pieces)
	var out []walked
	for n := range pieces {
		_, a := partner[2*n]
		_, b := partner[2*n+1]
		if seen[n] || !(a || b) {
			continue
		}
		in := head(partner, n)
		w := walked{}
		for {
			seen[in/2] = true
			w.links = append(w.links, in)
			next, ok := partner[in^1]
			if !ok {
				break
			}
			if next == w.links[0] {
				w.closed = true
				break
			}
			in = next
		}
		out = append(out, w)
	}
	return out
}

// head is the end piece n's chain starts at: walking back from n's start
// to the chain's first piece, or n's own start when the chain is a loop.
func head(partner map[int]int, n int) int {
	in := 2 * n
	for {
		p, ok := partner[in]
		if !ok {
			return in
		}
		if p/2 == n {
			return 2 * n
		}
		in = p ^ 1
	}
}

// chain names the walked chain's pieces.
func (w walked) chain(pieces []Piece, ends []end) Chain {
	c := Chain{Closed: w.closed}
	for _, in := range w.links {
		c.Links = append(c.Links, Link{Piece: pieces[in/2], Reversed: ends[in].last})
	}
	return c
}

// First is the chain's earliest piece in the document: the joined path
// takes its place.
func (c Chain) First() Piece {
	first := c.Links[0].Piece
	for _, l := range c.Links {
		if l.Piece.Shape < first.Shape || (l.Piece.Shape == first.Shape && l.Piece.Sub < first.Sub) {
			first = l.Piece
		}
	}
	return first
}

// Joined draws a chain as one subpath in the units of its first piece's
// shape, bridging each small gap with a straight line.
func Joined(shapes []design.Shape, c Chain) geom.Subpath {
	toFirst, _ := shapes[c.First().Shape].ToMM.Inverse()
	var out geom.Subpath
	for i, l := range c.Links {
		s := shapes[l.Piece.Shape]
		sub := s.Local.Subpaths[l.Piece.Sub]
		if l.Reversed {
			sub = sub.Reverse()
		}
		sub = geom.Path{Subpaths: []geom.Subpath{sub}}.Transform(s.ToMM.Then(toFirst)).Subpaths[0]
		if i == 0 {
			out.Start = sub.Start
		} else if at := out.Segments[len(out.Segments)-1].To; at != sub.Start {
			out.Segments = append(out.Segments, geom.LineTo(sub.Start))
		}
		out.Segments = append(out.Segments, sub.Segments...)
	}
	out.Closed = c.Closed
	return out
}
