package eddsa

import (
	"fmt"

	fp255 "github.com/cloudflare/circl/math/fp25519"
)

type point255R1 struct{ x, y, z, ta, tb fp255.Elt }
type point255R2 struct {
	point255R3
	z2 fp255.Elt
}
type point255R3 struct{ addYX, subYX, dt2 fp255.Elt }

func (P *point255R1) String() string {
	return fmt.Sprintf("\nx=  %v\ny=  %v\nta= %v\ntb= %v\nz=  %v",
		P.x, P.y, P.ta, P.tb, P.z)
}
func (P *point255R3) String() string {
	return fmt.Sprintf("\naddYX= %v\nsubYX= %v\ndt2=  %v",
		P.addYX, P.subYX, P.dt2)
}
func (P *point255R2) String() string {
	return fmt.Sprintf("%v\nz2=  %v", &P.point255R3, P.z2)
}

func (P *point255R1) neg() {
	fp255.Neg(&P.x, &P.x)
	fp255.Neg(&P.ta, &P.ta)
}

func (P *point255R1) copy() pointR1 { Q := *P; return &Q }

func (P *point255R1) SetIdentity() {
	fp255.SetZero(&P.x)
	fp255.SetOne(&P.y)
	fp255.SetOne(&P.z)
	fp255.SetZero(&P.ta)
	fp255.SetZero(&P.tb)
}

func (P *point255R1) SetGenerator() {
	copy(P.x[:], edwards25519.genX)
	copy(P.y[:], edwards25519.genY)
	fp255.SetOne(&P.z)
	P.ta = P.x
	P.tb = P.y
}

func (P *point255R1) toAffine() {
	fp255.Inv(&P.z, &P.z)
	fp255.Mul(&P.x, &P.x, &P.z)
	fp255.Mul(&P.y, &P.y, &P.z)
	fp255.Modp(&P.x)
	fp255.Modp(&P.y)
	fp255.SetOne(&P.z)
	P.ta = P.x
	P.tb = P.y
}

func (P *point255R1) ToBytes(k []byte) {
	P.toAffine()
	var x [32]byte
	fp255.ToBytes(k, &P.y)
	fp255.ToBytes(x[:], &P.x)
	b := x[0] & 1
	k[31] = k[31] | (b << 7)
}

func (P *point255R1) FromBytes(k []byte) bool {
	if len(k) != 32 {
		panic("wrong size")
	}
	signX := k[31] >> 7
	copy(P.y[:], k)
	P.y[31] &= 0x7F

	d := &fp255.Elt{}
	copy(d[:], edwards25519.paramD)

	one, u, v := &fp255.Elt{}, &fp255.Elt{}, &fp255.Elt{}
	fp255.SetOne(one)
	fp255.Sqr(u, &P.y)              // u = y^2
	fp255.Mul(v, u, d)              // v = dy^2
	fp255.Sub(u, u, one)            // u = y^2-1
	fp255.Add(v, v, one)            // v = dy^2+1
	ok := fp255.InvSqrt(&P.x, u, v) // x = sqrt(u/v)
	if !ok {
		return false
	}
	fp255.Modp(&P.x) // x = x mod p
	if fp255.IsZero(&P.x) && signX == 1 {
		return false
	}
	if signX != (P.x[0] & 1) {
		fp255.Neg(&P.x, &P.x)
	}
	P.ta = P.x
	P.tb = P.y
	fp255.SetOne(&P.z)
	return true
}

func (P *point255R1) double() {
	Px, Py, Pz, Pta, Ptb := &P.x, &P.y, &P.z, &P.ta, &P.tb
	a := Px
	b := Py
	c := Pz
	d := Pta
	e := Ptb
	f := b
	g := a
	fp255.Add(e, Px, Py)
	fp255.Sqr(a, Px)
	fp255.Sqr(b, Py)
	fp255.Sqr(c, Pz)
	fp255.Add(c, c, c)
	fp255.Add(d, a, b)
	fp255.Sqr(e, e)
	fp255.Sub(e, e, d)
	fp255.Sub(f, b, a)
	fp255.Sub(g, c, f)
	fp255.Mul(Pz, f, g)
	fp255.Mul(Px, e, g)
	fp255.Mul(Py, d, f)
}

func (P *point255R1) mixAdd(Q pointR3) {
	QQ, ok := Q.(*point255R3)
	if !ok {
		panic("wrong type")
	}
	addYX := &QQ.addYX
	subYX := &QQ.subYX
	dt2 := &QQ.dt2
	Px := &P.x
	Py := &P.y
	Pz := &P.z
	Pta := &P.ta
	Ptb := &P.tb
	a := Px
	b := Py
	c := &fp255.Elt{}
	d := b
	e := Pta
	f := a
	g := b
	h := Ptb
	fp255.Mul(c, Pta, Ptb)
	fp255.Sub(h, b, a)
	fp255.Add(b, b, a)
	fp255.Mul(a, h, subYX)
	fp255.Mul(b, b, addYX)
	fp255.Sub(e, b, a)
	fp255.Add(h, b, a)
	fp255.Add(d, Pz, Pz)
	fp255.Mul(c, c, dt2)
	fp255.Sub(f, d, c)
	fp255.Add(g, d, c)
	fp255.Mul(Pz, f, g)
	fp255.Mul(Px, e, f)
	fp255.Mul(Py, g, h)
}

func (P *point255R1) add(Q pointR2) {
	QQ, ok := Q.(*point255R2)
	if !ok {
		panic("wrong type")
	}
	addYX := &QQ.addYX
	subYX := &QQ.subYX
	dt2 := &QQ.dt2
	z2 := &QQ.z2
	Px := &P.x
	Py := &P.y
	Pz := &P.z
	Pta := &P.ta
	Ptb := &P.tb
	a := Px
	b := Py
	c := &fp255.Elt{}
	d := b
	e := Pta
	f := a
	g := b
	h := Ptb
	fp255.Mul(c, Pta, Ptb)
	fp255.Sub(h, b, a)
	fp255.Add(b, b, a)
	fp255.Mul(a, h, subYX)
	fp255.Mul(b, b, addYX)
	fp255.Sub(e, b, a)
	fp255.Add(h, b, a)
	fp255.Mul(d, Pz, z2)
	fp255.Mul(c, c, dt2)
	fp255.Sub(f, d, c)
	fp255.Add(g, d, c)
	fp255.Mul(Pz, f, g)
	fp255.Mul(Px, e, f)
	fp255.Mul(Py, g, h)
}

func (P *point255R1) oddMultiples(T []pointR2) {
	var R point255R2
	n := len(T)
	T[0] = new(point255R2)
	T[0].fromR1(P)
	_2P := *P
	_2P.double()
	R.fromR1(&_2P)
	Q := *P
	for i := 1; i < n; i++ {
		Q.add(&R)
		T[i] = new(point255R2)
		T[i].fromR1(&Q)
	}
}

func (P *point255R1) isEqual(Q pointR1) bool {
	QQ, ok := Q.(*point255R1)
	if !ok {
		panic("wrong type")
	}
	l, r := &fp255.Elt{}, &fp255.Elt{}
	fp255.Mul(l, &P.x, &QQ.z)
	fp255.Mul(r, &QQ.x, &P.z)
	fp255.Sub(l, l, r)
	b := fp255.IsZero(l)
	fp255.Mul(l, &P.y, &QQ.z)
	fp255.Mul(r, &QQ.y, &P.z)
	fp255.Sub(l, l, r)
	b = b && fp255.IsZero(l)
	fp255.Mul(l, &P.ta, &P.tb)
	fp255.Mul(l, l, &QQ.z)
	fp255.Mul(r, &QQ.ta, &QQ.tb)
	fp255.Mul(r, r, &P.z)
	fp255.Sub(l, l, r)
	b = b && fp255.IsZero(l)
	return b
}

func (P *point255R2) neg() pointR2 {
	Q := &point255R2{}
	Q.addYX = P.subYX
	Q.subYX = P.addYX
	fp255.Neg(&Q.dt2, &P.dt2)
	Q.z2 = P.z2
	return Q
}

func (P *point255R2) fromR1(Q pointR1) {
	QQ, ok := Q.(*point255R1)
	if !ok {
		panic("wrong type")
	}
	P.point255R3.fromR1(QQ)
	fp255.Add(&P.z2, &QQ.z, &QQ.z)
}

func (P *point255R3) neg() pointR3 {
	Q := &point255R3{}
	Q.addYX = P.subYX
	Q.subYX = P.addYX
	fp255.Neg(&Q.dt2, &P.dt2)
	return Q
}

func (P *point255R3) cneg(b int) {
	t := &fp255.Elt{}
	fp255.Cswap(&P.addYX, &P.subYX, uint(b))
	fp255.Neg(t, &P.dt2)
	fp255.Cmov(&P.dt2, t, uint(b))
}

func (P *point255R3) cmov(Q pointR3, b int) {
	QQ, ok := Q.(*point255R3)
	if !ok {
		panic("wrong type")
	}
	fp255.Cmov(&P.addYX, &QQ.addYX, uint(b))
	fp255.Cmov(&P.subYX, &QQ.subYX, uint(b))
	fp255.Cmov(&P.dt2, &QQ.dt2, uint(b))
}

func (P *point255R3) fromR1(Q pointR1) {
	QQ, ok := Q.(*point255R1)
	if !ok {
		panic("wrong type")
	}
	var d fp255.Elt
	copy(d[:], edwards25519.paramD)
	QQ.toAffine()
	fp255.Add(&P.addYX, &QQ.y, &QQ.x)
	fp255.Sub(&P.subYX, &QQ.y, &QQ.x)
	fp255.Mul(&P.dt2, &QQ.ta, &QQ.tb)
	fp255.Mul(&P.dt2, &P.dt2, &d)
	fp255.Add(&P.dt2, &P.dt2, &P.dt2)
}
