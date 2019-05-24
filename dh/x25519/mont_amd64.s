// +build amd64

#include "textflag.h"

// Depends on circl/math/fp25519 package
#include "../../math/fp25519/fp_amd64.h"

// CTE_A24 is (A+2)/4 from Curve25519
#define CTE_A24 121666

#define Size 32

// multiplyA24Leg multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R13, FLAGS
// Instr: x86_64, cmov
#define multiplyA24Leg(z,x) \
    MOVL $CTE_A24, AX; MULQ  0+x; MOVQ AX,  R8; MOVQ DX,  R9; \
    MOVL $CTE_A24, AX; MULQ  8+x; MOVQ AX, R12; MOVQ DX, R10; \
    MOVL $CTE_A24, AX; MULQ 16+x; MOVQ AX, R13; MOVQ DX, R11; \
    MOVL $CTE_A24, AX; MULQ 24+x; \
    ADDQ R12,  R9; \
    ADCQ R13, R10; \
    ADCQ  AX, R11; \
    ADCQ  $0,  DX; \
    MOVL $38,  AX; /* 2*C = 38 = 2^256 MOD 2^255-19*/ \
    IMULQ AX, DX; \
    ADDQ DX, R8; \
    ADCQ $0,  R9;  MOVQ  R9,  8+z; \
    ADCQ $0, R10;  MOVQ R10, 16+z; \
    ADCQ $0, R11;  MOVQ R11, 24+z; \
    MOVQ $0, DX; \
    CMOVQCS AX, DX; \
    ADDQ DX, R8;  MOVQ  R8,   0+z;

// multiplyA24Adx multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R12, FLAGS
// Instr: x86_64, cmov, bmi2
#define multiplyA24Adx(z,x) \
    MOVQ  $CTE_A24, DX; \
    MULXQ  0+x,  R8, R10; \
    MULXQ  8+x,  R9, R11;  ADDQ R10,  R9; \
    MULXQ 16+x, R10,  AX;  ADCQ R11, R10; \
    MULXQ 24+x, R11, R12;  ADCQ  AX, R11; \
    ;;;;;;;;;;;;;;;;;;;;;  ADCQ  $0, R12; \
    MOVL $38,  DX; /* 2*C = 38 = 2^256 MOD 2^255-19*/ \
    IMULQ DX, R12; \
    ADDQ R12, R8; \
    ADCQ $0,  R9;  MOVQ  R9,  8+z; \
    ADCQ $0, R10;  MOVQ R10, 16+z; \
    ADCQ $0, R11;  MOVQ R11, 24+z; \
    MOVQ $0, R12; \
    CMOVQCS DX, R12; \
    ADDQ R12, R8;  MOVQ  R8,  0+z;

#define mulA24Legacy \
    multiplyA24Leg(0(DI),0(SI))
#define mulA24Bmi2Adx \
    multiplyA24Adx(0(DI),0(SI))

// func mulA24(z, x *Elt)
TEXT ·mulA24(SB),NOSPLIT,$0-16
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    CHECK_BMI2ADX(LMA24, mulA24Legacy, mulA24Bmi2Adx)


#define ladderStepLeg          \
    addSub(x2,z2)              \
    addSub(x3,z3)              \
    integerMulLeg(b0,x2,z3)    \
    integerMulLeg(b1,x3,z2)    \
    reduceFromDoubleLeg(t0,b0) \
    reduceFromDoubleLeg(t1,b1) \
    addSub(t0,t1)              \
    cselect(x2,x3,regMove)     \
    cselect(z2,z3,regMove)     \
    integerSqrLeg(b0,t0)       \
    integerSqrLeg(b1,t1)       \
    reduceFromDoubleLeg(x3,b0) \
    reduceFromDoubleLeg(z3,b1) \
    integerMulLeg(b0,x1,z3)    \
    reduceFromDoubleLeg(z3,b0) \
    integerSqrLeg(b0,x2)       \
    integerSqrLeg(b1,z2)       \
    reduceFromDoubleLeg(x2,b0) \
    reduceFromDoubleLeg(z2,b1) \
    subtraction(t0,x2,z2)      \
    multiplyA24Leg(t1,t0)      \
    additionLeg(t1,t1,z2)      \
    integerMulLeg(b0,x2,z2)    \
    integerMulLeg(b1,t0,t1)    \
    reduceFromDoubleLeg(x2,b0) \
    reduceFromDoubleLeg(z2,b1)

#define ladderStepBmi2Adx      \
    addSub(x2,z2)              \
    addSub(x3,z3)              \
    integerMulAdx(b0,x2,z3)    \
    integerMulAdx(b1,x3,z2)    \
    reduceFromDoubleAdx(t0,b0) \
    reduceFromDoubleAdx(t1,b1) \
    addSub(t0,t1)              \
    cselect(x2,x3,regMove)     \
    cselect(z2,z3,regMove)     \
    integerSqrAdx(b0,t0)       \
    integerSqrAdx(b1,t1)       \
    reduceFromDoubleAdx(x3,b0) \
    reduceFromDoubleAdx(z3,b1) \
    integerMulAdx(b0,x1,z3)    \
    reduceFromDoubleAdx(z3,b0) \
    integerSqrAdx(b0,x2)       \
    integerSqrAdx(b1,z2)       \
    reduceFromDoubleAdx(x2,b0) \
    reduceFromDoubleAdx(z2,b1) \
    subtraction(t0,x2,z2)      \
    multiplyA24Adx(t1,t0)      \
    additionAdx(t1,t1,z2)      \
    integerMulAdx(b0,x2,z2)    \
    integerMulAdx(b1,t0,t1)    \
    reduceFromDoubleAdx(x2,b0) \
    reduceFromDoubleAdx(z2,b1)

// func ladderStep(w *[5]fp.Elt, move uint)
//  w contains variables used in the Montgomery's ladder step,
//  stack = (t0,t1) are two fp.Elt of fp.Size bytes, and
//          (b0,b1) are two fp.bigElt of 2*fp.Size bytes.
TEXT ·ladderStep(SB),NOSPLIT,$192-16
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
    MOVQ move+8(FP), regMove
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

#define difAddLeg              \
    addSub(x1,z1)              \
    integerMulLeg(b0,z1,ui)    \
    reduceFromDoubleLeg(z1,b0) \
    addSub(x1,z1)              \
    integerSqrLeg(b0,x1)       \
    integerSqrLeg(b1,z1)       \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1) \
    integerMulLeg(b0,x1,z2)    \
    integerMulLeg(b1,z1,x2)    \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1)

#define difAddBmi2Adx          \
    addSub(x1,z1)              \
    integerMulAdx(b0,z1,ui)    \
    reduceFromDoubleAdx(z1,b0) \
    addSub(x1,z1)              \
    integerSqrAdx(b0,x1)       \
    integerSqrAdx(b1,z1)       \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1) \
    integerMulAdx(b0,x1,z2)    \
    integerMulAdx(b1,z1,x2)    \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1)


// func difAdd(work *[4]fp.Elt, mu *fp.Elt, swap uint)
// Calculates a differential point addition
//  (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
//   work  = {x1,z1,x2,z2} are four fp.Elt of fp.Size bytes, and
//   stack = {b0,b1} are two fp.bigElt of 2*fp.Size bytes.
TEXT ·difAdd(SB),NOSPLIT,$128-24
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
    MOVQ mu+8(FP), regMu
    MOVQ swap+16(FP), regSwap
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

#define doubleLeg              \
    addSub(x1,z1)              \
    integerSqrLeg(b0,x1)       \
    integerSqrLeg(b1,z1)       \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1) \
    subtraction(x2,x1,z1)      \
    multiplyA24Leg(z2,x2)      \
    additionLeg(z2,z2,z1)      \
    integerMulLeg(b0,x1,z1)    \
    integerMulLeg(b1,x2,z2)    \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1)

#define doubleBmi2Adx          \
    addSub(x1,z1)              \
    integerSqrAdx(b0,x1)       \
    integerSqrAdx(b1,z1)       \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1) \
    subtraction(x2,x1,z1)      \
    multiplyA24Adx(z2,x2)      \
    additionAdx(z2,z2,z1)      \
    integerMulAdx(b0,x1,z1)    \
    integerMulAdx(b1,x2,z2)    \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1)


// func double(work *[4]fp.Elt)
// Calculates a point doubling (x1,z1) = 2*(x1,z1).
//   work  = {x1,z1,x2,z2} are four fp.Elt of fp.Size bytes, and
//   stack = {b0,b1} are two fp.bigElt of 2*fp.Size bytes.
// Variables x2,z2 are modified.
TEXT ·double(SB),NOSPLIT,$128-8
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
