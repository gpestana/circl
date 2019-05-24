// Package x25519 provides Diffie-Hellman functions as specified in RFC7748
// using the elliptic curve known as Curve25519.
//
// References:
//  - Curve25519 https://cr.yp.to/ecdh.html
//  - RFC7748 https://rfc-editor.org/rfc/rfc7748.txt
package x25519

import (
	fp "github.com/cloudflare/circl/math/fp25519"
)

// Size is the length in bytes of a key.
const Size = 32

const (
	xGenerator  = 9
	bits        = 255
	logCofactor = 3
	a24         = 121666
)

// Key is using for X25519 Diffie-Hellman protocol.
type Key [Size]byte

// clamp converts a key into a valid scalar
func (k *Key) clamp(in *Key) {
	*k = *in
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
}

// SetGenerator assigns k = 9 which is the x-coordinate of the generator.
func (k *Key) SetGenerator() { *k = Key{xGenerator} }

// KeyGen generates a public key k from Alice's secret key.
func (k *Key) KeyGen(secret *Key) {
	// The algorithm is the right-to-left Joye's ladder as described
	// in "How to precompute a ladder" in SAC'2017.
	k.clamp(secret)
	w := [4]fp.Elt{
		fp.Elt{1}, // x1 = 1
		fp.Elt{1}, // z1 = 1
		fp.Elt{ // x2 = G-S
			0xbd, 0xaa, 0x2f, 0xc8, 0xfe, 0xe1, 0x94, 0x7e,
			0xf8, 0xed, 0xb2, 0x14, 0xae, 0x95, 0xf0, 0xbb,
			0xe2, 0x48, 0x5d, 0x23, 0xb9, 0xa0, 0xc7, 0xad,
			0x34, 0xab, 0x7c, 0xe2, 0xee, 0xcd, 0xae, 0x1e},
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
	// [RFC-7748] When receiving such an array, implementations
	// of X25519 (but not X448) MUST mask the most significant
	// bit in the final byte.
	xP := (fp.Elt)(*public)
	xP[31] &= (1 << (bits % 8)) - 1
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
