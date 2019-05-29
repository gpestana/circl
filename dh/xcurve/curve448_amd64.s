// +build amd64

#include "textflag.h"

// Depends on circl/math/fp448 package
#include "../../math/fp448/fp_amd64.h"
#include "curve_amd64.h"

// CTE_A24 is (A+2)/4 from Curve448
#define CTE_A24 39082

#define Size 56

// multiplyA24Leg multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64, cmov, adx
#define multiplyA24Leg(z,x) \
    MOVQ $CTE_A24, R15; \
    MOVQ  0+x, AX; MULQ R15; MOVQ AX,  R8; ;;;;;;;;;;;;  MOVQ DX,  R9; \
    MOVQ  8+x, AX; MULQ R15; ADDQ AX,  R9; ADCQ $0, DX;  MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R15; ADDQ AX, R10; ADCQ $0, DX;  MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R15; ADDQ AX, R11; ADCQ $0, DX;  MOVQ DX, R12; \
    MOVQ 32+x, AX; MULQ R15; ADDQ AX, R12; ADCQ $0, DX;  MOVQ DX, R13; \
    MOVQ 40+x, AX; MULQ R15; ADDQ AX, R13; ADCQ $0, DX;  MOVQ DX, R14; \
    MOVQ 48+x, AX; MULQ R15; ADDQ AX, R14; ADCQ $0, DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; MOVQ $0, DX; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    ADCQ $0,  DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    MOVQ  R8,  0+z; \
    MOVQ  R9,  8+z; \
    MOVQ R10, 16+z; \
    MOVQ R11, 24+z; \
    MOVQ R12, 32+z; \
    MOVQ R13, 40+z; \
    MOVQ R14, 48+z;

// multiplyA24Adx multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R14, FLAGS
// Instr: x86_64, bmi2
#define multiplyA24Adx(z,x) \
    MOVQ $CTE_A24, DX; \
    MULXQ  0+x, R8,  R9; \
    MULXQ  8+x, AX, R10;  ADDQ AX,  R9; \
    MULXQ 16+x, AX, R11;  ADCQ AX, R10; \
    MULXQ 24+x, AX, R12;  ADCQ AX, R11; \
    MULXQ 32+x, AX, R13;  ADCQ AX, R12; \
    MULXQ 40+x, AX, R14;  ADCQ AX, R13; \
    MULXQ 48+x, AX,  DX;  ADCQ AX, R14; \
    ;;;;;;;;;;;;;;;;;;;;  ADCQ $0,  DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; MOVQ $0, DX; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    ADCQ $0,  DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    MOVQ  R8,  0+z; \
    MOVQ  R9,  8+z; \
    MOVQ R10, 16+z; \
    MOVQ R11, 24+z; \
    MOVQ R12, 32+z; \
    MOVQ R13, 40+z; \
    MOVQ R14, 48+z;

// func ladderStep448(w *[5]fp448.Elt, move uint)
// ladderStep448 calculates a point addition and doubling as follows:
// (x2,z2) = 2*(x2,z2) and (x3,z3) = (x2,z2)+(x3,z3) using as a difference (x1,-).
//    work = {x1,x2,z2,x3,z4} are five fp255.Elt of 56 bytes.
//   stack = (t0,t1) are two fp.Elt of 56 bytes, and
//           (b0,b1) are two fp.bigElt of 56 bytes.
TEXT ·ladderStep448(SB),NOSPLIT,$336-32
    // Parameters
    #define regWork DI
    #define regMove SI
    #define x1 0*Size(regWork)
    #define x2 1*Size(regWork)
    #define z2 2*Size(regWork)
    #define x3 3*Size(regWork)
    #define z3 4*Size(regWork)
    // Local variables
    #define t0 0*Size(SP)
    #define t1 1*Size(SP)
    #define b0 2*Size(SP)
    #define b1 4*Size(SP)
    MOVQ work+0(FP), regWork
    MOVQ move+24(FP), regMove
    CHECK_BMI2ADX(LLADSTEP, ladderStepLeg, ladderStepBmi2Adx)
    #undef regWork
    #undef regMove
    #undef x1
    #undef x2
    #undef z2
    #undef x3
    #undef z3
    #undef t0
    #undef t1
    #undef b0
    #undef b1

// func difAdd448(work *[4]fp.Elt, mu *fp.Elt, swap uint)
// diffAdd calculates a differential point addition using a precomputed point.
// (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
//    work = {x1,z1,x2,z2} are four fp.Elt of 56 bytes, and
//   stack = {b0,b1} are two fp.bigElt of 56 bytes each one.
// See Equation 7 at https://eprint.iacr.org/2017/264.
TEXT ·difAdd448(SB),NOSPLIT,$224-56
    // Parameters
    #define regWork DI
    #define regMu   CX
    #define regSwap SI
    #define ui 0(regMu)
    #define x1 0*Size(regWork)
    #define z1 1*Size(regWork)
    #define x2 2*Size(regWork)
    #define z2 3*Size(regWork)
    // Local variables
    #define b0 0*Size(SP)
    #define b1 2*Size(SP)
    MOVQ work+0(FP), regWork
    MOVQ mu+24(FP), regMu
    MOVQ swap+48(FP), regSwap
    cswap(x1,x2,regSwap)
    cswap(z1,z2,regSwap)
    CHECK_BMI2ADX(LDIFADD, difAddLeg, difAddBmi2Adx)
    #undef regWork
    #undef regMu
    #undef regSwap
    #undef ui
    #undef x1
    #undef z1
    #undef x2
    #undef z2
    #undef b0
    #undef b1

// func double448(work *[4]fp448.Elt)
// double calculates a point doubling (x1,z1) = 2*(x1,z1).
//   work  = {x1,z1,x2,z2} are four fp255.Elt of 32 bytes each one, and
//   stack = {b0,b1} are two fp.bigElt of 56 bytes each one.
// Variables x2,z2 are modified.
TEXT ·double448(SB),NOSPLIT,$224-24
    // Parameters
    #define regWork DI
    #define x1 0*Size(regWork)
    #define z1 1*Size(regWork)
    #define x2 2*Size(regWork)
    #define z2 3*Size(regWork)
    // Local variables
    #define b0 0*Size(SP)
    #define b1 2*Size(SP)
    MOVQ work+0(FP), regWork
    CHECK_BMI2ADX(LDOUB,doubleLeg,doubleBmi2Adx)
    #undef regWork
    #undef x1
    #undef z1
    #undef x2
    #undef z2
    #undef b0
    #undef b1
