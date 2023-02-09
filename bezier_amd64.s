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
    MOVSD     magic1<>(SB), X5
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
    MOVQ  ·x0+0(FP), X0
    MOVQ  ·y0+8(FP), X1
    MOVQ  ·x1+16(FP), X2
    MOVQ  ·y1+24(FP), X3
    MOVQ  ·x2+32(FP), X4
    MOVQ  ·y2+40(FP), X5
    MOVQ  ·x3+48(FP), X6
    MOVQ  ·y3+56(FP), X7
    MOVQ  ·p+72(FP), DI
    MOVQ  ·plen+80(FP), SI

    UNPCKLPD  X3, X2
    MOVAPD    X0, X3
    UNPCKLPD  X7, X6
    MOVQ      SI, DX
    MOVSD    ·ds+64(FP), X8
    UNPCKLPD  X5, X4
    UNPCKLPD  X1, X3
    TESTQ     SI, SI
    JLE       return
    MOVSD     magic1<>(SB), X7
    MOVSD     magic2<>(SB), X5
    XORQ      AX, AX
lop:
    PXOR      X9, X9
    MOVAPD    X7, X0
    ADDQ      $16, DI
    CVTSQ2SD  AX, X9
    ADDQ      $1, AX
    DIVSD     X8, X9
    SUBSD     X9, X0
    MOVAPD    X9, X10
    MULSD     X9, X10
    MOVAPD    X0, X11
    MOVAPD    X0, X1
    MULSD     X0, X11
    MULSD     X5, X1
    MULSD     X9, X10
    MULSD     X0, X11
    MULSD     X1, X0
    MULSD     X9, X1
    MULSD     X9, X0
    MULSD     X9, X1
    MOVAPD    X11, X9
    UNPCKLPD  X9, X9
    MULPD     X3, X9
    UNPCKLPD  X0, X0
    MULPD     X2, X0
    UNPCKLPD  X1, X1
    MULPD     X4, X1
    ADDPD     X9, X0
    ADDPD     X1, X0
    MOVAPD    X10, X1
    UNPCKLPD  X1, X1
    MULPD     X6, X1
    ADDPD     X1, X0
    MOVUPS    X0, -16(DI)
    CMPQ      AX, DX
    JNE       lop
return:
    RET
