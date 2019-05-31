package eddsa

import (
	"crypto/rand"
	mrand "math/rand"
	"testing"

	"github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/internal/test"
)

func TestDevel(t *testing.T) {
	t.FailNow()
	var P pointR1 = &point255R1{}
	var k [32]byte
	var l [32]byte
	_, _ = mrand.Read(k[:])
	_, _ = mrand.Read(l[:])
	// for i := range l {
	// 	k[i] = 0
	// 	l[i] = 0
	// }
	// k[31] = 3

	t.Logf("k: %v\n", conv.BytesLe2Hex(k[:]))
	t.Logf("l: %v\n", conv.BytesLe2Hex(l[:]))
	P.SetGenerator()
	P.double()
	t.Logf("P: %v\n", P)
	edwards25519.doubleMult(P, k[:], l[:])
	P.toAffine()
	t.Logf("P: %v\n", P)
}

func randomPoint(e *curve, P pointR1) {
	k := make([]byte, (e.b+7)/8)
	_, _ = rand.Read(k[:])
	_ = e.fixedMult(k)
}

func TestPoint(t *testing.T) {
	t.Run("ed25519", func(t *testing.T) {
		testPoint(t, edwards25519)
	})
	t.Run("ed448", func(t *testing.T) {
		testPoint(t, edwards448)
	})
}

func testPoint(t *testing.T, c *curve) {
	testTimes := 1 << 10

	t.Run("mixAdd", func(t *testing.T) {
		P := c.newPointR1()
		Q := c.newPointR1()
		R := c.newPointR3()
		S := c.newPointR2()
		for i := 0; i < testTimes; i++ {
			randomPoint(c, P)
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

			// 16P = P+P...+P
			Q.SetIdentity()
			for j := 0; j < 16; j++ {
				Q.add(S)
			}

			got = _16P.isEqual(Q)
			want = true
			if got != want {
				test.ReportError(t, got, want, P)
			}
		}
	})
}

func BenchmarkPoint(b *testing.B) {
	b.Run("ed25519", func(b *testing.B) {
		benchmarkPoint(b, edwards25519)
	})
	b.Run("ed448", func(b *testing.B) {
		benchmarkPoint(b, edwards448)
	})
}

func benchmarkPoint(b *testing.B, c *curve) {
	k := make([]byte, (c.b+7)/8)
	l := make([]byte, (c.b+7)/8)
	_, _ = rand.Read(k)
	_, _ = rand.Read(l)

	P := c.newPointR1()
	Q := c.newPointR2()
	R := c.newPointR3()
	P.SetGenerator()
	Q.fromR1(P)
	R.fromR1(P)
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
	b.Run("add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P.add(Q)
		}
	})
	b.Run("fixedMult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = c.fixedMult(k)
		}
	})
	b.Run("doubleMult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.doubleMult(P, k, l)
		}
	})
}
