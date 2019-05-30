// Package fp25519 provides prime field arithmetic over GF(2^255-19).
package fp25519

import "github.com/cloudflare/circl/internal/conv"

// Size in bytes of an element.
const Size = 32

// Elt is a prime field element.
type Elt [Size]byte

func (e Elt) String() string { return conv.BytesLe2Hex(e[:]) }

// p is the prime modulus 2^255-19
var p = Elt{
	0xed, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
}

// P returns the prime modulus 2^255-19.
func P() Elt { return p }

// ToBytes returns the little-endian byte representation of x.
func ToBytes(b []byte, x *Elt) {
	if len(b) != Size {
		panic("wrong size")
	}
	Modp(x)
	copy(b, x[:])
}
func IsZero(x *Elt) bool { Modp(x); return *x == Elt{} }
func SetZero(x *Elt)     { *x = Elt{} }
func SetOne(x *Elt)      { SetZero(x); x[0] = 1 }
func Neg(z, x *Elt)      { Sub(z, &p, x) }

// InvSqrt calculates z = sqrt(x/y)
func InvSqrt(z, x, y *Elt) {
	sqrtMinusOne := &Elt{
		0xc4, 0xee, 0x1b, 0x27, 0x4a, 0x0e, 0xa0, 0xb0,
		0x2f, 0x43, 0x18, 0x06, 0xad, 0x2f, 0xe4, 0x78,
		0x2b, 0x4d, 0x00, 0x99, 0x3d, 0xfb, 0xd7, 0xa7,
		0x2b, 0x83, 0x24, 0x80, 0x4f, 0xc1, 0xdf, 0x0b,
	}
	t0, t1, t2, t3 := &Elt{}, &Elt{}, &Elt{}, &Elt{}

	Mul(t0, x, y)   // t0 = u*v
	Sqr(t1, y)      // t1 = v^2
	Mul(t2, t0, t1) // t2 = u*v^3
	Sqr(t0, t1)     // t0 = v^4
	Mul(t1, t0, t2) // t1 = u*v^7

	var Tab [4]*Elt
	Tab[0] = &Elt{}
	Tab[1] = &Elt{}
	Tab[2] = t3
	Tab[3] = t1

	*Tab[0] = *t1
	Sqr(Tab[0], Tab[0])
	Sqr(Tab[1], Tab[0])
	Sqr(Tab[1], Tab[1])
	Mul(Tab[1], Tab[1], Tab[3])
	Mul(Tab[0], Tab[0], Tab[1])
	Sqr(Tab[0], Tab[0])
	Mul(Tab[0], Tab[0], Tab[1])
	Sqr(Tab[1], Tab[0])
	for i := 0; i < 4; i++ {
		Sqr(Tab[1], Tab[1])
	}
	Mul(Tab[1], Tab[1], Tab[0])
	Sqr(Tab[2], Tab[1])
	for i := 0; i < 4; i++ {
		Sqr(Tab[2], Tab[2])
	}
	Mul(Tab[2], Tab[2], Tab[0])
	Sqr(Tab[1], Tab[2])
	for i := 0; i < 14; i++ {
		Sqr(Tab[1], Tab[1])
	}
	Mul(Tab[1], Tab[1], Tab[2])
	Sqr(Tab[2], Tab[1])
	for i := 0; i < 29; i++ {
		Sqr(Tab[2], Tab[2])
	}
	Mul(Tab[2], Tab[2], Tab[1])
	Sqr(Tab[1], Tab[2])
	for i := 0; i < 59; i++ {
		Sqr(Tab[1], Tab[1])
	}
	Mul(Tab[1], Tab[1], Tab[2])
	for i := 0; i < 5; i++ {
		Sqr(Tab[1], Tab[1])
	}
	Mul(Tab[1], Tab[1], Tab[0])
	Sqr(Tab[2], Tab[1])
	for i := 0; i < 124; i++ {
		Sqr(Tab[2], Tab[2])
	}
	Mul(Tab[2], Tab[2], Tab[1])
	Sqr(Tab[2], Tab[2])
	Sqr(Tab[2], Tab[2])
	Mul(Tab[2], Tab[2], Tab[3])

	Mul(t0, z, t2) // uv^(p+3)/8 = uv3*uv^(p-5)/8
	// Checking whether v*uv_p38^2 == -u
	Sqr(t0, t0)    // uv3 = uv_p38^2
	Mul(t1, t2, y) // uv = uv3*v
	Neg(x, x)      //  u = -u
	Sub(t0, t0, x) // t0 = u-uv
	if IsZero(t0) {
		Mul(z, z, sqrtMinusOne) //  uv_p38 = uv_p38*sqrt(-1)
	}
}

// Inv calculates z = 1/x mod p
func Inv(z, x *Elt) {
	x0, x1, x2 := &Elt{}, &Elt{}, &Elt{}
	Sqr(x1, x)
	Sqr(x0, x1)
	Sqr(x0, x0)
	Mul(x0, x0, x)
	Mul(z, x0, x1)
	Sqr(x1, z)
	Mul(x0, x0, x1)
	Sqr(x1, x0)
	for i := 0; i < 4; i++ {
		Sqr(x1, x1)
	}
	Mul(x0, x0, x1)
	Sqr(x1, x0)
	for i := 0; i < 9; i++ {
		Sqr(x1, x1)
	}
	Mul(x1, x1, x0)
	Sqr(x2, x1)
	for i := 0; i < 19; i++ {
		Sqr(x2, x2)
	}
	Mul(x2, x2, x1)
	for i := 0; i < 10; i++ {
		Sqr(x2, x2)
	}
	Mul(x2, x2, x0)
	Sqr(x0, x2)
	for i := 0; i < 49; i++ {
		Sqr(x0, x0)
	}
	Mul(x0, x0, x2)
	Sqr(x1, x0)
	for i := 0; i < 99; i++ {
		Sqr(x1, x1)
	}
	Mul(x1, x1, x0)
	for i := 0; i < 50; i++ {
		Sqr(x1, x1)
	}
	Mul(x1, x1, x2)
	for i := 0; i < 5; i++ {
		Sqr(x1, x1)
	}
	Mul(z, z, x1)
}
