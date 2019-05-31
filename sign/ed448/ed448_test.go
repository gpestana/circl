package ed448

import (
	"crypto/rand"
	"testing"

	"github.com/cloudflare/circl/internal/test"
)

type vector struct {
	name   string
	scheme signScheme
	sk     []byte
	pk     []byte
	sig    []byte
	msg    []byte
	msgLen uint
	ctx    []byte
	ctxLen uint
}

var vectors = [...]vector{
	vector{
		name:   "-----TEST 1",
		scheme: schemeEd448,
		sk: []byte{
			0x9d, 0x61, 0xb1, 0x9d, 0xef, 0xfd, 0x5a, 0x60, 0xba, 0x84, 0x4a, 0xf4, 0x92, 0xec, 0x2c, 0xc4,
			0x44, 0x49, 0xc5, 0x69, 0x7b, 0x32, 0x69, 0x19, 0x70, 0x3b, 0xac, 0x03, 0x1c, 0xae, 0x7f, 0x60,
		},
		pk: []byte{
			0xd7, 0x5a, 0x98, 0x01, 0x82, 0xb1, 0x0a, 0xb7, 0xd5, 0x4b, 0xfe, 0xd3, 0xc9, 0x64, 0x07, 0x3a,
			0x0e, 0xe1, 0x72, 0xf3, 0xda, 0xa6, 0x23, 0x25, 0xaf, 0x02, 0x1a, 0x68, 0xf7, 0x07, 0x51, 0x1a,
		},
		msg:    []byte{},
		msgLen: 0,
		sig: []byte{
			0xe5, 0x56, 0x43, 0x00, 0xc3, 0x60, 0xac, 0x72, 0x90, 0x86, 0xe2, 0xcc, 0x80, 0x6e, 0x82, 0x8a,
			0x84, 0x87, 0x7f, 0x1e, 0xb8, 0xe5, 0xd9, 0x74, 0xd8, 0x73, 0xe0, 0x65, 0x22, 0x49, 0x01, 0x55,
			0x5f, 0xb8, 0x82, 0x15, 0x90, 0xa3, 0x3b, 0xac, 0xc6, 0x1e, 0x39, 0x70, 0x1c, 0xf9, 0xb4, 0x6b,
			0xd2, 0x5b, 0xf5, 0xf0, 0x59, 0x5b, 0xbe, 0x24, 0x65, 0x51, 0x41, 0x43, 0x8e, 0x7a, 0x10, 0x0b,
		},
		ctx:    []byte{},
		ctxLen: 0,
	},
}

func TestEd448(t *testing.T) {
	var scheme Ed448
	var public Pk
	var private Sk
	t.Run("keygen", func(t *testing.T) {
		var want Pk
		for _, v := range vectors {
			copy(private[:], v.sk)
			scheme.KeyGen(&public, &private)
			got := public
			copy(want[:], v.pk)

			if got != want {
				test.ReportError(t, got, want, v.sk)
			}
		}
	})

	t.Run("sign", func(t *testing.T) {
		var want Sig
		for _, v := range vectors {
			copy(private[:], v.sk)
			copy(public[:], v.pk)
			sig := scheme.Sign(v.msg, &public, &private)
			got := *sig
			copy(want[:], v.sig)

			if got != want {
				test.ReportError(t, got, want, v.name)
			}
		}
	})

	t.Run("Verify", func(t *testing.T) {
		var sig Sig
		for _, v := range vectors {
			copy(private[:], v.sk)
			copy(public[:], v.pk)
			copy(sig[:], v.sig)
			got := scheme.Verify(v.msg, &public, &sig)
			want := true

			if got != want {
				test.ReportError(t, got, want, v.name)
			}
		}
	})
}

func BenchmarkEd448(b *testing.B) {
	var scheme Ed448
	var public Pk
	var private Sk
	var sig Sig
	msg := make([]byte, Size*8)
	_, _ = rand.Read(msg)

	b.Run("keygen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scheme.KeyGen(&public, &private)
		}
	})
	b.Run("sign", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheme.Sign(msg, &public, &private)
		}
	})
	b.Run("verify", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scheme.Verify(msg, &public, &sig)
		}
	})
}
