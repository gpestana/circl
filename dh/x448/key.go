// Package x448 provides Diffie-Hellman functions as specified in RFC7748
// using the elliptic curve known as Curve448
//
// References:
//  - Curve448 and Goldilocks https://eprint.iacr.org/2015/625
//  - RFC7748 https://rfc-editor.org/rfc/rfc7748.txt
package x448

import (
	fp "github.com/cloudflare/circl/math/fp448"
)

// Size is the length in bytes of a key.
const Size = 56

const (
	xGenerator  = 5
	bits        = 448
	logCofactor = 2
	a24         = 39082
)

// Key is using for X25519 Diffie-Hellman protocol.
type Key [Size]byte

// clamp converts a key into a valid scalar
func (k *Key) clamp(in *Key) {
	*k = *in
	k[0] &= 252
	k[Size-1] |= 128
}

// SetGenerator assigns k = 5 which is the x-coordinate of the generator.
func (k *Key) SetGenerator() { *k = Key{xGenerator} }

// KeyGen generates a public key k from Alice's secret key.
func (k *Key) KeyGen(secret *Key) {
	// The algorithm is the right-to-left Joye's ladder as described
	// in "How to precompute a ladder" in SAC'2017.
	k.clamp(secret)
	w := [4]fp.Elt{
		fp.Elt{ // x1 = -1
			0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		fp.Elt{1}, // z1 = 1
		fp.Elt{ // x2 = G-S
			0x20, 0x27, 0x9d, 0xc9, 0x7d, 0x19, 0xb1, 0xac,
			0xf8, 0xba, 0x69, 0x1c, 0xff, 0x33, 0xac, 0x23,
			0x51, 0x1b, 0xce, 0x3a, 0x64, 0x65, 0xbd, 0xf1,
			0x23, 0xf8, 0xc1, 0x84, 0x9d, 0x45, 0x54, 0x29,
			0x67, 0xb9, 0x81, 0x1c, 0x03, 0xd1, 0xcd, 0xda,
			0x7b, 0xeb, 0xff, 0x1a, 0x88, 0x03, 0xcf, 0x3a,
			0x42, 0x44, 0x32, 0x01, 0x25, 0xb7, 0xfa, 0xf0},
		fp.Elt{1}, // z2 = 1
	}
	swap := uint(1)
	for s := 0; s < bits-logCofactor; s++ {
		i := (s + logCofactor) / 8
		j := (s + logCofactor) % 8
		bit := uint((k[i] >> uint(j)) & 1)
		difAdd(&w, &tableGenerator[s], swap^bit)
		swap = bit
	}
	double(&w)
	double(&w)
	x, z := &w[0], &w[1]
	fp.Inv(z, z)
	fp.Mul((*fp.Elt)(k), x, z)
	fp.Modp((*fp.Elt)(k))
}

// Shared generates a shared secret k using Alice's secret key and Bob's
// public key.
func (k *Key) Shared(secret, public *Key) {
	k.clamp(secret)
	xP := (fp.Elt)(*public)
	w := [5]fp.Elt{
		xP,        // x1 = xP
		fp.Elt{1}, // x2 = 1
		fp.Elt{0}, // z2 = 0
		xP,        // x3 = xP
		fp.Elt{1}, // z3 = 1
	}
	move := uint(0)
	for s := bits - 1; s >= 0; s-- {
		i := s / 8
		j := s % 8
		bit := uint((k[i] >> uint(j)) & 1)
		ladderStep(&w, move^bit)
		move = bit
	}
	x, z := &w[1], &w[2]
	fp.Inv(z, z)
	fp.Mul((*fp.Elt)(k), x, z)
	fp.Modp((*fp.Elt)(k))
}
