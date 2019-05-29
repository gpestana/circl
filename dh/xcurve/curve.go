package xcurve

import (
	fp255 "github.com/cloudflare/circl/math/fp25519"
	"github.com/cloudflare/circl/math/fp448"
)

type curve255 struct{}
type curve448 struct{}

var c255 curve255
var c448 curve448

// ladderJoye calculates a fixed-point multiplication with the generator point.
// The algorithm is the right-to-left Joye's ladder as described
// in "How to precompute a ladder" in SAC'2017.
func (c *curve255) ladderJoye(k *Key25519) {
	w := [5]fp255.Elt{} // [mu,x1,z1,x2,z2] order must be preserved.
	fp255.SetOne(&w[1]) // x1 = 1
	fp255.SetOne(&w[2]) // z1 = 1
	w[3] = fp255.Elt{   // x2 = G-S
		0xbd, 0xaa, 0x2f, 0xc8, 0xfe, 0xe1, 0x94, 0x7e,
		0xf8, 0xed, 0xb2, 0x14, 0xae, 0x95, 0xf0, 0xbb,
		0xe2, 0x48, 0x5d, 0x23, 0xb9, 0xa0, 0xc7, 0xad,
		0x34, 0xab, 0x7c, 0xe2, 0xee, 0xcd, 0xae, 0x1e,
	}
	fp255.SetOne(&w[4]) // z2 = 1

	const n = 255
	const h = 3
	swap := uint(1)
	for s := 0; s < n-h; s++ {
		i := (s + h) / 8
		j := (s + h) % 8
		bit := uint((k[i] >> uint(j)) & 1)
		copy(w[0][:], tableGenerator255[s*SizeX25519:(s+1)*SizeX25519])
		c.difAdd(&w, swap^bit)
		swap = bit
	}
	for s := 0; s < h; s++ {
		c.double(&w[1], &w[2])
	}
	c.toAffine((*[fp255.Size]byte)(k), &w[1], &w[2])
}

// ladderMontgomery calculates a generic scalar point multiplication
// The algorithm implemented is the left-to-right Montgomery's ladder.
func (c *curve255) ladderMontgomery(k, xP *Key25519) {
	w := [5]fp255.Elt{}      // [x1, x2, z2, x3, z3] order must be preserved.
	w[0] = *(*fp255.Elt)(xP) // x1 = xP
	fp255.SetOne(&w[1])      // x2 = 1
	w[3] = *(*fp255.Elt)(xP) // x3 = xP
	fp255.SetOne(&w[4])      // z3 = 1

	move := uint(0)
	for s := 255 - 1; s >= 0; s-- {
		i := s / 8
		j := s % 8
		bit := uint((k[i] >> uint(j)) & 1)
		c.ladderStep(&w, move^bit)
		move = bit
	}
	c.toAffine((*[fp255.Size]byte)(k), &w[1], &w[2])
}

// ladderJoye calculates a fixed-point multiplication with the generator point.
// The algorithm is the right-to-left Joye's ladder as described
// in "How to precompute a ladder" in SAC'2017.
func (c *curve448) ladderJoye(k *Key448) {
	w := [5]fp448.Elt{} // [mu,x1,z1,x2,z2] order must be preserved.
	w[1] = fp448.Elt{   // x1 = S
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	fp448.SetOne(&w[2]) // z1 = 1
	w[3] = fp448.Elt{   // x2 = G-S
		0x20, 0x27, 0x9d, 0xc9, 0x7d, 0x19, 0xb1, 0xac,
		0xf8, 0xba, 0x69, 0x1c, 0xff, 0x33, 0xac, 0x23,
		0x51, 0x1b, 0xce, 0x3a, 0x64, 0x65, 0xbd, 0xf1,
		0x23, 0xf8, 0xc1, 0x84, 0x9d, 0x45, 0x54, 0x29,
		0x67, 0xb9, 0x81, 0x1c, 0x03, 0xd1, 0xcd, 0xda,
		0x7b, 0xeb, 0xff, 0x1a, 0x88, 0x03, 0xcf, 0x3a,
		0x42, 0x44, 0x32, 0x01, 0x25, 0xb7, 0xfa, 0xf0,
	}
	fp448.SetOne(&w[4]) // z2 = 1

	const n = 448
	const h = 2
	swap := uint(1)
	for s := 0; s < n-h; s++ {
		i := (s + h) / 8
		j := (s + h) % 8
		bit := uint((k[i] >> uint(j)) & 1)
		copy(w[0][:], tableGenerator448[s*SizeX448:(s+1)*SizeX448])
		c.difAdd(&w, swap^bit)
		swap = bit
	}
	for s := 0; s < h; s++ {
		c.double(&w[1], &w[2])
	}
	c.toAffine((*[fp448.Size]byte)(k), &w[1], &w[2])
}

// ladderMontgomery calculates a generic scalar point multiplication
// The algorithm implemented is the left-to-right Montgomery's ladder.
func (c *curve448) ladderMontgomery(k, xP *Key448) {
	w := [5]fp448.Elt{}      // [x1, x2, z2, x3, z3] order must be preserved.
	w[0] = *(*fp448.Elt)(xP) // x1 = xP
	fp448.SetOne(&w[1])      // x2 = 1
	w[3] = *(*fp448.Elt)(xP) // x3 = xP
	fp448.SetOne(&w[4])      // z3 = 1

	move := uint(0)
	for s := 448 - 1; s >= 0; s-- {
		i := s / 8
		j := s % 8
		bit := uint((k[i] >> uint(j)) & 1)
		c.ladderStep(&w, move^bit)
		move = bit
	}
	c.toAffine((*[fp448.Size]byte)(k), &w[1], &w[2])
}

func (c *curve255) toAffine(k *[fp255.Size]byte, x, z *fp255.Elt) {
	fp255.Inv(z, z)
	fp255.Mul(x, x, z)
	fp255.ToBytes(k[:], x)
}

func (c *curve448) toAffine(k *[fp448.Size]byte, x, z *fp448.Elt) {
	fp448.Inv(z, z)
	fp448.Mul(x, x, z)
	fp448.ToBytes(k[:], x)
}
