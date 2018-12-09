//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _dotProductComplexComplexAVX(result, input, taps unsafe.Pointer, length uint)

func DotProductNativeComplexComplexAVX(input []complex64, taps []complex64) complex64 {
	var res = make([]float32, 2)

	var resPtr = unsafe.Pointer(&res[0])
	var inputPtr = unsafe.Pointer(&input[0])
	var tapsPtr = unsafe.Pointer(&taps[0])
	var cLen = uint(len(taps))

	if cLen > uint(len(input)) {
		cLen = uint(len(input))
	}

	_dotProductComplexComplexAVX(resPtr, inputPtr, tapsPtr, cLen)

	return complex(res[0], res[1])
}
