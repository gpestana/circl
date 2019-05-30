package eddsa

import (
	"fmt"

	"github.com/cloudflare/circl/math/fp448"
)

type point448R1 struct{ x, y, z, ta, tb fp448.Elt }
type point448R2 struct {
	point448R3
	z2 fp448.Elt
}
type point448R3 struct{ addYX, subYX, dt2 fp448.Elt }

func (P *point448R1) String() string {
	return fmt.Sprintf("\nx=  %v\ny=  %v\nta= %v\ntb= %v\nz=  %v",
		P.x, P.y, P.ta, P.tb, P.z)
}
func (P *point448R3) String() string {
	return fmt.Sprintf("\naddYX= %v\nsubYX= %v\ndt2=  %v",
		P.addYX, P.subYX, P.dt2)
}
func (P *point448R2) String() string {
	return fmt.Sprintf("%v\nz2=  %v", &P.point448R3, P.z2)
}

func (P *point448R1) neg() {
	fp448.Neg(&P.x, &P.x)
	fp448.Neg(&P.ta, &P.ta)
}

func (P *point448R1) copy() pointR1 { Q := *P; return &Q }

func (P *point448R1) SetIdentity() {
	fp448.SetZero(&P.x)
	fp448.SetOne(&P.y)
	fp448.SetOne(&P.z)
	fp448.SetZero(&P.ta)
	fp448.SetZero(&P.tb)
}

func (P *point448R1) SetGenerator() {
	copy(P.x[:], edwards448.genX)
	copy(P.y[:], edwards448.genY)
	fp448.SetOne(&P.z)
	P.ta = P.x
	P.tb = P.y
}

func (P *point448R1) toAffine() {
	fp448.Inv(&P.z, &P.z)
	fp448.Mul(&P.x, &P.x, &P.z)
	fp448.Mul(&P.y, &P.y, &P.z)
	fp448.Modp(&P.x)
	fp448.Modp(&P.y)
	fp448.SetOne(&P.z)
	P.ta = P.x
	P.tb = P.y
}

func (P *point448R1) ToBytes(k []byte) {
	P.toAffine()
	var x [32]byte
	fp448.ToBytes(k, &P.y)
	fp448.ToBytes(x[:], &P.x)
	b := x[0] & 1
	k[31] = k[31] | (b << 7)
}

func (P *point448R1) FromBytes(k []byte) bool {
	if len(k) != 32 {
		panic("wrong size")
	}
	signX := k[31] >> 7
	copy(P.y[:], k)
	P.y[31] &= 0x7F

	d := &fp448.Elt{}
	copy(d[:], edwards448.paramD)

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

func (P *point448R1) double() {
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

func (P *point448R1) mixAdd(Q pointR3) {
	QQ, ok := Q.(*point448R3)
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

func (P *point448R1) add(Q pointR2) {
	QQ, ok := Q.(*point448R2)
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

func (P *point448R1) oddMultiples(T []pointR2) {
	var R point448R2
	n := len(T)
	T[0] = new(point448R2)
	T[0].fromR1(P)
	_2P := *P
	_2P.double()
	R.fromR1(&_2P)
	Q := *P
	for i := 1; i < n; i++ {
		Q.add(&R)
		T[i] = new(point448R2)
		T[i].fromR1(&Q)
	}
}

func (P *point448R1) isEqual(Q pointR1) bool {
	QQ, ok := Q.(*point448R1)
	if !ok {
		panic("wrong type")
	}
	l, r := &fp448.Elt{}, &fp448.Elt{}
	fp448.Mul(l, &P.x, &QQ.z)
	fp448.Mul(r, &QQ.x, &P.z)
	fp448.Sub(l, l, r)
	b := fp448.IsZero(l)
	fp448.Mul(l, &P.y, &QQ.z)
	fp448.Mul(r, &QQ.y, &P.z)
	fp448.Sub(l, l, r)
	b = b && fp448.IsZero(l)
	fp448.Mul(l, &P.ta, &P.tb)
	fp448.Mul(l, l, &QQ.z)
	fp448.Mul(r, &QQ.ta, &QQ.tb)
	fp448.Mul(r, r, &P.z)
	fp448.Sub(l, l, r)
	b = b && fp448.IsZero(l)
	return b
}

func (P *point448R2) neg() pointR2 {
	Q := &point448R2{}
	Q.addYX = P.subYX
	Q.subYX = P.addYX
	fp448.Neg(&Q.dt2, &P.dt2)
	Q.z2 = P.z2
	return Q
}

func (P *point448R2) fromR1(Q pointR1) {
	QQ, ok := Q.(*point448R1)
	if !ok {
		panic("wrong type")
	}
	P.point448R3.fromR1(QQ)
	fp448.Add(&P.z2, &QQ.z, &QQ.z)
}

func (P *point448R3) neg() pointR3 {
	Q := &point448R3{}
	Q.addYX = P.subYX
	Q.subYX = P.addYX
	fp448.Neg(&Q.dt2, &P.dt2)
	return Q
}

func (P *point448R3) cneg(b int) {
	t := &fp448.Elt{}
	fp448.Cswap(&P.addYX, &P.subYX, uint(b))
	fp448.Neg(t, &P.dt2)
	fp448.Cmov(&P.dt2, t, uint(b))
}

func (P *point448R3) cmov(Q pointR3, b int) {
	QQ, ok := Q.(*point448R3)
	if !ok {
		panic("wrong type")
	}
	fp448.Cmov(&P.addYX, &QQ.addYX, uint(b))
	fp448.Cmov(&P.subYX, &QQ.subYX, uint(b))
	fp448.Cmov(&P.dt2, &QQ.dt2, uint(b))
}

func (P *point448R3) fromR1(Q pointR1) {
	QQ, ok := Q.(*point448R1)
	if !ok {
		panic("wrong type")
	}
	var d fp448.Elt
	copy(d[:], edwards448.paramD)
	QQ.toAffine()
	fp448.Add(&P.addYX, &QQ.y, &QQ.x)
	fp448.Sub(&P.subYX, &QQ.y, &QQ.x)
	fp448.Mul(&P.dt2, &QQ.ta, &QQ.tb)
	fp448.Mul(&P.dt2, &P.dt2, &d)
	fp448.Add(&P.dt2, &P.dt2, &P.dt2)
}
