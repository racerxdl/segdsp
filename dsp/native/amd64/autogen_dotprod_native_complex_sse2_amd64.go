//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _dotProductComplexSSE2(result, input, taps unsafe.Pointer, length uint)

func DotProductNativeComplexSSE2(input []complex64, taps []float32) complex64 {
	var res = make([]float32, 2)

	var resPtr = unsafe.Pointer(&res[0])
	var inputPtr = unsafe.Pointer(&input[0])
	var tapsPtr = unsafe.Pointer(&taps[0])
	var cLen = uint(len(taps))

	_dotProductComplexSSE2(resPtr, inputPtr, tapsPtr, cLen)

	return complex(res[0], res[1])
}
