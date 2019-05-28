package eddsa

import (
	"fmt"

	fp255 "github.com/cloudflare/circl/math/fp25519"
	"github.com/cloudflare/circl/math/fp448"
)

type pointR3 interface {
	cneg(int)
	cmov(pointR3, int)
	fromR1(pointR1)
}

type pointR1 interface {
	copy() pointR1
	SetIdentity()
	SetGenerator()
	isEqual(pointR1) bool
	toAffine()
	double()
	mixAdd(pointR3)
}

type point255R1 struct{ x, y, z, ta, tb fp255.Elt }
type point255R3 struct{ addYX, subYX, dt2 fp255.Elt }
type point448R1 struct{ x, y, z, ta, tb fp448.Elt }
type point448R3 struct{ addYX, subYX, dt2 fp448.Elt }

func (P *point255R1) String() string {
	return fmt.Sprintf("\nx=  %v\ny=  %v\nta= %v\ntb= %v\nz=  %v",
		P.x, P.y, P.ta, P.tb, P.z)
}
func (P *point255R3) String() string {
	return fmt.Sprintf("\naddYX= %v\nsubYX= %v\ndt2=  %v",
		P.addYX, P.subYX, P.dt2)
}
func (P *point448R1) String() string {
	return fmt.Sprintf("\nx=  %v\ny=  %v\nta= %v\ntb= %v\nz=  %v",
		P.x, P.y, P.ta, P.tb, P.z)
}
func (P *point448R3) String() string {
	return fmt.Sprintf("\naddYX= %v\nsubYX= %v\ndt2=  %v",
		P.addYX, P.subYX, P.dt2)
}

func (P *point255R1) ToBytes(k []byte) {
	P.toAffine()
	var x [32]byte
	fp255.ToBytes(k, &P.y)
	fp255.ToBytes(x[:], &P.x)
	b := x[0] & 1
	k[31] = k[31] | (b << 7)
}

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

func (P *point255R1) copy() pointR1 { Q := *P; return &Q }

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

func (P *point448R1) ToBytes(k []byte) {
}

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

func (P *point448R1) double() {
	//formula is different for a=1
	Px, Py, Pz, Pta, Ptb := &P.x, &P.y, &P.z, &P.ta, &P.tb
	a := Px
	b := Py
	c := Pz
	h := Pta
	e := Ptb
	f := a
	g := b
	fp448.Add(e, Px, Py)
	fp448.Sqr(a, Px)
	fp448.Sqr(b, Py)
	fp448.Sqr(c, Pz)
	fp448.Add(c, c, c)
	fp448.Add(h, a, b)
	fp448.Neg(h, h)
	fp448.Sqr(e, e)
	fp448.Add(e, e, h)
	fp448.Sub(g, b, a)
	fp448.Sub(f, g, c)
	fp448.Mul(Pz, f, g)
	fp448.Mul(Px, e, f)
	fp448.Mul(Py, g, h)
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

func (P *point448R1) copy() pointR1 { Q := *P; return &Q }

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
