// +build amd64

package xcurve

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/internal/test"
)

// Montgomery point doubling in projective (X:Z) coordintates.
func doubleBig(work [4]*big.Int, A24, p *big.Int) {
	x1, z1 := work[0], work[1]
	A, B, C := big.NewInt(0), big.NewInt(0), big.NewInt(0)

	A.Add(x1, z1).Mod(A, p)
	B.Sub(x1, z1).Mod(B, p)
	A.Mul(A, A)
	B.Mul(B, B)
	C.Sub(A, B)
	x1.Mul(A, B).Mod(x1, p)
	z1.Mul(C, A24).Add(z1, B).Mul(z1, C).Mod(z1, p)
}

// Equation 7 at https://eprint.iacr.org/2017/264
func diffAddBig(work [4]*big.Int, mu, p *big.Int, b uint) {
	x1, z1, x2, z2 := work[0], work[1], work[2], work[3]
	A, B := big.NewInt(0), big.NewInt(0)
	if b != 0 {
		t := new(big.Int)
		t.Set(x1)
		x1.Set(x2)
		x2.Set(t)
		t.Set(z1)
		z1.Set(z2)
		z2.Set(t)
	}
	A.Add(x1, z1)
	B.Sub(x1, z1)
	B.Mul(B, mu).Mod(B, p)
	x1.Add(A, B).Mod(x1, p)
	z1.Sub(A, B).Mod(z1, p)
	x1.Mul(x1, x1).Mul(x1, z2).Mod(x1, p)
	z1.Mul(z1, z1).Mul(z1, x2).Mod(z1, p)
	x2.Mod(x2, p)
	z2.Mod(z2, p)
}

func ladderStepBig(work [5]*big.Int, A24, p *big.Int, b uint) {
	x1 := work[0]
	x2, z2 := work[1], work[2]
	x3, z3 := work[3], work[4]
	A, B, C, D := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	DA, CB, E := big.NewInt(0), big.NewInt(0), big.NewInt(0)
	A.Add(x2, z2).Mod(A, p)
	B.Sub(x2, z2).Mod(B, p)
	C.Add(x3, z3).Mod(C, p)
	D.Sub(x3, z3).Mod(D, p)
	DA.Mul(D, A).Mod(DA, p)
	CB.Mul(C, B).Mod(CB, p)
	if b != 0 {
		t := new(big.Int)
		t.Set(A)
		A.Set(C)
		C.Set(t)
		t.Set(B)
		B.Set(D)
		D.Set(t)
	}
	AA := A.Mul(A, A).Mod(A, p)
	BB := B.Mul(B, B).Mod(B, p)
	E.Sub(AA, BB)
	x1.Mod(x1, p)
	x2.Mul(AA, BB).Mod(x2, p)
	z2.Mul(E, A24).Add(z2, BB).Mul(z2, E).Mod(z2, p)
	x3.Add(DA, CB)
	z3.Sub(DA, CB)
	x3.Mul(x3, x3).Mod(x3, p)
	z3.Mul(z3, z3).Mul(z3, x1).Mod(z3, p)
}

func TestCurve255(t *testing.T) {
	p := big.NewInt(1)
	p.Lsh(p, 255).Sub(p, big.NewInt(19))
	testCurve(t, c255, p)
}

func TestCurve448(t *testing.T) {
	p := big.NewInt(1)
	p.Lsh(p, 224)
	p.Sub(p, new(big.Int).SetInt64(1))
	p.Lsh(p, 224)
	p.Sub(p, new(big.Int).SetInt64(1))
	testCurve(t, c448, p)
}

func testCurve(t *testing.T, c *curve, p *big.Int) {
	numTests := 1 << 9
	n := c.size
	A24 := big.NewInt(int64(c.a24))
	mu := make([]byte, 1*c.size)
	work := make([]byte, 4*c.size)
	workLadder := make([]byte, 5*c.size)
	var bigWork [4]*big.Int
	var bigWorkLadder [5]*big.Int

	t.Run("double", func(t *testing.T) {
		for i := 0; i < numTests; i++ {
			_, _ = rand.Read(work[:])
			for j := range bigWork {
				bigWork[j] = conv.BytesLe2BigInt(work[j*n : (j+1)*n])
			}

			c.double(work[:])
			got0 := conv.BytesLe2BigInt(work[0*n : 1*n])
			got1 := conv.BytesLe2BigInt(work[1*n : 2*n])
			got0.Mod(got0, p)
			got1.Mod(got1, p)

			doubleBig(bigWork, A24, p)
			want0 := bigWork[0]
			want1 := bigWork[1]

			if got0.Cmp(want0) != 0 {
				test.ReportError(t, got0, want0, work)
			}
			if got1.Cmp(want1) != 0 {
				test.ReportError(t, got1, want1, work)
			}
		}
	})

	t.Run("diffAdd", func(t *testing.T) {
		for i := 0; i < numTests; i++ {
			_, _ = rand.Read(work[:])
			for j := range bigWork {
				bigWork[j] = conv.BytesLe2BigInt(work[j*n : (j+1)*n])
			}
			_, _ = rand.Read(mu[:])
			bigMu := conv.BytesLe2BigInt(mu[:])
			b := uint(mu[0] & 1)

			c.difAdd(work[:], mu[:], b)

			diffAddBig(bigWork, bigMu, p, b)

			for j := range bigWork {
				got := conv.BytesLe2BigInt(work[j*n : (j+1)*n])
				got.Mod(got, p)
				want := bigWork[j]
				if got.Cmp(want) != 0 {
					test.ReportError(t, got, want, work, mu, b)
				}
			}
		}
	})

	t.Run("ladder", func(t *testing.T) {
		for i := 0; i < numTests; i++ {
			_, _ = rand.Read(workLadder[:])
			for j := range bigWorkLadder {
				bigWorkLadder[j] = conv.BytesLe2BigInt(workLadder[j*n : (j+1)*n])
			}
			b := uint(workLadder[0] & 1)

			c.ladderStep(workLadder[:], b)

			ladderStepBig(bigWorkLadder, A24, p, b)

			for j := range bigWorkLadder {
				got := conv.BytesLe2BigInt(workLadder[j*n : (j+1)*n])
				got.Mod(got, p)
				want := bigWorkLadder[j]
				if got.Cmp(want) != 0 {
					test.ReportError(t, got, want, workLadder, b)
				}
			}
		}
	})
}
