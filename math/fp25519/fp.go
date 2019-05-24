// Package fp25519 provides prime field arithmetic over GF(2^255-19).
package fp25519

// Size in bytes of an element.
const Size = 32

// Elt is a prime field element.
type Elt [Size]byte

// p is the prime modulus 2^255-19
var p = Elt{
	0xed, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
}

// P returns the prime modulus 2^255-19.
func P() Elt { return p }

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
