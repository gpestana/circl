// +build amd64

// Package xcurve provides Diffie-Hellman functions as specified in RFC-7748
//
// References:
//  - Curve25519 https://cr.yp.to/ecdh.html
//  - Curve448 and Goldilocks https://eprint.iacr.org/2015/625
//  - RFC7748 https://rfc-editor.org/rfc/rfc7748.txt
package xcurve

const (
	// SizeX25519 is the length in bytes of a X25519 key.
	SizeX25519 = 32
	// SizeX448 is the length in bytes of a X448 key.
	SizeX448 = 56
)

// Key25519 represents a X25519 key.
type Key25519 [SizeX25519]byte

// Key448 represents a X448 key.
type Key448 [SizeX448]byte

// X25519 instantiates a receiver able to perform X25519 Diffie-Hellman operations.
type X25519 struct{}

// X448 instantiates a receiver able to perform X448 Diffie-Hellman operations.
type X448 struct{}

// KeyGen obtains a public key given a secret key.
func (x *X25519) KeyGen(public, secret *Key25519) {
	c255.ladderJoye(public.clamp(secret))
}

// Shared calculates Alice's shared key from Alice's secret key and Bob's public key.
func (x *X25519) Shared(shared, secret, public *Key25519) {
	p := *public
	p[31] &= (1 << (255 % 8)) - 1
	c255.ladderMontgomery(shared.clamp(secret), p[:])
}

// KeyGen obtains a public key given a secret key.
func (x *X448) KeyGen(public, secret *Key448) {
	c448.ladderJoye(public.clamp(secret))
}

// Shared calculates Alice's shared key from Alice's secret key and Bob's public key.
func (x *X448) Shared(shared, secret, public *Key448) {
	c448.ladderMontgomery(shared.clamp(secret), public[:])
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
