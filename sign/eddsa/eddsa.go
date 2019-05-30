// Package eddsa provides the signature schemes as described in RFC-8032.
package eddsa

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/sha3"
)

const (
	SizeEd25519 = 32
	SizeEd448   = 57
)

type EDDSA_SCHEME int

const (
	ED25519 EDDSA_SCHEME = iota
	ED25519ph
	ED25519ctx
	ED448
	ED448ph
)

type Pk25519 [SizeEd25519]byte
type Sk25519 [SizeEd25519]byte
type Sig25519 [2 * SizeEd25519]byte

type Pk448 [SizeEd448]byte
type Sk448 [SizeEd448]byte
type Sig448 [2 * SizeEd448]byte

type Ed25519 struct{}
type Ed448 struct{}

func (e *Ed25519) KeyGen(public *Pk25519, private *Sk25519) {
	H := sha512.New()
	_, _ = H.Write(private[:])
	ah := H.Sum([]byte{})
	ah[0] &= -(uint8(1) << edwards25519.lgCofactor)
	ah[31] = (ah[31] & 127) | 64

	edwards25519.reduceModOrder(ah[:32])

	var P point255R1
	var S point255R3
	edwards25519.fixedMult(&P, &S, ah[:32])
	P.ToBytes(public[:])
}

var prefix = "SigEd25519 no Ed25519 collisions"

func (e *Ed25519) Sign(msg []byte, public *Pk25519, private *Sk25519) *Sig25519 {
	H := sha512.New()
	_, _ = H.Write(private[:])
	ah := H.Sum([]byte{})
	ah[0] &= -(uint8(1) << edwards25519.lgCofactor)
	ah[31] = (ah[31] & 127) | 64
	// fmt.Printf("")
	// fmt.Printf("ah: %x\n", ah)

	H = sha512.New()
	_, _ = H.Write(ah[32:])
	_, _ = H.Write(msg)
	r := H.Sum([]byte{})
	// fmt.Printf("r: %x\n", r)
	edwards25519.reduceModOrder(r[:])
	// fmt.Printf("r: %x\n", r[:32])

	var P point255R1
	var S point255R3
	edwards25519.fixedMult(&P, &S, r[:32])
	var signature Sig25519
	P.ToBytes(signature[:32])
	// fmt.Printf("s0: %x\n", signature[:32])

	H = sha512.New()
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

func (e *Ed25519) Verify(msg []byte, public *Pk25519, sig *Sig25519) bool {
	var A point255R1
	fmt.Printf("pk: %x\n", public)
	if ok := A.FromBytes(public[:]); !ok {
		return false
	}
	fmt.Printf("A: %v\n", &A)

	H := sha512.New()
	_, _ = H.Write(sig[:32])
	_, _ = H.Write(public[:])
	_, _ = H.Write(msg)
	hRAM := H.Sum([]byte{})
	fmt.Printf("hRAM: %x\n", hRAM[:])
	edwards25519.reduceModOrder(hRAM[:])
	fmt.Printf("hRAM: %x\n", hRAM[:32])
	var P point255R1
	edwards25519.doubleMult(&P, &A, sig[:32], hRAM[:32])
	var encP [32]byte
	P.ToBytes(encP[:])
	return bytes.Equal(encP[:], sig[:32])
}

func (e *Ed448) KeyGen(public *Pk448, private *Sk448) {
	var dig [114]byte
	H := sha3.NewShake256()
	_, _ = H.Write(private[:])
	_, _ = H.Read(dig[:])
	dig[0] &= -(uint8(1) << edwards448.lgCofactor)
	// dig[56] = (dig[31] & 127) | 64

	edwards448.reduceModOrder(dig[:57])

	var P point448R1
	var S point448R3
	edwards448.fixedMult(&P, &S, dig[:57])
	P.ToBytes(public[:])
}
func (e *Ed448) Sign(msg []byte, private *Sk448) *Sig448 {
	var signature Sig448
	return &signature
}
func (e *Ed448) Verify(msg []byte, public *Pk448, sig *Sig448) bool {
	return false
}
