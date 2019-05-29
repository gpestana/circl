// +build amd64

package xcurve

import (
	"github.com/cloudflare/circl/utils/cpu"
)

var hasBmi2Adx = cpu.X86.HasBMI2 && cpu.X86.HasADX

// ladderStep255 calculates a point addition and doubling as follows:
// (x2,z2) = 2*(x2,z2) and (x3,z3) = (x2,z2)+(x3,z3) using as a difference (x1,-).
//   work  = {x1,x2,z2,x3,z4} are five fp255.Elt of 32 bytes.
//go:noescape
func ladderStep255(work []byte, move uint)

// diffAdd255 calculates a differential point addition using a precomputed point.
// (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
//    work = {x1,z1,x2,z2} are four fp.Elt of fp.Size bytes.
//      mu = {mu} is one element fp255.Elt of 32 bytes.
// See Equation 7 at https://eprint.iacr.org/2017/264.
//go:noescape
func difAdd255(work []byte, mu []byte, swap uint)

// double calculates a point doubling (x1,z1) = 2*(x1,z1).
//   work  = {x1,z1,x2,z2} are four fp255.Elt of 32 bytes each one.
// Variables x2,z2 are modified.
//go:noescape
func double255(work []byte)

// ladderStep448 calculates a point addition and doubling as follows:
// (x2,z2) = 2*(x2,z2) and (x3,z3) = (x2,z2)+(x3,z3) using as a difference (x1,-).
//   work  = {x1,x2,z2,x3,z4} are five fp448.Elt of 56 bytes.
//go:noescape
func ladderStep448(work []byte, move uint)

// diffAdd448 calculates a differential point addition using a precomputed point.
// (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
// work = {x1,z1,x2,z2} are four fp448.Elt of 56 bytes.
//   mu = {mu} is one element fp448.Elt of 56 bytes.
// See Equation 7 at https://eprint.iacr.org/2017/264.
//go:noescape
func difAdd448(work []byte, mu []byte, swap uint)

// double calculates a point doubling (x1,z1) = 2*(x1,z1).
//   work = {x1,z1,x2,z2} are four fp448.Elt of 56 bytes each one.
// Variables x2,z2 are modified.
//go:noescape
func double448(work []byte)
