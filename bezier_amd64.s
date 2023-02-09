#include "textflag.h"

DATA magic1<>+0x00(SB)/4, $0x0
DATA magic1<>+0x04(SB)/4, $0x3ff00000
GLOBL magic1<>(SB), RODATA, $8

DATA magic2<>+0x00(SB)/4, $0x0
DATA magic2<>+0x04(SB)/4, $0x40080000
GLOBL magic2<>(SB), RODATA, $8

// func quadratic(x0, y0, x1, y1, x2, y2, ds float64, p []Point)
TEXT ·quadratic(SB), NOSPLIT, $0-80
    MOVQ  ·x0+0(FP), X0
    MOVQ  ·y0+8(FP), X1
    MOVQ  ·x1+16(FP), X2
    MOVQ  ·y1+24(FP), X3
    MOVQ  ·x2+32(FP), X4
    MOVQ  ·y2+40(FP), X5
    MOVQ  ·ds+48(FP), X6
    MOVQ  ·p+56(FP), DI
    MOVQ  ·plen+64(FP), SI

    UNPCKLPD  X3, X2
    MOVAPD    X0, X3
    UNPCKLPD  X5, X4
    UNPCKLPD  X1, X3
    TESTQ     SI, SI
    JLE       return
    MOVQ      magic1<>(SB), X5
    XORQ      AX, AX
lop:
    PXOR      X1, X1
    MOVAPD    X5, X0
    ADDQ      $16, DI
    CVTSQ2SD  AX, X1
    ADDQ      $1, AX
    DIVSD     X6, X1
    SUBSD     X1, X0
    MOVAPD    X0, X7
    MULSD     X0, X7
    ADDSD     X0, X0
    MULSD     X1, X0
    MULSD     X1, X1
    UNPCKLPD  X7, X7
    MULPD     X3, X7
    UNPCKLPD  X0, X0
    MULPD     X2, X0
    UNPCKLPD  X1, X1
    MULPD     X4, X1
    ADDPD     X7, X0
    ADDPD     X1, X0
    MOVUPS    X0, -16(DI)
    CMPQ      AX, SI
    JNE       lop
return:
    RET

// func cubic(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point)
TEXT ·cubic(SB), NOSPLIT, $0-96
    MOVQ  ·x0+0(FP), AX
    MOVQ  ·y0+8(FP), BX
    MOVQ  ·x1+16(FP), DX
    MOVQ  ·y1+24(FP), SI
    MOVQ  ·x2+32(FP), R8
    MOVQ  ·y2+40(FP), R9
    MOVQ  ·x3+48(FP), R10
    MOVQ  ·y3+56(FP), R11
    MOVQ  ·p+72(FP), DI
    MOVQ  ·plen+80(FP), CX
    XORQ  R12, R12

lop:
    MOVQ  AX, X0
    MOVQ  BX, X1
    MOVQ  DX, X2
    MOVQ  SI, X3
    MOVQ  R8, X4
    MOVQ  R9, X5
    MOVQ  R10, X6
    MOVQ  R11, X7
    MOVQ  magic1<>(SB), X8
    MOVQ  ·ds+64(FP), X9
    MOVQ  magic2<>(SB), X10
    CVTSQ2SD  R12, X11
    DIVSD  X9, X11

    SUBSD  X11, X8
    MOVAPD X11, X12
    MULSD  X11, X12
    MULSD  X8, X10
    MOVAPD X8, X9
    MULSD  X8, X9
    MULSD  X8, X9
    MULSD  X10, X8
    MULSD  X11, X10
    MULSD  X9, X0
    MULSD  X11, X8
    MULSD  X11, X10
    MULSD  X9, X1
    MULSD  X12, X11
    MULSD  X8, X2
    MULSD  X8, X3
    MULSD  X10, X4
    MULSD  X10, X5
    MULSD  X11, X6
    ADDSD  X2, X0
    MULSD  X11, X7
    ADDSD  X3, X1
    ADDSD  X4, X0
    ADDSD  X5, X1
    ADDSD  X6, X0
    ADDSD  X7, X1

    MOVQ  X0, 0(DI*1)
    MOVQ  X1, 8(DI*1)
    ADDQ  $16, DI
    INCQ  R12
    DECQ  CX
    JA    lop
    RET
