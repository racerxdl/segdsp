//+build !noasm
//+build !appengine
package arm64

import (
"unsafe"
)

//go:noescape
func _dotProductFloatNeon(result, input, taps unsafe.Pointer, length uint)

func DotProductNativeFloatNeon(input []float32, taps []float32) float32 {
    var res = make([]float32, 1)

    var resPtr = unsafe.Pointer(&res[0])
    var inputPtr = unsafe.Pointer(&input[0])
    var tapsPtr = unsafe.Pointer(&taps[0])
    var cLen = uint(len(taps))

    if cLen == 0 {
        return 0
    }

    _dotProductFloatNeon(resPtr, inputPtr, tapsPtr, cLen)

    return res[0]
}
