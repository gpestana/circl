// Package ed25519 provides the signature scheme Ed25519 as described in RFC-8032.
package ed25519

import (
	"bytes"
	"crypto/sha512"
	// "fmt"
	"github.com/cloudflare/circl/internal/conv"
)

// Size is the size in bytes of Ed25519 keys.
const Size = 32

type signScheme int

const (
	schemeEd25519 signScheme = iota
	schemeEd25519ph
	schemeEd25519ctx
)

// Pk25519 represents a public key of Ed25519.
type Pk25519 [Size]byte

// Sk25519 represents a private key of Ed25519.
type Sk25519 [Size]byte

// Sig25519 represents an Ed25519 signature.
type Sig25519 [2 * Size]byte

// Ed25519 is used to instantiate an object able to perform Ed25519 operations.
type Ed25519 struct{}

// KeyGen generates a public key from a secret key.
func (e *Ed25519) KeyGen(public *Pk25519, private *Sk25519) {
	k := sha512.Sum512(private[:])
	e.clamp(k[:])
	e.reduceModOrder(k[:Size])
	var P pointR1
	P.fixedMult(k[:Size])
	P.ToBytes(public[:])
}

// Sign creates the signature of a message using both the private and public
// keys of the signer.
func (e *Ed25519) Sign(msg []byte, public *Pk25519, private *Sk25519) *Sig25519 {
	k := sha512.Sum512(private[:])
	e.clamp(k[:])
	// fmt.Printf("")
	// fmt.Printf("ah: %x\n", ah)
	H := sha512.New()
	var r [sha512.Size]byte
	_, _ = H.Write(k[Size:])
	_, _ = H.Write(msg)
	H.Sum(r[:0])
	// fmt.Printf("r: %x\n", r)
	e.reduceModOrder(r[:])
	// fmt.Printf("r: %x\n", r[:32])
	var P pointR1
	P.fixedMult(r[:Size])
	signature := &Sig25519{}
	P.ToBytes(signature[:Size])
	// fmt.Printf("s0: %x\n", signature[:32])
	var hRAM [sha512.Size]byte
	H.Reset()
	_, _ = H.Write(signature[:Size])
	_, _ = H.Write(public[:])
	_, _ = H.Write(msg)
	H.Sum(hRAM[:0])
	// fmt.Printf("hRAM: %x\n", hRAM[:])
	e.reduceModOrder(hRAM[:])
	// fmt.Printf("hRAM: %x\n", hRAM[:32])
	// fmt.Printf("s1: %x\n", signature[32:])
	e.calculateS(signature[Size:], r[:Size], hRAM[:Size], k[:Size])
	return signature
}

// Verify returns false if the signature is invalid or when the public key can
// not be decoded; otherwise, returns true.
func (e *Ed25519) Verify(msg []byte, public *Pk25519, sig *Sig25519) bool {
	var P pointR1
	// fmt.Printf("pk: %x\n", public)
	if ok := P.FromBytes(public[:]); !ok {
		return false
	}
	P.neg()
	// fmt.Printf("A: %v\n", &A)

	hRAM := [sha512.Size]byte{}
	H := sha512.New()
	_, _ = H.Write(sig[:Size])
	_, _ = H.Write(public[:])
	_, _ = H.Write(msg)
	H.Sum(hRAM[:0])
	// fmt.Printf("hRAM: %x\n", hRAM[:])
	e.reduceModOrder(hRAM[:])
	// fmt.Printf("s: %v\n", conv.BytesLe2Hex(sig[32:]))
	// fmt.Printf("h: %v\n", conv.BytesLe2Hex(hRAM[:32]))
	// fmt.Printf("P: %v\n", &P)
	if ok := e.verifyS(sig[Size:]); !ok {
		return false
	}
	var Q pointR1
	Q.doubleMult(&P, sig[Size:], hRAM[:Size])
	// fmt.Printf("Q: %v\n", &P)
	// fmt.Printf("aQ: %v\n", &P)
	var enc [Size]byte
	Q.ToBytes(enc[:])
	// fmt.Printf("encP: %x\n", encP)
	// fmt.Printf("sig0: %x\n", sig[:32])
	return bytes.Equal(enc[:], sig[:Size])
}

func (e *Ed25519) clamp(k []byte) {
	k[0] &= 248
	k[Size-1] = (k[Size-1] & 127) | 64
}

func (e *Ed25519) reduceModOrder(k []byte) {
	bigK := conv.BytesLe2BigInt(k)
	orderBig := conv.Uint64Le2BigInt(curve.order[:])
	bigK.Mod(bigK, orderBig)
	conv.BigInt2BytesLe(k, bigK)
}

// calculateS performs s= r+k*a mod L
func (e *Ed25519) calculateS(s, r, k, a []byte) {
	R := conv.BytesLe2BigInt(r)
	K := conv.BytesLe2BigInt(k)
	A := conv.BytesLe2BigInt(a)
	order := conv.Uint64Le2BigInt(curve.order[:])
	S := K.Mul(K, A).Add(K, R)
	S.Mod(S, order)
	conv.BigInt2BytesLe(s, S)
}

func (e *Ed25519) verifyS([]byte) bool {
	return true
}
