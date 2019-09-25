// +build amd64,!noasm

#include "textflag.h"

// p434 + 1
#define P434P1_3 $0xFDC1767AE3000000
#define P434P1_4 $0x7BC65C783158AEA3
#define P434P1_5 $0x6CFC5FD681C52056
#define P434P1_6 $0x0002341F27177344

// p434 x 2
#define P434X2_0 $0xFFFFFFFFFFFFFFFE
#define P434X2_1 $0xFFFFFFFFFFFFFFFF
#define P434X2_3 $0xFB82ECF5C5FFFFFF
#define P434X2_4 $0xF78CB8F062B15D47
#define P434X2_5 $0xD9F8BFAD038A40AC
#define P434X2_6 $0x0004683E4E2EE688

TEXT ·cswapP434(SB),NOSPLIT,$0-17

	MOVQ	x+0(FP), DI
	MOVQ	y+8(FP), SI
	MOVB	choice+16(FP), AL	// AL = 0 or 1
	MOVBLZX	AL, AX				// AX = 0 or 1
	NEGQ	AX					// AX = 0x00..00 or 0xff..ff

#ifndef CSWAP_BLOCK
#define CSWAP_BLOCK(idx) 	\
	MOVQ	(idx*8)(DI), BX	\ // BX = x[idx]
	MOVQ 	(idx*8)(SI), CX	\ // CX = y[idx]
	MOVQ	CX, DX			\ // DX = y[idx]
	XORQ	BX, DX			\ // DX = y[idx] ^ x[idx]
	ANDQ	AX, DX			\ // DX = (y[idx] ^ x[idx]) & mask
	XORQ	DX, BX			\ // BX = (y[idx] ^ x[idx]) & mask) ^ x[idx] = x[idx] or y[idx]
	XORQ	DX, CX			\ // CX = (y[idx] ^ x[idx]) & mask) ^ y[idx] = y[idx] or x[idx]
	MOVQ	BX, (idx*8)(DI)	\
	MOVQ	CX, (idx*8)(SI)
#endif

	CSWAP_BLOCK(0)
	CSWAP_BLOCK(1)
	CSWAP_BLOCK(2)
	CSWAP_BLOCK(3)
	CSWAP_BLOCK(4)
	CSWAP_BLOCK(5)
	CSWAP_BLOCK(6)

#ifdef CSWAP_BLOCK
#undef CSWAP_BLOCK
#endif

	RET

TEXT ·addP434(SB),NOSPLIT,$0-24
	MOVQ	z+0(FP), DX
	MOVQ	x+8(FP), DI
	MOVQ	y+16(FP), SI

	// Used later to calculate a mask
	XORQ    CX, CX

	// [R8-R14]: z = x + y
	MOVQ	( 0)(DI), R8;	ADDQ	( 0)(SI), R8
	MOVQ	( 8)(DI), R9;	ADCQ	( 8)(SI), R9
	MOVQ	(16)(DI), R10;	ADCQ	(16)(SI), R10
	MOVQ	(24)(DI), R11;	ADCQ	(24)(SI), R11
	MOVQ	(32)(DI), R12;	ADCQ	(32)(SI), R12
	MOVQ	(40)(DI), R13;	ADCQ	(40)(SI), R13
	MOVQ	(48)(DI), R14;	ADCQ	(48)(SI), R14

	XORQ	DI, DI

	MOVQ	P434X2_0, AX;	SUBQ	AX, R8
	MOVQ	P434X2_1, AX;	SBBQ	AX, R9
					SBBQ	AX, R10
	MOVQ	P434X2_3, AX;	SBBQ	AX, R11
	MOVQ	P434X2_4, AX;	SBBQ	AX, R12
	MOVQ	P434X2_5, AX;	SBBQ	AX, R13
	MOVQ	P434X2_6, AX;	SBBQ	AX, R14

	// mask
	SBBQ	$0, CX

	// if z<0 add P434x2 back
	MOVQ	P434X2_0, R15;	ANDQ	CX, R15;
	MOVQ	P434X2_1, AX;	ANDQ	CX, AX;
	MOVQ	P434X2_1, BX;	ANDQ	CX, BX; // not needed OZAPTF

	ADDQ	R8, R15;	MOVQ R15, ( 0)(DX)
	ADCQ	R9, AX;	MOVQ  AX, ( 8)(DX)
	ADCQ	R10,BX;	MOVQ  BX, (16)(DX)

	ADCQ 	$0, DI
	MOVQ	P434X2_3, R15;	ANDQ	CX, R15;
	MOVQ	P434X2_4,  R8;	ANDQ	CX, R8;
	MOVQ	P434X2_5,  R9;	ANDQ	CX, R9;
	MOVQ	P434X2_6, R10;	ANDQ	CX, R10;
	BTQ	$0, DI

	ADCQ	R11, R15;	MOVQ R15, (24)(DX)
	ADCQ	R12, R8;	MOVQ R8,  (32)(DX)
	ADCQ	R13, R9;	MOVQ R9,  (40)(DX)
	ADCQ	R14, R10;	MOVQ R10, (48)(DX)

	RET

TEXT ·adlP434(SB),NOSPLIT,$0-24
	MOVQ	z+0(FP), DX
	MOVQ	x+8(FP), DI
	MOVQ	y+16(FP),SI

	MOVQ	( 0)(DI), R8
	ADDQ	( 0)(SI), R8
	MOVQ	( 8)(DI), R9
	ADCQ	( 8)(SI), R9
	MOVQ	(16)(DI), R10
	ADCQ	(16)(SI), R10
	MOVQ	(24)(DI), R11
	ADCQ	(24)(SI), R11
	MOVQ	(32)(DI), R12
	ADCQ	(32)(SI), R12
	MOVQ	(40)(DI), R13
	ADCQ	(40)(SI), R13
	MOVQ	(48)(DI), R14
	ADCQ	(48)(SI), R14
	MOVQ	(56)(DI), R15
	ADCQ	(56)(SI), R15
	MOVQ	(64)(DI), AX
	ADCQ	(64)(SI), AX
	MOVQ	(72)(DI), BX
	ADCQ	(72)(SI), BX
	MOVQ	(80)(DI), CX
	ADCQ	(80)(SI), CX

	MOVQ	R8, ( 0)(DX)
	MOVQ	R9, ( 8)(DX)
	MOVQ	R10,(16)(DX)
	MOVQ	R11,(24)(DX)
	MOVQ	R12,(32)(DX)
	MOVQ	R13,(40)(DX)
	MOVQ	R14,(48)(DX)
	MOVQ	R15,(56)(DX)
	MOVQ	AX, (64)(DX)
	MOVQ	BX, (72)(DX)
	MOVQ	CX, (80)(DX)

	MOVQ	(88)(DI), R8
	ADCQ	(88)(SI), R8
	MOVQ	(96)(DI), R9
	ADCQ	(96)(SI), R9
	MOVQ	(104)(DI), R10
	ADCQ	(104)(SI), R10

	MOVQ	R8, (88)(DX)
	MOVQ	R9, (96)(DX)
	MOVQ	R10,(104)(DX)

	RET
TEXT ·subP434(SB),NOSPLIT,$0-24
	RET

TEXT ·sulP434(SB),NOSPLIT,$0-24

	MOVQ z+0(FP), DX
	MOVQ x+8(FP), DI
	MOVQ y+16(FP), SI

	// Used later to store result of 0-borrow
	XORQ CX, CX

	// SUBC for first 11 limbs
	MOVQ	( 0)(DI), R8
	MOVQ	( 8)(DI), R9
	MOVQ	(16)(DI), R10
	MOVQ	(24)(DI), R11
	MOVQ	(32)(DI), R12
	MOVQ	(40)(DI), R13
	MOVQ	(48)(DI), R14
	MOVQ	(56)(DI), R15
	MOVQ	(64)(DI), AX
	MOVQ	(72)(DI), BX

	SUBQ	( 0)(SI), R8
	SBBQ	( 8)(SI), R9
	SBBQ	(16)(SI), R10
	SBBQ	(24)(SI), R11
	SBBQ	(32)(SI), R12
	SBBQ	(40)(SI), R13
	SBBQ	(48)(SI), R14
	SBBQ	(56)(SI), R15
	SBBQ	(64)(SI), AX
	SBBQ	(72)(SI), BX

	MOVQ	 R8, ( 0)(DX)
	MOVQ	 R9, ( 8)(DX)
	MOVQ	R10, (16)(DX)
	MOVQ	R11, (24)(DX)
	MOVQ	R12, (32)(DX)
	MOVQ	R13, (40)(DX)
	MOVQ	R14, (48)(DX)
	MOVQ	R15, (56)(DX)
	MOVQ	 AX, (64)(DX)
	MOVQ	 BX, (72)(DX)

	// SUBC for last 5 limbs
	MOVQ	( 80)(DI), 	R8
	MOVQ	( 88)(DI), 	R9
	MOVQ	( 96)(DI), 	R10
	MOVQ	(104)(DI), 	R11

	SBBQ	( 80)(SI), R8
	SBBQ	( 88)(SI), R9
	SBBQ	( 96)(SI), R10
	SBBQ	(104)(SI), R11

	MOVQ	R8,  ( 80)(DX)
	MOVQ	R9,  ( 88)(DX)
	MOVQ	R10, ( 96)(DX)
	MOVQ	R11, (104)(DX)



	// Now the carry flag is 1 if x-y < 0.  If so, add p*2^512.
	SBBQ	$0, CX

	// Load p into registers:
	MOVQ	P434_0, R8
	// P434_{1,2} = P434_0, so reuse R8
	MOVQ	P434_3, R9
	MOVQ	P434_4, R10
	MOVQ	P434_5, R11
	MOVQ	P434_6, R12
	MOVQ	P434_7, R13

	ANDQ	CX, R8
	ANDQ	CX, R9
	ANDQ	CX, R10
	ANDQ	CX, R11
	ANDQ	CX, R12
	ANDQ	CX, R13

	MOVQ   (64   )(DX), AX; ADDQ R8,  AX; MOVQ AX, (64   )(DX)
	MOVQ   (64+ 8)(DX), AX; ADCQ R8,  AX; MOVQ AX, (64+ 8)(DX)
	MOVQ   (64+16)(DX), AX; ADCQ R8,  AX; MOVQ AX, (64+16)(DX)
	MOVQ   (64+24)(DX), AX; ADCQ R9,  AX; MOVQ AX, (64+24)(DX)
	MOVQ   (64+32)(DX), AX; ADCQ R10, AX; MOVQ AX, (64+32)(DX)
	MOVQ   (64+40)(DX), AX; ADCQ R11, AX; MOVQ AX, (64+40)(DX)
	MOVQ   (64+48)(DX), AX; ADCQ R12, AX; MOVQ AX, (64+48)(DX)
	MOVQ   (64+56)(DX), AX; ADCQ R13, AX; MOVQ AX, (64+56)(DX)

	RET

TEXT ·modP434(SB),NOSPLIT,$0-8
	RET
TEXT ·mulP434(SB),NOSPLIT,$104-24
	RET
TEXT ·rdcP434(SB),$0-16
	RET
