// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

// +build amd64,!noasm arm64,!noasm

package p434

import (
	. "github.com/cloudflare/circl/dh/sidh/internal/common"
)

// If choice = 0, leave x,y unchanged. If choice = 1, set x,y = y,x.
// If choice is neither 0 nor 1 then behaviour is undefined.
// This function executes in constant time.
//go:noescape
func cswapP434(x, y *Fp, choice uint8)

// Compute z = x + y (mod p).
//go:noescape
func addP434(z, x, y *Fp)

// Compute z = x - y (mod p).
//go:noescape
func subP434(z, x, y *Fp)

// Compute z = x + y, without reducing mod p.
//go:noescape
func adlP434(z, x, y *FpX2)

// Compute z = x - y, without reducing mod p.
//go:noescape
func sulP434(z, x, y *FpX2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func modP434(x *Fp)

// Computes z = x * y.
//go:noescape
func mulP434(z *FpX2, x, y *Fp)

// Computes the Montgomery reduction z = x R^{-1} (mod 2*p). On return value
// of x may be changed. z=x not allowed.
//go:noescape
func rdcP434(z *Fp, x *FpX2)
