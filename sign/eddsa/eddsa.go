// Package eddsa provides the signature schemes as described in RFC-8032.
package eddsa

import (
	"bytes"
	"crypto/sha512"
	// "fmt"
	"golang.org/x/crypto/sha3"
	// "github.com/cloudflare/circl/internal/conv"
)

const (
	// SizeKey25519 is the size in bytes of Ed25519 keys.
	SizeKey25519 = 32
	// SizeKey448 is the size in bytes of Ed448 keys.
	SizeKey448 = 57
)

type signScheme int

const (
	schemeEd25519 signScheme = iota
	schemeEd25519ph
	schemeEd25519ctx
	schemeEd448
	schemeEd448ph
)

// Pk25519 represents a public key of Ed25519.
type Pk25519 [SizeKey25519]byte

// Sk25519 represents a private key of Ed25519.
type Sk25519 [SizeKey25519]byte

// Sig25519 represents an Ed25519 signature.
type Sig25519 [2 * SizeKey25519]byte

// Pk448 represents a public key of Ed448.
type Pk448 [SizeKey448]byte

// Sk448 represents a private key of Ed448.
type Sk448 [SizeKey448]byte

// Sig448 represents an Ed448 signature.
type Sig448 [2 * SizeKey448]byte

// Ed25519 is used to instantiate an object able to perform Ed25519 operations.
type Ed25519 struct{}

// Ed448 is used to instantiate an object able to perform Ed448 operations.
type Ed448 struct{}

// KeyGen generates a public key from a secret key.
func (e *Ed25519) KeyGen(public *Pk25519, private *Sk25519) {
	H := sha512.New()
	_, _ = H.Write(private[:])
	ah := H.Sum([]byte{})
	ah[0] &= -(uint8(1) << edwards25519.lgCofactor)
	ah[31] = (ah[31] & 127) | 64

	edwards25519.reduceModOrder(ah[:32])

	var P point255R1
	edwards25519.fixedMult(&P, ah[:32])
	P.ToBytes(public[:])
}

// Sign creates the signature of a message using both the private and public
// keys of the signer.
func (e *Ed25519) Sign(msg []byte, public *Pk25519, private *Sk25519) *Sig25519 {
	H := sha512.New()
	_, _ = H.Write(private[:])
	ah := H.Sum([]byte{})
	ah[0] &= -(uint8(1) << edwards25519.lgCofactor)
	ah[31] = (ah[31] & 127) | 64
	// fmt.Printf("")
	// fmt.Printf("ah: %x\n", ah)

	H.Reset()
	_, _ = H.Write(ah[32:])
	_, _ = H.Write(msg)
	r := H.Sum([]byte{})
	// fmt.Printf("r: %x\n", r)
	edwards25519.reduceModOrder(r[:])
	// fmt.Printf("r: %x\n", r[:32])

	var P point255R1
	edwards25519.fixedMult(&P, r[:32])
	var signature Sig25519
	P.ToBytes(signature[:32])
	// fmt.Printf("s0: %x\n", signature[:32])

	H.Reset()
	_, _ = H.Write(signature[:32])
	_, _ = H.Write(public[:])
	_, _ = H.Write(msg)
	hRAM := H.Sum([]byte{})
	// fmt.Printf("hRAM: %x\n", hRAM[:])
	edwards25519.reduceModOrder(hRAM[:])
	// fmt.Printf("hRAM: %x\n", hRAM[:32])
	edwards25519.calculateS(signature[32:], r[:32], hRAM[:32], ah[:32])
	// fmt.Printf("s1: %x\n", signature[32:])
	return &signature
}

// Verify returns false if the signature is invalid or when the public key can
// not be decoded; otherwise, returns true.
func (e *Ed25519) Verify(msg []byte, public *Pk25519, sig *Sig25519) bool {
	var P point255R1
	// fmt.Printf("pk: %x\n", public)
	if ok := P.FromBytes(public[:]); !ok {
		return false
	}
	// fmt.Printf("A: %v\n", &A)

	H := sha512.New()
	_, _ = H.Write(sig[:32])
	_, _ = H.Write(public[:])
	_, _ = H.Write(msg)
	hRAM := H.Sum([]byte{})
	// fmt.Printf("hRAM: %x\n", hRAM[:])
	edwards25519.reduceModOrder(hRAM[:])
	// fmt.Printf("s: %v\n", conv.BytesLe2Hex(sig[32:]))
	// fmt.Printf("h: %v\n", conv.BytesLe2Hex(hRAM[:32]))
	P.neg()
	// fmt.Printf("P: %v\n", &P)
	edwards25519.doubleMult(&P, sig[32:], hRAM[:32])
	// fmt.Printf("Q: %v\n", &P)
	P.toAffine()
	// fmt.Printf("aQ: %v\n", &P)
	var encP [32]byte
	P.ToBytes(encP[:])
	// fmt.Printf("encP: %x\n", encP)
	// fmt.Printf("sig0: %x\n", sig[:32])
	return bytes.Equal(encP[:], sig[:32])
}

// KeyGen generates a public key from a secret key.
func (e *Ed448) KeyGen(public *Pk448, private *Sk448) {
	var dig [114]byte
	H := sha3.NewShake256()
	_, _ = H.Write(private[:])
	_, _ = H.Read(dig[:])
	dig[0] &= -(uint8(1) << edwards448.lgCofactor)
	// dig[56] = (dig[31] & 127) | 64

	edwards448.reduceModOrder(dig[:57])

	var P point448R1
	edwards448.fixedMult(&P, dig[:57])
	P.ToBytes(public[:])
}

// Sign creates the signature of a message using both the private and public
// keys of the signer.
func (e *Ed448) Sign(msg []byte, public *Pk448, private *Sk448) *Sig448 {
	var signature Sig448
	return &signature
}

// Verify returns false if the signature is invalid or when the public key can
// not be decoded; otherwise, returns true.
func (e *Ed448) Verify(msg []byte, public *Pk448, sig *Sig448) bool {
	return false
}
