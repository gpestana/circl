// Package xcurve provides Diffie-Hellman functions as specified in RFC7748
//
// References:
//  - Curve25519 https://cr.yp.to/ecdh.html
//  - Curve448 and Goldilocks https://eprint.iacr.org/2015/625
//  - RFC7748 https://rfc-editor.org/rfc/rfc7748.txt
package xcurve

import (
	fp255 "github.com/cloudflare/circl/math/fp25519"
	"github.com/cloudflare/circl/math/fp448"
)

const (
	SizeX25519 = 32 // Byte length of a X25519 key.
	SizeX448   = 56 // Byte length of a X448 key.
)

type Key25519 [SizeX25519]byte
type Key448 [SizeX448]byte

type X25519 struct{} // Implements X25519 Diffie-Hellman
type X448 struct{}   // Implements X448 Diffie-Hellman

func (x *X25519) Generator() Key25519 { return Key25519{c255.xCoord} }
func (x *X448) Generator() Key448     { return Key448{c448.xCoord} }
func (x *X25519) KeyGen(public, secret *Key25519) {
	const n = SizeX25519
	var w [4 * n]byte
	c255.ladderJoye(public.clamp(secret), w[:])
	x.toAffine(public, w[0*n:1*n], w[1*n:2*n])
}

func (x *X25519) Shared(shared, secret, public *Key25519) {
	const n = SizeX25519
	p := *public
	p[n-1] &= (1 << (255 % 8)) - 1
	var w [5 * n]byte
	c255.ladderMontgomery(shared.clamp(secret), p[:], w[:])
	x.toAffine(shared, w[1*n:2*n], w[2*n:3*n])
}

func (x *X25519) toAffine(k *Key25519, x1, z1 []byte) {
	X, Z := &fp255.Elt{}, &fp255.Elt{}
	copy(X[:], x1)
	copy(Z[:], z1)
	fp255.Inv(Z, Z)
	fp255.Mul(X, X, Z)
	fp255.Modp(X)
	copy(k[:], X[:])
}

func (x *X448) KeyGen(public, secret *Key448) {
	const n = SizeX448
	var w [4 * n]byte
	c448.ladderJoye(public.clamp(secret), w[:])
	x.toAffine(public, w[0*n:1*n], w[1*n:2*n])
}

func (x *X448) Shared(shared, secret, public *Key448) {
	const n = SizeX448
	var w [5 * n]byte
	c448.ladderMontgomery(shared.clamp(secret), public[:], w[:])
	x.toAffine(shared, w[1*n:2*n], w[2*n:3*n])
}

func (x *X448) toAffine(k *Key448, x1, z1 []byte) {
	X, Z := &fp448.Elt{}, &fp448.Elt{}
	copy(X[:], x1)
	copy(Z[:], z1)
	fp448.Inv(Z, Z)
	fp448.Mul(X, X, Z)
	fp448.Modp(X)
	copy(k[:], X[:])
}

// clamp converts a Key into a valid scalar
func (k *Key25519) clamp(in *Key25519) []byte {
	*k = *in
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
	return k[:]
}

func (k *Key448) clamp(in *Key448) []byte {
	*k = *in
	k[0] &= 252
	k[55] |= 128
	return k[:]
}
