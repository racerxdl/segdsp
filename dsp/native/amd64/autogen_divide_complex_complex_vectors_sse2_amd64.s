//+build !noasm !appengine
// AUTO-GENERATED BY C2GOASM -- DO NOT EDIT

TEXT ·_divideComplexComplexVectorsSSE2(SB), $0-24

    MOVQ A+0(FP), DI
    MOVQ B+8(FP), SI
    MOVQ length+16(FP), DX

    WORD $0xd285
	JE LBB0_17
    LONG $0xfa83d389; BYTE $0x03
	JA LBB0_5
LBB0_2:
    WORD $0x3145; BYTE $0xc0
LBB0_3:
    QUAD $0x4c01c08300048d43
    WORD $0xc329
LBB0_4:
    QUAD $0x9704100ff3ff508d; QUAD $0xf38f0c100ff3c189
    QUAD $0x1c100ff39614100f; QUAD $0xe4590ff3e2280f8e
    QUAD $0xf3ea590ff3e8280f; QUAD $0xf3cb590ff3d1590f
    QUAD $0xf3db590ff3c3590f; QUAD $0xf3cd580ff3dc580f
    QUAD $0x970c110ff3cb5e0f; QUAD $0xd35e0ff3d05c0ff3
    QUAD $0x02c0838f14110ff3
    LONG $0xffc38348
	JNE LBB0_4
LBB0_17:
    RET
LBB0_5:
    QUAD $0x89c38949ff438d48; QUAD $0x49d1920f41d201c2
    QUAD $0x00000008ba20ebc1; QUAD $0x4104578d48e2f748
    QUAD $0x920fc20148d2900f; QUAD $0x45d0920ff80148d2
    LONG $0x854dc031; BYTE $0xdb
	JNE LBB0_3
    WORD $0x8445; BYTE $0xc9
	JNE LBB0_3
	JNE LBB0_3
    WORD $0xd284
	JNE LBB0_3
    WORD $0x8445; BYTE $0xd2
	JNE LBB0_3
    WORD $0xc084
	JNE LBB0_3
    WORD $0x8445; BYTE $0xd2
	JNE LBB0_3
    LONG $0xde048d48; WORD $0x3948; BYTE $0xf8
	JBE LBB0_14
    LONG $0xdf048d48; WORD $0x3948; BYTE $0xf0
	JA LBB0_2
LBB0_14:
    QUAD $0xbafce08341d88941
    LONG $0x00000001; WORD $0x894d; BYTE $0xc2
LBB0_15:
    QUAD $0x0f8704100fff428d; QUAD $0x0fe0280f10875410
    QUAD $0xc60fe8280f88e2c6; QUAD $0x1c100fd18941ddea
    QUAD $0x280f108674100f86; QUAD $0xcd590f88cec60fcb
    QUAD $0x0fddeec60feb280f; QUAD $0xd6590fcd5c0fec59
    QUAD $0x590fc3590ff6590f; QUAD $0x88e6c60fe3280fdb
    QUAD $0x0fdc580fdddec60f; QUAD $0xc60f88e2c60fe028
    QUAD $0xc35e0fc4580fddc2; QUAD $0x140fd0280fcb5e0f
    QUAD $0x44110f42c1150fd1; QUAD $0xfc8f54110f420c8f
    LONG $0x4908c283; WORD $0xc283; BYTE $0xfc
	JNE LBB0_15
    WORD $0x3949; BYTE $0xd8
	JNE LBB0_3
	JMP LBB0_17