//+build !noasm !appengine

TEXT ·_dotProductFloatNeon(SB), $0-32

    MOVD result+0(FP), R0
    MOVD input+8(FP), R1
    MOVD taps+16(FP), R2
    MOVD length+24(FP), R3

    WORD $0x1e2703e0 //    fmov    s0, wzr
    WORD $0x34000003 //    cbz    w3, .LBB0_2
LBB0_1:
    WORD $0xbc404421 //    ldr    s1, [x1], #4
    WORD $0xbc404442 //    ldr    s2, [x2], #4
    WORD $0x51000463 //    sub    w3, w3, #1
    WORD $0x1e220821 //    fmul    s1, s1, s2
    WORD $0x1e212800 //    fadd    s0, s0, s1
    WORD $0x35FFFF63 //    cbnz    w3, .LBB0_1

LBB0_2:
    WORD $0xbd000000 //    str        s0, [x0]
    WORD $0xd65f03c0 //    ret
