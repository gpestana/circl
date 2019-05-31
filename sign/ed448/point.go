package ed448

import (
	"encoding/binary"
	"fmt"

	"github.com/cloudflare/circl/math/fp448"
)

type pointR1 struct{ x, y, z, ta, tb fp448.Elt }
type pointR2 struct {
	pointR3
	z2 fp448.Elt
}
type pointR3 struct{ addYX, subYX, dt2 fp448.Elt }

func (P pointR1) String() string {
	return fmt.Sprintf("\nx=  %v\ny=  %v\nta= %v\ntb= %v\nz=  %v",
		P.x, P.y, P.ta, P.tb, P.z)
}
func (P pointR3) String() string {
	return fmt.Sprintf("\naddYX= %v\nsubYX= %v\ndt2=  %v",
		P.addYX, P.subYX, P.dt2)
}
func (P pointR2) String() string {
	return fmt.Sprintf("%v\nz2=  %v", &P.pointR3, P.z2)
}

func (P *pointR1) neg() {
	fp448.Neg(&P.x, &P.x)
	fp448.Neg(&P.ta, &P.ta)
}

func (P *pointR1) SetIdentity() {
	fp448.SetZero(&P.x)
	fp448.SetOne(&P.y)
	fp448.SetOne(&P.z)
	fp448.SetZero(&P.ta)
	fp448.SetZero(&P.tb)
}

func (P *pointR1) toAffine() {
	fp448.Inv(&P.z, &P.z)
	fp448.Mul(&P.x, &P.x, &P.z)
	fp448.Mul(&P.y, &P.y, &P.z)
	fp448.Modp(&P.x)
	fp448.Modp(&P.y)
	fp448.SetOne(&P.z)
	P.ta = P.x
	P.tb = P.y
}

func (P *pointR1) ToBytes(k []byte) {
	P.toAffine()
	var x [fp448.Size]byte
	fp448.ToBytes(k[:fp448.Size], &P.y)
	fp448.ToBytes(x[:], &P.x)
	b := x[0] & 1
	k[Size-1] = k[Size-1] | (b << 7)
}

func isGreaterThanP(x *fp448.Elt) bool {
	p := fp448.P()
	n := 8
	x0 := binary.LittleEndian.Uint64(x[0*n : 1*n])
	x1 := binary.LittleEndian.Uint64(x[1*n : 2*n])
	x2 := binary.LittleEndian.Uint64(x[2*n : 3*n])
	x3 := binary.LittleEndian.Uint64(x[3*n : 4*n])
	p0 := binary.LittleEndian.Uint64(p[0*n : 1*n])
	p1 := binary.LittleEndian.Uint64(p[1*n : 2*n])
	p2 := binary.LittleEndian.Uint64(p[2*n : 3*n])
	p3 := binary.LittleEndian.Uint64(p[3*n : 4*n])

	if x3 >= p3 {
		return true
	} else if x2 >= p2 {
		return true
	} else if x1 >= p1 {
		return true
	} else if x0 >= p0 {
		return true
	}
	return false
}

func (P *pointR1) FromBytes(k []byte) bool {
	if len(k) != Size {
		panic("wrong size")
	}
	signX := k[Size-1] >> 7
	copy(P.y[:], k)
	if isGreaterThanP(&P.y) {
		return false
	}
	d := &fp448.Elt{}
	copy(d[:], curve.paramD[:fp448.Size])
	one, u, v := &fp448.Elt{}, &fp448.Elt{}, &fp448.Elt{}
	fp448.SetOne(one)
	fp448.Sqr(u, &P.y)              // u = y^2
	fp448.Mul(v, u, d)              // v = dy^2
	fp448.Sub(u, u, one)            // u = y^2-1
	fp448.Add(v, v, one)            // v = dy^2+1
	ok := fp448.InvSqrt(&P.x, u, v) // x = sqrt(u/v)
	if !ok {
		return false
	}
	fp448.Modp(&P.x) // x = x mod p
	if fp448.IsZero(&P.x) && signX == 1 {
		return false
	}
	if signX != (P.x[0] & 1) {
		fp448.Neg(&P.x, &P.x)
	}
	P.ta = P.x
	P.tb = P.y
	fp448.SetOne(&P.z)
	return true
}

func (P *pointR1) double() {
	Px, Py, Pz, Pta, Ptb := &P.x, &P.y, &P.z, &P.ta, &P.tb
	a := Px
	b := Py
	c := Pz
	d := Pta
	e := Ptb
	f := b
	g := a
	fp448.Add(e, Px, Py)
	fp448.Sqr(a, Px)
	fp448.Sqr(b, Py)
	fp448.Sqr(c, Pz)
	fp448.Add(c, c, c)
	fp448.Add(d, a, b)
	fp448.Sqr(e, e)
	fp448.Sub(e, e, d)
	fp448.Sub(f, b, a)
	fp448.Sub(g, c, f)
	fp448.Mul(Pz, f, g)
	fp448.Mul(Px, e, g)
	fp448.Mul(Py, d, f)
}

func (P *pointR1) mixAdd(Q *pointR3) {
	addYX := &Q.addYX
	subYX := &Q.subYX
	dt2 := &Q.dt2
	Px := &P.x
	Py := &P.y
	Pz := &P.z
	Pta := &P.ta
	Ptb := &P.tb
	a := Px
	b := Py
	c := &fp448.Elt{}
	d := b
	e := Pta
	f := a
	g := b
	h := Ptb
	fp448.Mul(c, Pta, Ptb)
	fp448.Sub(h, b, a)
	fp448.Add(b, b, a)
	fp448.Mul(a, h, subYX)
	fp448.Mul(b, b, addYX)
	fp448.Sub(e, b, a)
	fp448.Add(h, b, a)
	fp448.Add(d, Pz, Pz)
	fp448.Mul(c, c, dt2)
	fp448.Sub(f, d, c)
	fp448.Add(g, d, c)
	fp448.Mul(Pz, f, g)
	fp448.Mul(Px, e, f)
	fp448.Mul(Py, g, h)
}

func (P *pointR1) add(Q *pointR2) {
	addYX := &Q.addYX
	subYX := &Q.subYX
	dt2 := &Q.dt2
	z2 := &Q.z2
	Px := &P.x
	Py := &P.y
	Pz := &P.z
	Pta := &P.ta
	Ptb := &P.tb
	a := Px
	b := Py
	c := &fp448.Elt{}
	d := b
	e := Pta
	f := a
	g := b
	h := Ptb
	fp448.Mul(c, Pta, Ptb)
	fp448.Sub(h, b, a)
	fp448.Add(b, b, a)
	fp448.Mul(a, h, subYX)
	fp448.Mul(b, b, addYX)
	fp448.Sub(e, b, a)
	fp448.Add(h, b, a)
	fp448.Mul(d, Pz, z2)
	fp448.Mul(c, c, dt2)
	fp448.Sub(f, d, c)
	fp448.Add(g, d, c)
	fp448.Mul(Pz, f, g)
	fp448.Mul(Px, e, f)
	fp448.Mul(Py, g, h)
}

func (P *pointR1) oddMultiples(T []pointR2) {
	var R pointR2
	n := len(T)
	T[0].fromR1(P)
	_2P := *P
	_2P.double()
	R.fromR1(&_2P)
	for i := 1; i < n; i++ {
		P.add(&R)
		T[i].fromR1(P)
	}
}

func (P *pointR1) isEqual(Q *pointR1) bool {
	l, r := &fp448.Elt{}, &fp448.Elt{}
	fp448.Mul(l, &P.x, &Q.z)
	fp448.Mul(r, &Q.x, &P.z)
	fp448.Sub(l, l, r)
	b := fp448.IsZero(l)
	fp448.Mul(l, &P.y, &Q.z)
	fp448.Mul(r, &Q.y, &P.z)
	fp448.Sub(l, l, r)
	b = b && fp448.IsZero(l)
	fp448.Mul(l, &P.ta, &P.tb)
	fp448.Mul(l, l, &Q.z)
	fp448.Mul(r, &Q.ta, &Q.tb)
	fp448.Mul(r, r, &P.z)
	fp448.Sub(l, l, r)
	b = b && fp448.IsZero(l)
	return b
}

func (P *pointR3) neg() {
	P.addYX, P.subYX = P.subYX, P.addYX
	fp448.Neg(&P.dt2, &P.dt2)
}

func (P *pointR2) fromR1(Q *pointR1) {
	d := &fp448.Elt{}
	copy(d[:], curve.paramD[:fp448.Size])
	fp448.Add(&P.addYX, &Q.y, &Q.x)
	fp448.Sub(&P.subYX, &Q.y, &Q.x)
	fp448.Mul(&P.dt2, &Q.ta, &Q.tb)
	fp448.Mul(&P.dt2, &P.dt2, d)
	fp448.Add(&P.dt2, &P.dt2, &P.dt2)
	fp448.Add(&P.z2, &Q.z, &Q.z)
}

func (P *pointR3) cneg(b int) {
	t := &fp448.Elt{}
	fp448.Cswap(&P.addYX, &P.subYX, uint(b))
	fp448.Neg(t, &P.dt2)
	fp448.Cmov(&P.dt2, t, uint(b))
}

func (P *pointR3) cmov(Q *pointR3, b int) {
	fp448.Cmov(&P.addYX, &Q.addYX, uint(b))
	fp448.Cmov(&P.subYX, &Q.subYX, uint(b))
	fp448.Cmov(&P.dt2, &Q.dt2, uint(b))
}
