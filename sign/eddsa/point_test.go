package eddsa

import (
	"crypto/rand"
	"testing"

	// "github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/internal/test"
)

//
// func TestDevel(t *testing.T) {
// 	var P pointR1 = &point255R1{}
// 	var Q pointR3 = &point255R3{}
// 	var k [32]byte
// 	_, _ = rand.Read(k[:])
//
// 	t.Logf("k: %v\n", conv.BytesLe2Hex(k[:]))
// 	edwards25519.fixedMult(P, Q, k[:])
// 	P.toAffine()
// 	t.Logf("P: %v\n", P)
// }

func randomPoint(e *curve, P pointR1, Q pointR3) {
	k := make([]byte, e.size)
	_, _ = rand.Read(k[:])
	e.fixedMult(P, Q, k)
}

func TestPoint(t *testing.T) {
	var P, Q pointR1
	var R pointR3
	t.Run("ed25519", func(t *testing.T) {
		P = &point255R1{}
		Q = &point255R1{}
		R = &point255R3{}
		testPoint(t, P, Q, R, edwards25519)
	})
}

func testPoint(t *testing.T, P, Q pointR1, R pointR3, c *curve) {
	testTimes := 1 << 10
	t.Run("addition", func(t *testing.T) {
		for i := 0; i < testTimes; i++ {
			randomPoint(c, P, R)
			_16P := P.copy()
			R.fromR1(P)
			// 16P = 2^4P
			for j := 0; j < 4; j++ {
				_16P.double()
			}
			// 16P = P+P...+P
			Q.SetIdentity()
			for j := 0; j < 16; j++ {
				Q.mixAdd(R)
			}

			got := _16P.isEqual(Q)
			want := true
			if got != want {
				test.ReportError(t, got, want, P)
			}
		}
	})
}

func BenchmarkPoint(b *testing.B) {
	var P, Q pointR1
	var R pointR3
	b.Run("ed25519", func(b *testing.B) {
		P = &point255R1{}
		Q = &point255R1{}
		R = &point255R3{}
		benchmarkPoint(b, P, Q, R, edwards25519)
	})
	// b.Run("ed448", func(b *testing.B) {
	// 	P = &point448R1{}
	// 	Q = &point448R3{}
	// 	R = &point448R3{}
	// 	benchmarkPoint(b, P, Q, R, edwards448)
	// })
}

func benchmarkPoint(b *testing.B, P, Q pointR1, R pointR3, c *curve) {
	k := make([]byte, (c.size+7)/8)
	l := make([]byte, (c.size+7)/8)
	_, _ = rand.Read(k)
	P.SetGenerator()
	b.Run("toAffine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.toAffine()
		}
	})
	b.Run("double", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.double()
		}
	})
	b.Run("mixadd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.mixAdd(R)
		}
	})
	b.Run("fixedMult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.fixedMult(P, R, k)
		}
	})
	b.Run("doubleMult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.doubleMult(P, Q, k, l)
		}
	})
}
