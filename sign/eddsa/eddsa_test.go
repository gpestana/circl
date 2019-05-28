package eddsa

import (
	"crypto/rand"
	"testing"
)

func BenchmarkEdDSA(b *testing.B) {
	var scheme Ed25519
	var public Pk25519
	var private Sk25519
	var sig Sig25519
	msg := make([]byte, 256)
	_, _ = rand.Read(msg)

	b.Run("keygen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scheme.KeyGen(&public, &private)
		}
	})
	b.Run("sign", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = scheme.Sign(msg, &private)
		}
	})
	b.Run("verify", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scheme.Verify(msg, &public, &sig)
		}
	})
}
