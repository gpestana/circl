// +build amd64

#include "textflag.h"
#include "fp_amd64.h"

// func Cmov(x, y *Elt, n uint)
TEXT ·Cmov(SB),NOSPLIT,$0-24
    MOVQ x+0(FP), DI
    MOVQ y+8(FP), SI
    MOVQ n+16(FP), BX
    cselect(0(DI),0(SI),BX)
    RET

// func Cswap(x, y *Elt, n uint)
TEXT ·Cswap(SB),NOSPLIT,$0-24
    MOVQ x+0(FP), DI
    MOVQ y+8(FP), SI
    MOVQ n+16(FP), BX
    cswap(0(DI),0(SI),BX)
    RET

// func Sub(z, x, y *Elt)
TEXT ·Sub(SB),NOSPLIT,$0-24
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    MOVQ y+16(FP), BX
    subtraction(0(DI),0(SI),0(BX))
    RET

// func AddSub(x, y *Elt)
TEXT ·AddSub(SB),NOSPLIT,$0-16
    MOVQ x+0(FP), DI
    MOVQ y+8(FP), SI
    addSub(0(DI),0(SI))
    RET

#define addLegacy \
    additionLeg(0(DI),0(SI),0(BX))
#define addBmi2Adx \
    additionAdx(0(DI),0(SI),0(BX))

#define mulLegacy \
    integerMulLeg(0(SP),0(SI),0(BX)) \
    reduceFromDoubleLeg(0(DI),0(SP))
#define mulBmi2Adx \
    integerMulAdx(0(SP),0(SI),0(BX)) \
    reduceFromDoubleAdx(0(DI),0(SP))

#define sqrLegacy \
    integerSqrLeg(0(SP),0(SI)) \
    reduceFromDoubleLeg(0(DI),0(SP))
#define sqrBmi2Adx \
    integerSqrAdx(0(SP),0(SI)) \
    reduceFromDoubleAdx(0(DI),0(SP))

// func Add(z, x, y *Elt)
TEXT ·Add(SB),NOSPLIT,$0-24
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    MOVQ y+16(FP), BX
    CHECK_BMI2ADX(LADD, addLegacy, addBmi2Adx)

// func Mul(z, x, y *Elt)
TEXT ·Mul(SB),NOSPLIT,$112-24
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    MOVQ y+16(FP), BX
    CHECK_BMI2ADX(LMUL, mulLegacy, mulBmi2Adx)

// func Sqr(z, x *Elt)
TEXT ·Sqr(SB),NOSPLIT,$112-16
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    CHECK_BMI2ADX(LSQR, sqrLegacy, sqrBmi2Adx)
