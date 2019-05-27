package eddsa

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/internal/test"
)

func TestScalar(t *testing.T) {
	testScalar(t, edwards25519)
	testScalar(t, edwards448)
}

func testScalar(t *testing.T, c *curve) {
	testTimes := 1 << 12
	x := make([]uint64, len(c.order)+1)
	xx := make([]uint64, len(c.order)+1)
	max := big.NewInt(1)
	max.Lsh(max, uint(c.size))
	two64 := big.NewInt(1)
	two64.Lsh(two64, 64)
	bigOrder := conv.Uint64Le2BigInt(c.order[:])

	t.Run("div2subY", func(t *testing.T) {
		want := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigX, _ := rand.Int(rand.Reader, max)
			conv.BigInt2Uint64Le(xx[:], bigX)
			copy(x, xx)

			bigY, _ := rand.Int(rand.Reader, two64)
			y := bigY.Int64()
			bigY.SetInt64(y)

			c.div2subY(x[:], y, len(c.order))
			got := conv.Uint64Le2BigInt(x[:])

			want.Rsh(bigX, 1).Sub(want, bigY)

			if got.Cmp(want) != 0 {
				test.ReportError(t, got, want, bigX, y)
			}
		}
	})

	t.Run("condAddOrderN", func(t *testing.T) {
		for i := 0; i < testTimes; i++ {
			bigX, _ := rand.Int(rand.Reader, max)
			conv.BigInt2Uint64Le(xx[:], bigX)
			copy(x, xx)
			c.condAddOrderN(x[:])
			got := conv.Uint64Le2BigInt(x[:])

			want := new(big.Int).Set(bigX)
			if want.Bit(0) == 0 {
				want.Add(want, bigOrder)
			}

			if got.Cmp(want) != 0 {
				test.ReportError(t, got, want, bigX)
			}
		}
	})
}
