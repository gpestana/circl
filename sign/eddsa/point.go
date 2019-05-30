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

type pointR2 interface {
	SetIdentity()
	copy() pointR2
	neg() pointR2
	double()
	add(pointR2)
	oddMultiples([]pointR2)
}

type point255R1 struct{ x, y, z, ta, tb fp255.Elt }
type point255R2 struct{ x, y, z, t fp255.Elt }
type point255R3 struct{ addYX, subYX, dt2 fp255.Elt }
type point448R1 struct{ x, y, z, ta, tb fp448.Elt }
type point448R2 struct{ x, y, z, t fp448.Elt }
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

func (P *point255R2) toAffine() {
	fp255.Inv(&P.z, &P.z)
	fp255.Mul(&P.x, &P.x, &P.z)
	fp255.Mul(&P.y, &P.y, &P.z)
	fp255.Mul(&P.t, &P.x, &P.y)
	fp255.Modp(&P.x)
	fp255.Modp(&P.y)
	fp255.Modp(&P.t)
	fp255.SetOne(&P.z)
}

func (P *point255R2) ToBytes(k []byte) {
	P.toAffine()
	var x [32]byte
	fp255.ToBytes(k, &P.y)
	fp255.ToBytes(x[:], &P.x)
	b := x[0] & 1
	k[31] = k[31] | (b << 7)
}

func (P *point255R2) FromBytes(k []byte) bool {
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
	fp255.Mul(&P.t, &P.x, &P.y)
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

func (P *point255R2) copy() pointR2 { Q := *P; return &Q }

func (P *point255R2) SetIdentity() {
	fp255.SetZero(&P.x)
	fp255.SetOne(&P.y)
	fp255.SetZero(&P.t)
	fp255.SetOne(&P.z)
}
func (P *point255R2) neg() pointR2 {
	Q := &point255R2{}
	fp255.Neg(&Q.x, &P.x)
	fp255.Neg(&Q.t, &P.t)
	Q.y = P.y
	Q.z = P.z
	return Q
}
func (P *point255R2) double() {
	Px, Py, Pz, Pt := &P.x, &P.y, &P.z, &P.t
	a := Px
	b := Py
	c := Pz
	d := Pt
	e := &fp255.Elt{}
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
	fp255.Mul(Pt, e, d)
}

func (P *point255R2) add(Q pointR2) {
	QQ, ok := Q.(*point255R2)
	if !ok {
		panic("wrong type")
	}
	_2D := &fp255.Elt{}
	x1, y1, z1, t1 := &P.x, &P.y, &P.z, &P.t
	x2, y2, z2, t2 := &QQ.x, &QQ.y, &QQ.z, &QQ.t
	a := &fp255.Elt{}
	b := &fp255.Elt{}
	c := &fp255.Elt{}
	d := &fp255.Elt{}
	e := t1
	f := x1
	g := d
	h := y1
	fp255.AddSub(y1, x1)
	fp255.AddSub(y2, x2)
	fp255.Mul(a, x1, x2)
	fp255.Mul(b, y1, y2)
	fp255.Mul(c, t1, t2)
	fp255.Mul(c, c, _2D)
	fp255.Mul(d, z1, z2)
	fp255.Add(d, d, d)
	fp255.Sub(e, b, a)
	fp255.Add(h, b, a)
	fp255.Sub(f, d, c)
	fp255.Add(g, d, c)
	fp255.Mul(&P.z, f, g)
	fp255.Mul(&P.x, e, f)
	fp255.Mul(&P.t, e, h)
	fp255.Mul(&P.y, g, h)
}

func (P *point255R2) oddMultiples(T []pointR2) {
	n := len(T)
	T[0] = P
	_2P := *P
	_2P.double()
	for i := 1; i < n; i++ {
		T[i] = T[i-1].copy()
		T[i].add(&_2P)
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

func (P *point448R2) SetIdentity() {
	fp448.SetZero(&P.x)
	fp448.SetOne(&P.y)
	fp448.SetZero(&P.t)
	fp448.SetOne(&P.z)
}

func (P *point448R2) copy() pointR2 { Q := *P; return &Q }

func (P *point448R2) neg() pointR2 {
	Q := &point448R2{}
	fp448.Neg(&Q.x, &P.x)
	fp448.Neg(&Q.t, &P.t)
	Q.y = P.y
	Q.z = P.z
	return Q
}

func (P *point448R2) double() {
	Px, Py, Pz, Pt := &P.x, &P.y, &P.z, &P.t
	a := Px
	b := Py
	c := Pz
	d := Pt
	e := &fp448.Elt{}
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
	fp448.Mul(Pt, e, d)
}

func (P *point448R2) add(Q pointR2) {

}

func (P *point448R2) oddMultiples(T []pointR2) {
	n := len(T)
	T[0] = P
	_2P := *P
	_2P.double()
	for i := 1; i < n; i++ {
		T[i] = T[i-1].copy()
		T[i].add(&_2P)
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
