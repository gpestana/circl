package eddsa

import (
	"crypto/subtle"
	"encoding/binary"
	// "fmt"
	"math/bits"

	"github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/math"
)

type pointR2 interface {
	neg() pointR2
	fromR1(pointR1)
}

type pointR3 interface {
	neg() pointR3
	cneg(int)
	cmov(pointR3, int)
	fromR1(pointR1)
}

type pointR1 interface {
	neg()
	copy() pointR1
	SetIdentity()
	SetGenerator()
	isEqual(pointR1) bool
	toAffine()
	double()
	mixAdd(pointR3)
	add(pointR2)
	oddMultiples([]pointR2)
}

const (
	idEd25519 = iota
	idEd448
)

type curve struct {
	id          int
	b           int
	lgCofactor  uint
	fixedParams struct{ t, v, w int }
	order       []uint64
	paramD      []byte
	genX, genY  []byte
	TabSign     *[2][_2w1]pointR3
	TabVerif    *[numPointsVerif]pointR3
}

func (ecc *curve) newPointR1() pointR1 {
	if ecc.id == idEd25519 {
		return &point255R1{}
	}
	if ecc.id == idEd448 {
		return &point448R1{}
	}
	return nil
}

func (ecc *curve) newPointR2() pointR2 {
	if ecc.id == idEd25519 {
		return &point255R2{}
	}
	if ecc.id == idEd448 {
		return &point448R2{}
	}
	return nil
}

func (ecc *curve) newPointR3() pointR3 {
	if ecc.id == idEd25519 {
		return &point255R3{}
	}
	if ecc.id == idEd448 {
		return &point448R3{}
	}
	return nil
}

// condAddOrderN updates x = x+order if x is even, otherwise x remains unchanged
func (ecc *curve) condAddOrderN(x []uint64) {
	if len(x) != len(ecc.order)+1 {
		panic("wrong size")
	}
	isOdd := uint64((x[0] & 0x1) - 1)
	c := uint64(0)
	for i := range ecc.order {
		o := isOdd & ecc.order[i]
		x0, c0 := bits.Add64(x[i], o, c)
		x[i] = x0
		c = c0
	}
	x[len(ecc.order)], _ = bits.Add64(x[len(ecc.order)], 0, c)
}

// div2subY update x = (x/2) - y
func (ecc *curve) div2subY(x []uint64, y int64, l int) {
	s := uint64(y >> 63)
	for i := 0; i < l-1; i++ {
		x[i] = (x[i] >> 1) | (x[i+1] << 63)
	}
	x[l-1] = (x[l-1] >> 1)

	b := uint64(0)
	x0, b0 := bits.Sub64(x[0], uint64(y), b)
	x[0] = x0
	b = b0
	for i := 1; i < l-1; i++ {
		x0, b0 := bits.Sub64(x[i], s, b)
		x[i] = x0
		b = b0
	}
	x[l-1], _ = bits.Sub64(x[l-1], s, b)
}

// mLSBRecoding is the odd-only modified LSB-set.
//
// Reference:
//  "Efficient and secure algorithms for GLV-based scalar multiplication and
//   their implementation on GLV–GLS curves" by (Faz-Hernandez et al.)
//   http://doi.org/10.1007/s13389-014-0085-7
func (ecc *curve) mLSBRecoding(L []int8, k []byte) {
	fx_t := ecc.fixedParams.t
	fx_v := ecc.fixedParams.v
	fx_w := ecc.fixedParams.w
	e := (fx_t + fx_w*fx_v - 1) / (fx_w * fx_v)
	d := e * fx_v
	l := d * fx_w
	if len(L) == (l + 1) {
		m := make([]uint64, len(ecc.order)+1)
		for i := 0; i < len(ecc.order); i++ {
			m[i] = binary.LittleEndian.Uint64(k[8*i : 8*i+8])
		}
		ecc.condAddOrderN(m[:])
		L[d-1] = 1
		for i := 0; i < d-1; i++ {
			kip1 := (m[(i+1)/64] >> (uint(i+1) % 64)) & 0x1
			L[i] = int8(kip1<<1) - 1
		}
		{ // right-shift by d
			right := uint(d % 64)
			left := uint(64 - right)
			lim := ((len(ecc.order)+1)*64 - d) / 64
			j := d / 64
			for i := 0; i < lim; i++ {
				m[i] = (m[i+j] >> right) | (m[i+j+1] << left)
			}
			m[lim] = m[lim+j] >> right
		}
		for i := d; i < l; i++ {
			L[i] = L[i%d] * int8(m[0]&0x1)
			ecc.div2subY(m[:], int64(L[i]>>1), 4)
		}
		L[l] = int8(m[0])
	}
}

func (ecc *curve) fixedMult(P pointR1, scalar []byte) {
	fx_t := ecc.fixedParams.t
	fx_v := ecc.fixedParams.v
	fx_w := ecc.fixedParams.w
	var e = (fx_t + fx_w*fx_v - 1) / (fx_w * fx_v)
	var d = e * fx_v
	var l = d * fx_w
	var fx_2w1 = 1 << (uint(fx_w) - 1)

	L := make([]int8, l+1)
	ecc.mLSBRecoding(L[:], scalar)
	P.SetIdentity()
	S := ecc.newPointR3()
	for ii := e - 1; ii >= 0; ii-- {
		P.double()
		for j := 0; j < fx_v; j++ {
			dig := L[fx_w*d-j*e+ii-e]
			for i := (fx_w-1)*d - j*e + ii - e; i >= (2*d - j*e + ii - e); i = i - d {
				dig = 2*dig + L[i]
			}
			idx := absolute(int32(dig))
			sig := L[d-j*e+ii-e]
			Tabj := &ecc.TabSign[fx_v-j-1]
			for k := 0; k < fx_2w1; k++ {
				S.cmov(Tabj[k], subtle.ConstantTimeEq(int32(k), int32(idx)))
			}
			S.cneg(subtle.ConstantTimeEq(int32(sig), -1))
			P.mixAdd(S)
		}
	}
}

// doubleMult calculates P = mP+nG
func (ecc *curve) doubleMult(P pointR1, m, n []byte) {
	nafFix := math.OmegaNAF(conv.BytesLe2BigInt(m), omegaFix)
	nafVar := math.OmegaNAF(conv.BytesLe2BigInt(n), omegaVar)

	if len(nafFix) > len(nafVar) {
		nafVar = append(nafVar, make([]int32, len(nafFix)-len(nafVar))...)
	} else if len(nafFix) < len(nafVar) {
		nafFix = append(nafFix, make([]int32, len(nafVar)-len(nafFix))...)
	}
	// fmt.Printf("nafVar[")
	// for i := range nafVar {
	// 	fmt.Printf("%v, ", nafVar[i])
	// }
	// fmt.Printf("]\n")
	// fmt.Printf("nafFix[")
	// for i := range nafFix {
	// 	fmt.Printf("%v, ", nafFix[i])
	// }
	// fmt.Printf("]\n")

	var TabP [1 << (omegaVar - 2)]pointR2
	// fmt.Println("doubleMult")
	P.oddMultiples(TabP[:])
	// P is now used as an output value
	P.SetIdentity()
	for i := len(nafFix) - 1; i >= 0; i-- {
		// fmt.Printf("i:%v\n", i)
		P.double()
		// Generator point
		if nafFix[i] != 0 {
			idxM := absolute(nafFix[i]) >> 1
			R := ecc.TabVerif[idxM]
			if nafFix[i] < 0 {
				R = R.neg()
			}
			P.mixAdd(R)
		}
		// Variable input point
		if nafVar[i] != 0 {
			idxN := absolute(nafVar[i]) >> 1
			Q := TabP[idxN]
			if nafVar[i] < 0 {
				Q = Q.neg()
			}
			P.add(Q)
		}
	}
}

func (ecc *curve) reduceModOrder(k []byte) {
	bigK := conv.BytesLe2BigInt(k)
	orderBig := conv.Uint64Le2BigInt(ecc.order[:])
	bigK.Mod(bigK, orderBig)
	conv.BigInt2BytesLe(k, bigK)
}

// calculateS performs s= r+k*a mod L
func (ecc *curve) calculateS(s, r, k, a []byte) {
	R := conv.BytesLe2BigInt(r)
	K := conv.BytesLe2BigInt(k)
	A := conv.BytesLe2BigInt(a)
	order := conv.Uint64Le2BigInt(ecc.order[:])
	S := K.Mul(K, A).Add(K, R)
	S.Mod(S, order)
	conv.BigInt2BytesLe(s, S)
}

// absolute returns always a positive value.
func absolute(x int32) int32 {
	mask := x >> 31
	return (x + mask) ^ mask
}
