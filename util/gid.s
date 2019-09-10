#include "textflag.h"
#include "go_asm.h"
#include "go_tls.h"

// func Getg() unsafe.Pointer
TEXT ·Getg(SB), NOSPLIT, $0-8
    get_tls(CX)
    MOVQ    g(CX), AX
    MOVQ    AX, ret+0(FP)
    RET