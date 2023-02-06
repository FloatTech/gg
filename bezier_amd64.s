#include "textflag.h"

DATA magic1<>+0x00(SB)/4, $0x0
DATA magic1<>+0x04(SB)/4, $0x3ff00000
GLOBL magic1<>(SB), RODATA, $8

DATA magic2<>+0x00(SB)/4, $0x0
DATA magic2<>+0x04(SB)/4, $0x40080000
GLOBL magic2<>(SB), RODATA, $8

// func quadratic(x0, y0, x1, y1, x2, y2, t float64) (x, y float64)
TEXT ·quadratic(SB), NOSPLIT, $0-72
    MOVQ  ·x0+0(FP), X0
    MOVQ  ·y0+8(FP), X1
    MOVQ  ·x1+16(FP), X2
    MOVQ  ·y1+24(FP), X3
    MOVQ  ·x2+32(FP), X4
    MOVQ  ·y2+40(FP), X5
    MOVQ  ·t+48(FP), X6
    MOVQ  magic1<>(SB), X7

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

    MOVQ  X0, ·x+56(FP)
    MOVQ  X3, ·y+64(FP)
    RET

// func cubic(x0, y0, x1, y1, x2, y2, x3, y3, t float64) (x, y float64)
TEXT ·cubic(SB), NOSPLIT, $0-88
    MOVQ  ·x0+0(FP), X0
    MOVQ  ·y0+8(FP), X1
    MOVQ  ·x1+16(FP), X2
    MOVQ  ·y1+24(FP), X3
    MOVQ  ·x2+32(FP), X4
    MOVQ  ·y2+40(FP), X5
    MOVQ  ·x3+48(FP), X6
    MOVQ  ·y3+56(FP), X7
    MOVQ  magic1<>(SB), X8
    MOVQ  ·t+64(FP), X11
    MOVQ  magic2<>(SB), X10

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

    MOVQ  X0, ·x+72(FP)
    MOVQ  X1, ·y+80(FP)
    RET
