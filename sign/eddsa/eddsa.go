package eddsa

import (
	"crypto/sha512"
	"golang.org/x/crypto/sha3"
)

const (
	SizeEd25519 = 32
	SizeEd448   = 57
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
	dig := H.Sum([]byte{})
	dig[0] &= -(uint8(1) << edwards25519.lgCofactor)
	dig[31] = (dig[31] & 127) | 64

	edwards25519.reduceModOrder(dig[:32])

	var P point255R1
	var S point255R3
	edwards25519.fixedMult(&P, &S, dig[:32])
	P.ToBytes(public[:])
}
func (e *Ed25519) Sign(msg []byte, private *Sk25519) *Sig25519 {
	return nil
}
func (e *Ed25519) Verify(msg []byte, public *Pk25519, sig *Sig25519) bool {
	return false
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
func (e *Ed448) Sign(msg []byte, private *Sk448) {

}
func (e *Ed448) Verify(msg []byte, public *Pk448, sig *Sig448) bool {
	return false
}
