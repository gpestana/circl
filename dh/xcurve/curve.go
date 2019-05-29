// +build amd64

package xcurve

import (
	fp255 "github.com/cloudflare/circl/math/fp25519"
	"github.com/cloudflare/circl/math/fp448"
)

type curve struct {
	size       int                        // Size in bytes of prime field elements
	n          int                        // Size in bits of the prime field
	h          int                        // Cofactor of the elliptic curve group
	a24        int                        // (A+2)/4
	table      []byte                     // Precomputed multiples of the generator point
	pointS     []byte                     // The x-coord of a point of order four
	pointGS    []byte                     // The x-coord of the generator point minus the point S
	ladderStep func([]byte, uint)         // Differential addition and doubling
	double     func([]byte)               // Point doubling
	difAdd     func([]byte, []byte, uint) // Differential point addition
}

var c255, c448 *curve

func init() {
	var ( // Coordinates of points according to ia.cr/2017/264
		pointS255  = fp255.Elt{1}
		pointGS255 = fp255.Elt{0xbd, 0xaa, 0x2f, 0xc8, 0xfe, 0xe1, 0x94, 0x7e, 0xf8, 0xed, 0xb2, 0x14, 0xae, 0x95, 0xf0, 0xbb, 0xe2, 0x48, 0x5d, 0x23, 0xb9, 0xa0, 0xc7, 0xad, 0x34, 0xab, 0x7c, 0xe2, 0xee, 0xcd, 0xae, 0x1e}
		pointS448  = fp448.Elt{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		pointGS448 = fp448.Elt{0x20, 0x27, 0x9d, 0xc9, 0x7d, 0x19, 0xb1, 0xac, 0xf8, 0xba, 0x69, 0x1c, 0xff, 0x33, 0xac, 0x23, 0x51, 0x1b, 0xce, 0x3a, 0x64, 0x65, 0xbd, 0xf1, 0x23, 0xf8, 0xc1, 0x84, 0x9d, 0x45, 0x54, 0x29, 0x67, 0xb9, 0x81, 0x1c, 0x3, 0xd1, 0xcd, 0xda, 0x7b, 0xeb, 0xff, 0x1a, 0x88, 0x3, 0xcf, 0x3a, 0x42, 0x44, 0x32, 0x1, 0x25, 0xb7, 0xfa, 0xf0}
	)

	c255 = &curve{
		size:       fp255.Size,
		n:          255,
		h:          3,
		a24:        121666,
		table:      tableGenerator255[:],
		pointS:     pointS255[:],
		pointGS:    pointGS255[:],
		ladderStep: ladderStep255,
		double:     double255,
		difAdd:     difAdd255,
	}

	c448 = &curve{
		size:       fp448.Size,
		n:          448,
		h:          2,
		a24:        39082,
		table:      tableGenerator448[:],
		pointS:     pointS448[:],
		pointGS:    pointGS448[:],
		ladderStep: ladderStep448,
		double:     double448,
		difAdd:     difAdd448,
	}
}

// ladderJoye calculates a fixed-point multiplication with the generator point.
// The algorithm is the right-to-left Joye's ladder as described
// in "How to precompute a ladder" in SAC'2017.
// w = [x1,z1,x2,z2] order must be preserved.
func (c *curve) ladderJoye(k, w []byte) {
	n := c.size
	copy(w[0*n:1*n], c.pointS)  // x1 = S
	w[1*n] = 1                  // z1 = 1
	copy(w[2*n:3*n], c.pointGS) // x2 = G-S
	w[3*n] = 1                  // z2 = 1

	swap := uint(1)
	for s := 0; s < c.n-c.h; s++ {
		i := (s + c.h) / 8
		j := (s + c.h) % 8
		bit := uint((k[i] >> uint(j)) & 1)
		mu := c.table[s*n : (s+1)*n]
		c.difAdd(w, mu, swap^bit)
		swap = bit
	}
	for s := 0; s < c.h; s++ {
		c.double(w)
	}
}

// ladderMontgomery calculates a generic scalar point multiplication
// The algorithm implemented is the left-to-right Montgomery's ladder.
// w = [x1,x2,z2,x3,z3] order must be preserved.
func (c *curve) ladderMontgomery(k, xP, w []byte) {
	n := c.size
	copy(w[0*n:1*n], xP) // x1 = xP
	w[1*n] = 1           // x2 = 1
	copy(w[3*n:4*n], xP) // x3 = xP
	w[4*n] = 1           // z3 = 1

	move := uint(0)
	for s := c.n - 1; s >= 0; s-- {
		i := s / 8
		j := s % 8
		bit := uint((k[i] >> uint(j)) & 1)
		c.ladderStep(w, move^bit)
		move = bit
	}
}
