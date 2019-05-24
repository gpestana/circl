// +build amd64

package x448

import (
	fp "github.com/cloudflare/circl/math/fp448"
	"github.com/cloudflare/circl/utils/cpu"
)

var hasBmi2Adx = cpu.X86.HasBMI2 && cpu.X86.HasADX

//go:noescape
func ladderStep(work *[5]fp.Elt, move uint)

// diffAdd calculates a differential point addition using a precomputed point.
// (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
// work = {x1,z1,x2,z2} are four fp.Elt of fp.Size bytes.
// See Equation 7 at https://eprint.iacr.org/2017/264.
//go:noescape
func difAdd(work *[4]fp.Elt, mu *fp.Elt, swap uint)

// double calculates a point doubling (x1,z1) = 2*(x1,z1).
//   work  = {x1,z1,x2,z2} are four fp.Elt of fp.Size bytes, and
// Variables x2,z2 are modified.
//go:noescape
func double(work *[4]fp.Elt)

//go:noescape
func mulA24(z, x *fp.Elt)
