package eddsa

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
	var P point255R1
	var S point255R3
	edwards448.fixedMult(&P, &S, private[:])
	P.ToBytes(public)
}
func (e *Ed25519) Sign(msg []byte, private *Sk25519) {

}
func (e *Ed25519) Verify(msg []byte, public *Pk25519, sig *Sig25519) bool {
	return false
}

func (e *Ed448) KeyGen(public *Pk448, private *Sk448) {

}
func (e *Ed448) Sign(msg []byte, private *Sk448) {

}
func (e *Ed448) Verify(msg []byte, public *Pk448, sig *Sig448) bool {
	return false
}
