//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _dotProductFloatAVX(result, input, taps unsafe.Pointer, length uint)

func DotProductNativeFloatAVX(input []float32, taps []float32) float32 {
	var res = make([]float32, 1)

	var resPtr = unsafe.Pointer(&res[0])
	var inputPtr = unsafe.Pointer(&input[0])
	var tapsPtr = unsafe.Pointer(&taps[0])
	var cLen = uint(len(taps))

	if cLen > uint(len(input)) {
		cLen = uint(len(input))
	}

	_dotProductFloatAVX(resPtr, inputPtr, tapsPtr, cLen)

	return res[0]
}
