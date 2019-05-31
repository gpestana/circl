// +build amd64,!noasm

#include "textflag.h"

// Multipies 512-bit value by 64-bit value. Uses MULQ instruction to
// multiply 2 64-bit values.
// 
// Result: x = (y * z) mod 2^512
//
// Registers used: AX, CX, DX, SI, DI, R8
//
// func mul512(a, b *u512, c uint64)
TEXT ·mul512(SB), NOSPLIT, $0-24
	MOVQ	x+ 0(FP), DI	// result
	MOVQ	y+ 8(FP), SI	// multiplicand

	// Check wether to use optimized implementation
	CMPB    ·HasBMI2(SB), $1
	JE      mul512_mulx

	MOVQ z+16(FP), R10	// 64 bit multiplier, used by MULQ
	MOVQ R10, AX; MULQ  0(SI);				              MOVQ DX, R11; MOVQ AX,  0(DI) //x[0]
	MOVQ R10, AX; MULQ  8(SI); ADDQ R11, AX; ADCQ $0, DX; MOVQ DX, R11; MOVQ AX,  8(DI) //x[1]
	MOVQ R10, AX; MULQ 16(SI); ADDQ R11, AX; ADCQ $0, DX; MOVQ DX, R11; MOVQ AX, 16(DI) //x[2]
	MOVQ R10, AX; MULQ 24(SI); ADDQ R11, AX; ADCQ $0, DX; MOVQ DX, R11; MOVQ AX, 24(DI) //x[3]
	MOVQ R10, AX; MULQ 32(SI); ADDQ R11, AX; ADCQ $0, DX; MOVQ DX, R11; MOVQ AX, 32(DI) //x[4]
	MOVQ R10, AX; MULQ 40(SI); ADDQ R11, AX; ADCQ $0, DX; MOVQ DX, R11; MOVQ AX, 40(DI) //x[5]
	MOVQ R10, AX; MULQ 48(SI); ADDQ R11, AX; ADCQ $0, DX; MOVQ DX, R11; MOVQ AX, 48(DI) //x[6]
	MOVQ R10, AX; MULQ 56(SI); ADDQ R11, AX;	                        MOVQ AX, 56(DI) //x[7]
	RET

// Optimized for CPUs with BMI2
mul512_mulx:
	MOVQ	 z+16(FP), DX	// 64 bit multiplier, used by MULX
	MULXQ	 0(SI), AX, R10; MOVQ AX, 0(DI) // x[0]
	MULXQ	 8(SI), AX, R11; ADDQ R10, AX; MOVQ AX,  8(DI) // x[1]
	MULXQ	16(SI), AX, R10; ADCQ R11, AX; MOVQ AX, 16(DI) // x[2]
	MULXQ	24(SI), AX, R11; ADCQ R10, AX; MOVQ AX, 24(DI) // x[3]
	MULXQ	32(SI), AX, R10; ADCQ R11, AX; MOVQ AX, 32(DI) // x[4]
	MULXQ	40(SI), AX, R11; ADCQ R10, AX; MOVQ AX, 40(DI) // x[5]
	MULXQ	48(SI), AX, R10; ADCQ R11, AX; MOVQ AX, 48(DI) // x[6]
	MULXQ	56(SI), AX, R11; ADCQ R10, AX; MOVQ AX, 56(DI) // x[7]
	RET

// x = y + z
// func add512(x, y, z *u512) uint64
TEXT ·add512(SB), NOSPLIT, $0-32
	MOVQ	x+ 0(FP), DI	// result
	MOVQ	y+ 8(FP), SI	// first summand
	MOVQ	z+16(FP), DX	// second summand

	XORQ	AX, AX

	MOVQ	 0(SI), R8;	ADDQ	 0(DX), R8;	MOVQ	R8,  0(DI)	// x[0]
	MOVQ	 8(SI), R8;	ADCQ	 8(DX), R8;	MOVQ	R8,  8(DI)	// x[1]
	MOVQ	16(SI), R8;	ADCQ	16(DX), R8;	MOVQ	R8, 16(DI)	// x[2]
	MOVQ	24(SI), R8;	ADCQ	24(DX), R8;	MOVQ	R8, 24(DI)	// x[3]
	MOVQ	32(SI), R8;	ADCQ	32(DX), R8;	MOVQ	R8, 32(DI)	// x[4]
	MOVQ	40(SI), R8;	ADCQ	40(DX), R8;	MOVQ	R8, 40(DI)	// x[5]
	MOVQ	48(SI), R8;	ADCQ	48(DX), R8;	MOVQ	R8, 48(DI)	// x[6]
	MOVQ	56(SI), R8;	ADCQ	56(DX), R8;	MOVQ	R8, 56(DI)	// x[7]

	// return carry
	ADCQ	AX, AX
	MOVQ	AX, ret+24(FP)
	RET


// x = y - z
// func sub512(x, y, z *u512) uint64
TEXT ·sub512(SB), NOSPLIT, $0-32
	MOVQ	x+ 0(FP), DI	// result
	MOVQ	y+ 8(FP), SI	// minuend
	MOVQ	z+16(FP), DX	// subtrahend

	XORQ	AX, AX

	MOVQ	 0(SI), R8;	SUBQ	 0(DX), R8;	MOVQ	R8,  0(DI)	// x[0]
	MOVQ	 8(SI), R8;	SBBQ	 8(DX), R8;	MOVQ	R8,  8(DI)	// x[1]
	MOVQ	16(SI), R8;	SBBQ	16(DX), R8;	MOVQ	R8, 16(DI)	// x[2]
	MOVQ	24(SI), R8;	SBBQ	24(DX), R8;	MOVQ	R8, 24(DI)	// x[3]
	MOVQ	32(SI), R8;	SBBQ	32(DX), R8;	MOVQ	R8, 32(DI)	// x[4]
	MOVQ	40(SI), R8;	SBBQ	40(DX), R8;	MOVQ	R8, 40(DI)	// x[5]
	MOVQ	48(SI), R8;	SBBQ	48(DX), R8;	MOVQ	R8, 48(DI)	// x[6]
	MOVQ	56(SI), R8;	SBBQ	56(DX), R8;	MOVQ	R8, 56(DI)	// x[7]

	// return borrow
	ADCQ	AX, AX
	MOVQ	AX, ret+24(FP)

	RET

TEXT ·cswap512(SB),NOSPLIT,$0-17
	MOVQ    x+0(FP), DI
	MOVQ    y+8(FP), SI
    MOVBLZX choice+16(FP), AX       // AL = 0 or 1

	// Make AX, so that either all bits are set or non
	// AX = 0 or 1
	NEGQ    AX

	// Fill xmm15. After this step first half of XMM15 is
	// just zeros and second half is whatever in AX
	MOVQ    AX, X15

	// Copy lower double word everywhere else. So that
	// XMM15=AL|AL|AL|AL. As AX has either all bits set
	// or non result will be that XMM15 has also either
	// all bits set or non of them.
	PSHUFD $0, X15, X15

#ifndef CSWAP_BLOCK
#define CSWAP_BLOCK(idx)       \
	MOVOU   (idx*16)(DI), X0 \
	MOVOU   (idx*16)(SI), X1 \
	\ // X2 = mask & (X0 ^ X1)
	MOVO     X1, X2 \
	PXOR     X0, X2 \
	PAND    X15, X2 \
	\
	PXOR     X2, X0 \
	PXOR     X2, X1 \
	\
	MOVOU    X0, (idx*16)(DI) \
	MOVOU    X1, (idx*16)(SI)
#endif

	CSWAP_BLOCK(0)
	CSWAP_BLOCK(1)
	CSWAP_BLOCK(2)
	CSWAP_BLOCK(3)

	RET

// val = val<p?val:val-p
TEXT ·crdc512(SB),NOSPLIT,$0-8
	MOVQ val+0(FP), DI

	MOVQ ( 0)(DI),  SI; SUBQ ·p+ 0(SB),  SI
	MOVQ ( 8)(DI),  DX; SBBQ ·p+ 8(SB),  DX
	MOVQ (16)(DI),  CX; SBBQ ·p+16(SB),  CX
	MOVQ (24)(DI),  R8; SBBQ ·p+24(SB),  R8
	MOVQ (32)(DI),  R9; SBBQ ·p+32(SB),  R9
	MOVQ (40)(DI), R10; SBBQ ·p+40(SB), R10
	MOVQ (48)(DI), R11; SBBQ ·p+48(SB), R11
	MOVQ (56)(DI), R12; SBBQ ·p+56(SB), R12

	MOVQ ( 0)(DI), AX; CMOVQCC  SI, AX; MOVQ AX, ( 0)(DI)
	MOVQ ( 8)(DI), AX; CMOVQCC  DX, AX; MOVQ AX, ( 8)(DI)
	MOVQ (16)(DI), AX; CMOVQCC  CX, AX; MOVQ AX, (16)(DI)
	MOVQ (24)(DI), AX; CMOVQCC  R8, AX; MOVQ AX, (24)(DI)
	MOVQ (32)(DI), AX; CMOVQCC  R9, AX; MOVQ AX, (32)(DI)
	MOVQ (40)(DI), AX; CMOVQCC R10, AX; MOVQ AX, (40)(DI)
	MOVQ (48)(DI), AX; CMOVQCC R11, AX; MOVQ AX, (48)(DI)
	MOVQ (56)(DI), AX; CMOVQCC R12, AX; MOVQ AX, (56)(DI)

	RET

// val = b?val+p:val
TEXT ·csubrdc512(SB),NOSPLIT,$0-16
	MOVQ val+0(FP), DI
	MOVQ choice+8(FP), SI

	XORQ  R8,  R8
	XORQ  R9,  R9
	XORQ R10, R10
	XORQ R11, R11
	XORQ R12, R12
	XORQ R13, R13
	XORQ R14, R14
	XORQ R15, R15

	TESTQ SI, SI
	CMOVQNE ·p+ 0(SB), R8
	CMOVQNE ·p+ 8(SB), R9
	CMOVQNE ·p+16(SB), R10
	CMOVQNE ·p+24(SB), R11
	CMOVQNE ·p+32(SB), R12
	CMOVQNE ·p+40(SB), R13
	CMOVQNE ·p+48(SB), R14
	CMOVQNE ·p+56(SB), R15

	MOVQ ( 0)(DI), DX; ADDQ  R8, DX; MOVQ DX, ( 0)(DI)
	MOVQ ( 8)(DI), DX; ADCQ  R9, DX; MOVQ DX, ( 8)(DI)
	MOVQ (16)(DI), DX; ADCQ R10, DX; MOVQ DX, (16)(DI)
	MOVQ (24)(DI), DX; ADCQ R11, DX; MOVQ DX, (24)(DI)
	MOVQ (32)(DI), DX; ADCQ R12, DX; MOVQ DX, (32)(DI)
	MOVQ (40)(DI), DX; ADCQ R13, DX; MOVQ DX, (40)(DI)
	MOVQ (48)(DI), DX; ADCQ R14, DX; MOVQ DX, (48)(DI)
	MOVQ (56)(DI), DX; ADCQ R15, DX; MOVQ DX, (56)(DI)

	RET

// mul function implements montgomery multiplication interleaved with rdc.
// It takes advantage of the fact that inversion of 'p' has only 64-bits
//
// z = x*y mod p
TEXT ·mul(SB),NOSPLIT,$32-24
	// Check wether to use optimized implementation
	CMPB    ·HasADXandBMI2(SB), $1
	JE      mul_with_mulx_adcx_adox	


	// Generic x86 implementation (below) uses variant of Karatsuba method.
	//
	// Here we store the destination in CX instead of in REG_P3 because the
	// multiplication instructions use DX as an implicit destination
	// operand: MULQ $REG sets DX:AX <-- AX * $REG.

	// RAX and RDX will be used for a mask (0-borrow)
	XORQ	AX, AX

	// RCX[0-3]: U1+U0
	MOVQ	(32)(REG_P1), R8
	MOVQ	(40)(REG_P1), R9
	MOVQ	(48)(REG_P1), R10
	MOVQ	(56)(REG_P1), R11
	ADDQ	( 0)(REG_P1), R8
	ADCQ	( 8)(REG_P1), R9
	ADCQ	(16)(REG_P1), R10
	ADCQ	(24)(REG_P1), R11
	MOVQ	R8,  ( 0)(CX)
	MOVQ	R9,  ( 8)(CX)
	MOVQ	R10, (16)(CX)
	MOVQ	R11, (24)(CX)

	SBBQ	$0, AX

	// R12-R15: V1+V0
	XORQ	DX, DX
	MOVQ	(32)(REG_P2), R12
	MOVQ	(40)(REG_P2), R13
	MOVQ	(48)(REG_P2), R14
	MOVQ	(56)(REG_P2), R15
	ADDQ	( 0)(REG_P2), R12
	ADCQ	( 8)(REG_P2), R13
	ADCQ	(16)(REG_P2), R14
	ADCQ	(24)(REG_P2), R15

	SBBQ	$0, DX

	// Store carries on stack
	MOVQ	AX, (64)(SP)
	MOVQ	DX, (72)(SP)

	// (SP[0-3],R8,R9,R10,R11) <- (U0+U1)*(V0+V1).
	// MUL using comba; In comments below U=U0+U1 V=V0+V1

	// U0*V0
	MOVQ    (CX), AX
	MULQ    R12
	MOVQ    AX, (SP)        // C0
	MOVQ    DX, R8

	// U0*V1
	XORQ    R9, R9
	MOVQ    (CX), AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9

	// U1*V0
	XORQ    R10, R10
	MOVQ    (8)(CX), AX
	MULQ    R12
	ADDQ    AX, R8
	MOVQ    R8, (8)(SP)     // C1
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U0*V2
	XORQ    R8, R8
	MOVQ    (CX), AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U2*V0
	MOVQ    (16)(CX), AX
	MULQ    R12
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U1*V1
	MOVQ    (8)(CX), AX
	MULQ    R13
	ADDQ    AX, R9
	MOVQ    R9, (16)(SP)        // C2
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U0*V3
	XORQ    R9, R9
	MOVQ    (CX), AX
	MULQ    R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U3*V0
	MOVQ    (24)(CX), AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U1*V2
	MOVQ    (8)(CX), AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U2*V1
	MOVQ    (16)(CX), AX
	MULQ    R13
	ADDQ    AX, R10
	MOVQ    R10, (24)(SP)       // C3
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U1*V3
	XORQ    R10, R10
	MOVQ    (8)(CX), AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U3*V1
	MOVQ    (24)(CX), AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V2
	MOVQ    (16)(CX), AX
	MULQ    R14
	ADDQ    AX, R8
	MOVQ    R8, (32)(SP)        // C4
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V3
	XORQ    R11, R11
	MOVQ    (16)(CX), AX
	MULQ    R15
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R11

	// U3*V2
	MOVQ    (24)(CX), AX
	MULQ    R14
	ADDQ    AX, R9              // C5
	ADCQ    DX, R10
	ADCQ    $0, R11

	// U3*V3
	MOVQ    (24)(CX), AX
	MULQ    R15
	ADDQ    AX, R10             // C6
	ADCQ    DX, R11             // C7

	MOVQ    (64)(SP), AX
	ANDQ    AX, R12
	ANDQ    AX, R13
	ANDQ    AX, R14
	ANDQ    AX, R15
	ADDQ    R8, R12
	ADCQ    R9, R13
	ADCQ    R10, R14
	ADCQ    R11, R15

	MOVQ    (72)(SP), AX
	MOVQ    (CX), R8
	MOVQ    (8)(CX), R9
	MOVQ    (16)(CX), R10
	MOVQ    (24)(CX), R11
	ANDQ    AX, R8
	ANDQ    AX, R9
	ANDQ    AX, R10
	ANDQ    AX, R11
	ADDQ    R12, R8
	ADCQ    R13, R9
	ADCQ    R14, R10
	ADCQ    R15, R11
	MOVQ    R8, (32)(SP)
	MOVQ    R9, (40)(SP)
	MOVQ    R10, (48)(SP)
	MOVQ    R11, (56)(SP)

	// CX[0-7] <- AL*BL

	// U0*V0
	MOVQ    (REG_P1), R11
	MOVQ    (REG_P2), AX
	MULQ    R11
	XORQ    R9, R9
	MOVQ    AX, (CX)            // C0
	MOVQ    DX, R8

	// U0*V1
	MOVQ    (16)(REG_P1), R14
	MOVQ    (8)(REG_P2), AX
	MULQ    R11
	XORQ    R10, R10
	ADDQ    AX, R8
	ADCQ    DX, R9

	// U1*V0
	MOVQ    (8)(REG_P1), R12
	MOVQ    (REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	MOVQ    R8, (8)(CX)         // C1
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U0*V2
	XORQ    R8, R8
	MOVQ    (16)(REG_P2), AX
	MULQ    R11
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U2*V0
	MOVQ    (REG_P2), R13
	MOVQ    R14, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U1*V1
	MOVQ    (8)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R9
	MOVQ    R9, (16)(CX)        // C2
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U0*V3
	XORQ    R9, R9
	MOVQ    (24)(REG_P2), AX
	MULQ    R11
	MOVQ    (24)(REG_P1), R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U3*V1
	MOVQ    R15, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U2*V2
	MOVQ    (16)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U2*V3
	MOVQ    (8)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R10
	MOVQ    R10, (24)(CX)       // C3
	ADCQ    DX, R8
	ADCQ    $0, R9

	// U3*V2
	XORQ    R10, R10
	MOVQ    (24)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U3*V1
	MOVQ    (8)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V2
	MOVQ    (16)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R8
	MOVQ    R8, (32)(CX)		// C4
	ADCQ    DX, R9
	ADCQ    $0, R10

	// U2*V3
	XORQ    R8, R8
	MOVQ    (24)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U3*V2
	MOVQ    (16)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R9
	MOVQ    R9, (40)(CX)		// C5
	ADCQ    DX, R10
	ADCQ    $0, R8

	// U3*V3
	MOVQ    (24)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R10
	MOVQ    R10, (48)(CX)		// C6
	ADCQ    DX, R8
	MOVQ    R8, (56)(CX)		// C7

	// CX[8-15] <- U1*V1
	MOVQ    (32)(REG_P1), R11
	MOVQ    (32)(REG_P2), AX
	MULQ    R11
	XORQ    R9, R9
	MOVQ    AX, (64)(CX)        // C0
	MOVQ    DX, R8

	MOVQ    (48)(REG_P1), R14
	MOVQ    (40)(REG_P2), AX
	MULQ    R11
	XORQ    R10, R10
	ADDQ    AX, R8
	ADCQ    DX, R9

	MOVQ    (40)(REG_P1), R12
	MOVQ    (32)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	MOVQ    R8, (72)(CX)        // C1
	ADCQ    DX, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    (48)(REG_P2), AX
	MULQ    R11
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (32)(REG_P2), R13
	MOVQ    R14, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (40)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R9
	MOVQ    R9, (80)(CX)        // C2
	ADCQ    DX, R10
	ADCQ    $0, R8

	XORQ    R9, R9
	MOVQ    (56)(REG_P2), AX
	MULQ    R11
	MOVQ    (56)(REG_P1), R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    R15, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    (48)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9

	MOVQ    (40)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R10
	MOVQ    R10, (88)(CX)       // C3
	ADCQ    DX, R8
	ADCQ    $0, R9

	XORQ    R10, R10
	MOVQ    (56)(REG_P2), AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    (40)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10

	MOVQ    (48)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R8
	MOVQ    R8, (96)(CX)        // C4
	ADCQ    DX, R9
	ADCQ    $0, R10

	XORQ    R8, R8
	MOVQ    (56)(REG_P2), AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (48)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R9
	MOVQ    R9, (104)(CX)       // C5
	ADCQ    DX, R10
	ADCQ    $0, R8

	MOVQ    (56)(REG_P2), AX
	MULQ    R15
	ADDQ    AX, R10
	MOVQ    R10, (112)(CX)      // C6
	ADCQ    DX, R8
	MOVQ    R8, (120)(CX)       // C7

	// [R8-R15] <- (U0+U1)*(V0+V1) - U1*V1
	MOVQ    (SP), R8
	SUBQ    (CX), R8
	MOVQ    (8)(SP), R9
	SBBQ    (8)(CX), R9
	MOVQ    (16)(SP), R10
	SBBQ    (16)(CX), R10
	MOVQ    (24)(SP), R11
	SBBQ    (24)(CX), R11
	MOVQ    (32)(SP), R12
	SBBQ    (32)(CX), R12
	MOVQ    (40)(SP), R13
	SBBQ    (40)(CX), R13
	MOVQ    (48)(SP), R14
	SBBQ    (48)(CX), R14
	MOVQ    (56)(SP), R15
	SBBQ    (56)(CX), R15

	// [R8-R15] <- (U0+U1)*(V0+V1) - U1*V0 - U0*U1
	MOVQ    ( 64)(CX), AX;	SUBQ    AX, R8
	MOVQ    ( 72)(CX), AX;	SBBQ    AX, R9
	MOVQ    ( 80)(CX), AX;	SBBQ    AX, R10
	MOVQ    ( 88)(CX), AX;	SBBQ    AX, R11
	MOVQ    ( 96)(CX), AX;	SBBQ    AX, R12
	MOVQ    (104)(CX), DX;	SBBQ    DX, R13
	MOVQ    (112)(CX), DI;	SBBQ    DI, R14
	MOVQ    (120)(CX), SI;	SBBQ    SI, R15

	// Final result
	ADDQ    (32)(CX), R8;	MOVQ    R8,  (32)(CX)
	ADCQ    (40)(CX), R9;	MOVQ    R9,  (40)(CX)
	ADCQ    (48)(CX), R10;	MOVQ    R10, (48)(CX)
	ADCQ    (56)(CX), R11;	MOVQ    R11, (56)(CX)
	ADCQ    (64)(CX), R12;	MOVQ    R12, (64)(CX)
	ADCQ    (72)(CX), R13;	MOVQ    R13, (72)(CX)
	ADCQ    (80)(CX), R14;	MOVQ    R14, (80)(CX)
	ADCQ    (88)(CX), R15;	MOVQ    R15, (88)(CX)
	ADCQ    $0, AX;        	MOVQ    AX,  (96)(CX)
	ADCQ    $0, DX;        	MOVQ    DX, (104)(CX)
	ADCQ    $0, DI;         MOVQ    DI, (112)(CX)
	ADCQ    $0, SI;     	MOVQ    SI, (120)(CX)
	RET

// Optimized for CPUs with BMI2 and ADCX/ADOX instructions
mul_with_mulx_adcx_adox:
	MOVQ y+ 8(FP), DI // multiplicand
	MOVQ z+16(FP), SI // multiplier

	XORQ  R8,  R8
	XORQ  R9,  R9
	XORQ R10, R10
	XORQ R11, R11
	XORQ R12, R12
	XORQ R13, R13
	XORQ R14, R14
	XORQ R15, R15

	MOVQ BP, 24(SP) // OZAPTF: thats maybe wrong
	XORQ BP, BP

// Uses BMI2 (MULX)
#ifdef MULS_MULX_512
#undef MULS_MULX_512
#endif
#define MULS_MULX_512(idx, r0, r1, r2, r3, r4, r5, r6, r7, r8) \
	\ // Reduction step
	MOVQ  ( 0)(SI), DX 		\
	MULXQ ( 8*idx)(DI), DX, CX 	\
	ADDQ  r0, DX 			\
	MULXQ ·pNegInv(SB), DX, CX	\
	\
	XORQ  AX, AX \
	MULXQ ·p+ 0(SB), AX, BX;             ; ADOXQ AX, r0 \
	MULXQ ·p+ 8(SB), AX, CX; ADCXQ BX, r1; ADOXQ AX, r1 \
	MULXQ ·p+16(SB), AX, BX; ADCXQ CX, r2; ADOXQ AX, r2 \
	MULXQ ·p+24(SB), AX, CX; ADCXQ BX, r3; ADOXQ AX, r3 \
	MULXQ ·p+32(SB), AX, BX; ADCXQ CX, r4; ADOXQ AX, r4 \
	MULXQ ·p+40(SB), AX, CX; ADCXQ BX, r5; ADOXQ AX, r5 \
	MULXQ ·p+48(SB), AX, BX; ADCXQ CX, r6; ADOXQ AX, r6 \
	MULXQ ·p+56(SB), AX, CX; ADCXQ BX, r7; ADOXQ AX, r7 \
	MOVQ  $0, AX           ; ADCXQ CX, r8; ADOXQ AX, r8 \
	\ // Multiplication step
	MOVQ (8*idx)(DI), DX \
	\
	XORQ  AX, AX \
	MULXQ ( 0)(SI), AX, BX; ADOXQ AX, r0 \
	MULXQ ( 8)(SI), AX, CX; ADCXQ BX, r1; ADOXQ AX, r1 \
	MULXQ (16)(SI), AX, BX; ADCXQ CX, r2; ADOXQ AX, r2 \
	MULXQ (24)(SI), AX, CX; ADCXQ BX, r3; ADOXQ AX, r3 \
	MULXQ (32)(SI), AX, BX; ADCXQ CX, r4; ADOXQ AX, r4 \
	MULXQ (40)(SI), AX, CX; ADCXQ BX, r5; ADOXQ AX, r5 \
	MULXQ (48)(SI), AX, BX; ADCXQ CX, r6; ADOXQ AX, r6 \
	MULXQ (56)(SI), AX, CX; ADCXQ BX, r7; ADOXQ AX, r7 \
	MOVQ  $0, AX          ; ADCXQ CX, r8; ADOXQ AX, r8

	MULS_MULX_512(0,  R8,  R9, R10, R11, R12, R13, R14, R15,  BP)
	MULS_MULX_512(1,  R9, R10, R11, R12, R13, R14, R15,  BP,  R8)
	MULS_MULX_512(2, R10, R11, R12, R13, R14, R15,  BP,  R8,  R9)
	MULS_MULX_512(3, R11, R12, R13, R14, R15,  BP,  R8,  R9, R10)
	MULS_MULX_512(4, R12, R13, R14, R15,  BP,  R8,  R9, R10, R11)
	MULS_MULX_512(5, R13, R14, R15,  BP,  R8,  R9, R10, R11, R12)
	MULS_MULX_512(6, R14, R15,  BP,  R8,  R9, R10, R11, R12, R13)
	MULS_MULX_512(7, R15,  BP,  R8,  R9, R10, R11, R12, R13, R14)
#undef MULS_MULX_512

	MOVQ x+0(FP), DI
	MOVQ  BP, ( 0)(DI)
	MOVQ  R8, ( 8)(DI)
	MOVQ  R9, (16)(DI)
	MOVQ R10, (24)(DI)
	MOVQ R11, (32)(DI)
	MOVQ R12, (40)(DI)
	MOVQ R13, (48)(DI)
	MOVQ R14, (56)(DI)
	MOVQ 24(SP), BP

	// NOW DI needs to be reduced if > p
	RET

// Checks if x>y. Returns 1 if true otherwise, 0
TEXT ·checkBigger(SB),NOSPLIT,$0-24
	MOVQ	x+ 0(FP), DI	// minuend
	MOVQ	y+ 8(FP), SI	// subtrahend

	XORQ	AX, AX
	MOVQ	 0(SI), R8;	SUBQ	 0(DI), R8
	MOVQ	 8(SI), R8;	SBBQ	 8(DI), R8
	MOVQ	16(SI), R8;	SBBQ	16(DI), R8
	MOVQ	24(SI), R8;	SBBQ	24(DI), R8
	MOVQ	32(SI), R8;	SBBQ	32(DI), R8
	MOVQ	40(SI), R8;	SBBQ	40(DI), R8
	MOVQ	48(SI), R8;	SBBQ	48(DI), R8
	MOVQ	56(SI), R8;	SBBQ	56(DI), R8

	// return borrow
	ADCQ	AX, AX
	MOVQ	AX, 24(SP)

	RET
