package eddsa

import (
	"crypto/subtle"
	"encoding/binary"
	// "fmt"
	"math/bits"
)

type curve struct {
	size        int
	b           int
	fixedParams struct{ t, v, w int }
	order       []uint64
	paramD      []byte
	genX, genY  []byte
	Table       *[2][_2w1]pointR3
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

func (ecc *curve) fixedMult(P pointR1, S pointR3, scalar []byte) {
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
	for ii := e - 1; ii >= 0; ii-- {
		P.double()
		for j := 0; j < fx_v; j++ {
			dig := L[fx_w*d-j*e+ii-e]
			for i := (fx_w-1)*d - j*e + ii - e; i >= (2*d - j*e + ii - e); i = i - d {
				dig = 2*dig + L[i]
			}
			idx := absolute(int32(dig))
			sig := L[d-j*e+ii-e]
			Tabj := &ecc.Table[fx_v-j-1]
			for k := 0; k < fx_2w1; k++ {
				S.cmov(Tabj[k], subtle.ConstantTimeEq(int32(k), int32(idx)))
			}
			S.cneg(subtle.ConstantTimeEq(int32(sig), -1))
			P.mixAdd(S)
		}
	}
}

// absolute returns always a positive value.
func absolute(x int32) int32 {
	mask := x >> 31
	return (x + mask) ^ mask
}
