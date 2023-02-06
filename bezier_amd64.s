#include "textflag.h"

DATA magic1<>+0x00(SB)/4, $0x0
DATA magic1<>+0x04(SB)/4, $0x3ff00000
GLOBL magic1<>(SB), RODATA, $8

DATA magic2<>+0x00(SB)/4, $0x0
DATA magic2<>+0x04(SB)/4, $0x40080000
GLOBL magic2<>(SB), RODATA, $8

// func quadratic(x0, y0, x1, y1, x2, y2, ds float64, p []Point)
TEXT ·quadratic(SB), NOSPLIT, $0-72
    MOVQ  ·x0+0(FP), AX
    MOVQ  ·y0+8(FP), BX
    MOVQ  ·x1+16(FP), DX
    MOVQ  ·y1+24(FP), SI
    MOVQ  ·x2+32(FP), R8
    MOVQ  ·y2+40(FP), R9
    MOVQ  ·ds+48(FP), R10
    MOVQ  magic1<>(SB), R11
    MOVQ  ·p+56(FP), DI
    MOVQ  ·plen+64(FP), CX
    XORQ  R12, R12

lop:
    MOVQ  AX, X0
    MOVQ  BX, X1
    MOVQ  DX, X2
    MOVQ  SI, X3
    MOVQ  R8, X4
    MOVQ  R9, X5

    CVTSQ2SD  R12, X6
    MOVQ  R10, X8
    DIVSD X8, X6

    MOVQ  R11, X7

    SUBSD  X6, X7
    MOVAPD X7, X8
    MULSD  X7, X8
    ADDSD  X7, X7
    MULSD  X6, X7
    MULSD  X6, X6
    MULSD  X8, X0
    MULSD  X8, X1
    MULSD  X7, X2
    MULSD  X7, X3
    MULSD  X6, X4
    MULSD  X6, X5
    ADDSD  X2, X0
    ADDSD  X1, X3
    ADDSD  X4, X0
    ADDSD  X5, X3

    MOVQ  X0, 0(DI*1)
    MOVQ  X3, 8(DI*1)
    ADDQ  $16, DI
    INCQ  R12
    DECQ  CX
    JA    lop
    RET

// func cubic(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point)
TEXT ·cubic(SB), NOSPLIT, $0-88
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
