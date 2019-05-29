// +build !amd64,go1.12

package fp25519

import (
	"math/bits"
	"unsafe"
)

type elt64 [4]uint64

// Cmov assigns y to x if n is 1.
func Cmov(x, y *Elt, n uint) {
	x64 := (*elt64)(unsafe.Pointer(x))
	y64 := (*elt64)(unsafe.Pointer(y))
	cmov64(x64, y64, n)
}

// Cswap interchages x and y if n is 1.
func Cswap(x, y *Elt, n uint) {
	x64 := (*elt64)(unsafe.Pointer(x))
	y64 := (*elt64)(unsafe.Pointer(y))
	cswap64(x64, y64, n)
}

// Add calculates z = x+y mod p
func Add(z, x, y *Elt) {
	x64 := (*elt64)(unsafe.Pointer(x))
	y64 := (*elt64)(unsafe.Pointer(y))
	z64 := (*elt64)(unsafe.Pointer(z))
	add64(z64, x64, y64)
}

// Sub calculates z = x-y mod p
func Sub(z, x, y *Elt) {
	x64 := (*elt64)(unsafe.Pointer(x))
	y64 := (*elt64)(unsafe.Pointer(y))
	z64 := (*elt64)(unsafe.Pointer(z))
	sub64(z64, x64, y64)
}

// AddSub calculates (x,y) = (x+y mod p, x-y mod p)
func AddSub(x, y *Elt) {
	x64 := (*elt64)(unsafe.Pointer(x))
	y64 := (*elt64)(unsafe.Pointer(y))
	z64 := &elt64{}
	add64(z64, x64, y64)
	sub64(y64, x64, y64)
	*x64 = *z64
}

// Mul calculates z = x*y mod p
func Mul(z, x, y *Elt) {
	x64 := (*elt64)(unsafe.Pointer(x))
	y64 := (*elt64)(unsafe.Pointer(y))
	z64 := (*elt64)(unsafe.Pointer(z))
	mul64(z64, x64, y64)
}

// Sqr calculates z = x^2 mod p
func Sqr(z, x *Elt) {
	x64 := (*elt64)(unsafe.Pointer(x))
	z64 := (*elt64)(unsafe.Pointer(z))
	sqr64(z64, x64)
}

// Modp calculates z is between [0,p-1]
func Modp(z *Elt) { modp64((*elt64)(unsafe.Pointer(z))) }

func cmov64(x, y *elt64, n uint) {
	m := -uint64(n & 0x1)
	x[0] = (x[0] &^ m) | (y[0] & m)
	x[1] = (x[1] &^ m) | (y[1] & m)
	x[2] = (x[2] &^ m) | (y[2] & m)
	x[3] = (x[3] &^ m) | (y[3] & m)
}

func cswap64(x, y *elt64, n uint) {
	m := -uint64(n & 0x1)
	t0 := m & (x[0] ^ y[0])
	t1 := m & (x[1] ^ y[1])
	t2 := m & (x[2] ^ y[2])
	t3 := m & (x[3] ^ y[3])
	x[0] ^= t0
	x[1] ^= t1
	x[2] ^= t2
	x[3] ^= t3
	y[0] ^= t0
	y[1] ^= t1
	y[2] ^= t2
	y[3] ^= t3
}

func add64(z, x, y *elt64) {
	z0, c0 := bits.Add64(x[0], y[0], 0)
	z1, c1 := bits.Add64(x[1], y[1], c0)
	z2, c2 := bits.Add64(x[2], y[2], c1)
	z3, c3 := bits.Add64(x[3], y[3], c2)

	z0, c0 = bits.Add64(z0, (-c3)&38, 0)
	z[1], c1 = bits.Add64(z1, 0, c0)
	z[2], c2 = bits.Add64(z2, 0, c1)
	z[3], c3 = bits.Add64(z3, 0, c2)
	z[0], _ = bits.Add64(z0, (-c3)&38, 0)
}

func sub64(z, x, y *elt64) {
	z0, c0 := bits.Sub64(x[0], y[0], 0)
	z1, c1 := bits.Sub64(x[1], y[1], c0)
	z2, c2 := bits.Sub64(x[2], y[2], c1)
	z3, c3 := bits.Sub64(x[3], y[3], c2)

	z0, c0 = bits.Sub64(z0, (-c3)&38, 0)
	z[1], c1 = bits.Sub64(z1, 0, c0)
	z[2], c2 = bits.Sub64(z2, 0, c1)
	z[3], c3 = bits.Sub64(z3, 0, c2)
	z[0], _ = bits.Sub64(z0, (-c3)&38, 0)
}

func modp64(x *elt64) {
	x3 := x[3]
	// CX = C[255] ? 38 : 19
	cx := uint64(19) << (x3 >> 63)
	// PUT BIT 255 IN CARRY FLAG AND CLEAR
	x3 &^= 1 << 63

	x0, c0 := bits.Add64(x[0], cx, 0)
	x1, c1 := bits.Add64(x[1], 0, c0)
	x2, c2 := bits.Add64(x[2], 0, c1)
	x3, _ = bits.Add64(x3, 0, c2)

	// TEST FOR BIT 255 AGAIN; ONLY TRIGGERED ON OVERFLOW MODULO 2^255-19
	// cx = C[255] ? 0 : 19
	cx = uint64(19) &^ (-(x3 >> 63))
	// CLEAR BIT 255
	x3 &^= 1 << 63

	x[0], c0 = bits.Sub64(x0, cx, 0)
	x[1], c1 = bits.Sub64(x1, 0, c0)
	x[2], c2 = bits.Sub64(x2, 0, c1)
	x[3], _ = bits.Sub64(x3, 0, c2)
}

func mul64(z, x, y *elt64) {
	x0, x1, x2, x3 := x[0], x[1], x[2], x[3]
	yi := y[0]
	h0, l0 := bits.Mul64(x0, yi)
	h1, l1 := bits.Mul64(x1, yi)
	h2, l2 := bits.Mul64(x2, yi)
	h3, l3 := bits.Mul64(x3, yi)

	b0 := l0
	a0, c0 := bits.Add64(h0, l1, 0)
	a1, c1 := bits.Add64(h1, l2, c0)
	a2, c2 := bits.Add64(h2, l3, c1)
	a3, _ := bits.Add64(h3, 0, c2)

	yi = y[1]
	h0, l0 = bits.Mul64(x0, yi)
	h1, l1 = bits.Mul64(x1, yi)
	h2, l2 = bits.Mul64(x2, yi)
	h3, l3 = bits.Mul64(x3, yi)

	b1, c0 := bits.Add64(a0, l0, 0)
	h0, c1 = bits.Add64(h0, l1, c0)
	h1, c2 = bits.Add64(h1, l2, c1)
	h2, c3 := bits.Add64(h2, l3, c2)
	h3, _ = bits.Add64(h3, 0, c3)

	a0, c0 = bits.Add64(a1, h0, 0)
	a1, c1 = bits.Add64(a2, h1, c0)
	a2, c2 = bits.Add64(a3, h2, c1)
	a3, _ = bits.Add64(0, h3, c2)

	yi = y[2]
	h0, l0 = bits.Mul64(x0, yi)
	h1, l1 = bits.Mul64(x1, yi)
	h2, l2 = bits.Mul64(x2, yi)
	h3, l3 = bits.Mul64(x3, yi)

	b2, c0 := bits.Add64(a0, l0, 0)
	h0, c1 = bits.Add64(h0, l1, c0)
	h1, c2 = bits.Add64(h1, l2, c1)
	h2, c3 = bits.Add64(h2, l3, c2)
	h3, _ = bits.Add64(h3, 0, c3)

	a0, c0 = bits.Add64(a1, h0, 0)
	a1, c1 = bits.Add64(a2, h1, c0)
	a2, c2 = bits.Add64(a3, h2, c1)
	a3, _ = bits.Add64(0, h3, c2)

	yi = y[3]
	h0, l0 = bits.Mul64(x0, yi)
	h1, l1 = bits.Mul64(x1, yi)
	h2, l2 = bits.Mul64(x2, yi)
	h3, l3 = bits.Mul64(x3, yi)

	b3, c0 := bits.Add64(a0, l0, 0)
	h0, c1 = bits.Add64(h0, l1, c0)
	h1, c2 = bits.Add64(h1, l2, c1)
	h2, c3 = bits.Add64(h2, l3, c2)
	h3, _ = bits.Add64(h3, 0, c3)

	b4, c0 := bits.Add64(a1, h0, 0)
	b5, c1 := bits.Add64(a2, h1, c0)
	b6, c2 := bits.Add64(a3, h2, c1)
	b7, _ := bits.Add64(0, h3, c2)

	// Reduction
	h0, l0 = bits.Mul64(b4, 38)
	h1, l1 = bits.Mul64(b5, 38)
	h2, l2 = bits.Mul64(b6, 38)
	h3, l3 = bits.Mul64(b7, 38)

	l1, c0 = bits.Add64(h0, l1, 0)
	l2, c1 = bits.Add64(h1, l2, c0)
	l3, c2 = bits.Add64(h2, l3, c1)
	l4, _ := bits.Add64(h3, 0, c2)

	l0, c0 = bits.Add64(l0, b0, 0)
	l1, c1 = bits.Add64(l1, b1, c0)
	l2, c2 = bits.Add64(l2, b2, c1)
	l3, c3 = bits.Add64(l3, b3, c2)
	l4, _ = bits.Add64(l4, 0, c3)

	_, l4 = bits.Mul64(l4, 38)
	l0, c0 = bits.Add64(l0, l4, 0)
	z[1], c1 = bits.Add64(l1, 0, c0)
	z[2], c2 = bits.Add64(l2, 0, c1)
	z[3], c3 = bits.Add64(l3, 0, c2)
	z[0], _ = bits.Add64(l0, (-c3)&38, 0)
}

func sqr64(z, x *elt64) { Sqrn(z, x, 1) }

// Sqrn calculates z = x^(2^n) mod p
func Sqrn(z, x *elt64, n uint) {
	z0 := x[0]
	z1 := x[1]
	z2 := x[2]
	z3 := x[3]
	for {
		if n == 0 {
			z[0] = z0
			z[1] = z1
			z[2] = z2
			z[3] = z3
			return
		}
		n--

		h0, a0 := bits.Mul64(z0, z1)
		h1, l1 := bits.Mul64(z0, z2)
		h2, l2 := bits.Mul64(z0, z3)
		h3, l3 := bits.Mul64(z3, z1)
		h4, l4 := bits.Mul64(z3, z2)
		h, l := bits.Mul64(z1, z2)

		a1, c0 := bits.Add64(l1, h0, 0)
		a2, c1 := bits.Add64(l2, h1, c0)
		a3, c2 := bits.Add64(l3, h2, c1)
		a4, c3 := bits.Add64(l4, h3, c2)
		a5, _ := bits.Add64(h4, 0, c3)

		a2, c0 = bits.Add64(a2, l, 0)
		a3, c1 = bits.Add64(a3, h, c0)
		a4, c2 = bits.Add64(a4, 0, c1)
		a5, c3 = bits.Add64(a5, 0, c2)
		a6, _ := bits.Add64(0, 0, c3)

		a0, c0 = bits.Add64(a0, a0, 0)
		a1, c1 = bits.Add64(a1, a1, c0)
		a2, c2 = bits.Add64(a2, a2, c1)
		a3, c3 = bits.Add64(a3, a3, c2)
		a4, c4 := bits.Add64(a4, a4, c3)
		a5, c5 := bits.Add64(a5, a5, c4)
		a6, _ = bits.Add64(a6, a6, c5)

		b1, b0 := bits.Mul64(z0, z0)
		b3, b2 := bits.Mul64(z1, z1)
		b5, b4 := bits.Mul64(z2, z2)
		b7, b6 := bits.Mul64(z3, z3)

		b1, c0 = bits.Add64(b1, a0, 0)
		b2, c1 = bits.Add64(b2, a1, c0)
		b3, c2 = bits.Add64(b3, a2, c1)
		b4, c3 = bits.Add64(b4, a3, c2)
		b5, c4 = bits.Add64(b5, a4, c3)
		b6, c5 = bits.Add64(b6, a5, c4)
		b7, _ = bits.Add64(b7, a6, c5)

		// Reduction
		h0, l0 := bits.Mul64(b4, 38)
		h1, l1 = bits.Mul64(b5, 38)
		h2, l2 = bits.Mul64(b6, 38)
		h3, l3 = bits.Mul64(b7, 38)

		l1, c0 = bits.Add64(h0, l1, 0)
		l2, c1 = bits.Add64(h1, l2, c0)
		l3, c2 = bits.Add64(h2, l3, c1)
		l4, _ = bits.Add64(h3, 0, c2)

		l0, c0 = bits.Add64(l0, b0, 0)
		l1, c1 = bits.Add64(l1, b1, c0)
		l2, c2 = bits.Add64(l2, b2, c1)
		l3, c3 = bits.Add64(l3, b3, c2)
		l4, _ = bits.Add64(l4, 0, c3)

		_, l4 = bits.Mul64(l4, 38)
		z0, c0 = bits.Add64(l0, l4, 0)
		z1, c1 = bits.Add64(l1, 0, c0)
		z2, c2 = bits.Add64(l2, 0, c1)
		z3, c3 = bits.Add64(l3, 0, c2)
		z0, _ = bits.Add64(z0, (-c3)&38, 0)
	}
}
