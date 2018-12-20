//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _firFilterSSE2(result, input, taps unsafe.Pointer, lengthTaps, length uint)

//go:noescape
func _firFilterDecimateSSE2(result, input, taps unsafe.Pointer, decimate, lengthTaps, length uint)

func FirFilterSSE2(input []complex64, output []complex64, taps []float32) {
	var inputPtr = unsafe.Pointer(&input[0])
	var outputPtr = unsafe.Pointer(&output[0])
	var tapsPtr = unsafe.Pointer(&taps[0])

	var iLen = uint(len(input))
	var oLen = uint(len(output))
	var tLen = uint(len(taps))

	var workLen = iLen

	if workLen > oLen {
		workLen = oLen
	}

	_firFilterSSE2(outputPtr, inputPtr, tapsPtr, tLen, workLen)
}

func FirFilterDecimateSSE2(decimation uint, input []complex64, output []complex64, taps []float32) {
	var inputPtr = unsafe.Pointer(&input[0])
	var outputPtr = unsafe.Pointer(&output[0])
	var tapsPtr = unsafe.Pointer(&taps[0])

	var iLen = uint(len(input))
	var oLen = uint(len(output))
	var tLen = uint(len(taps))

	var workLen = iLen

	if workLen > oLen {
		workLen = oLen
	}

	_firFilterDecimateSSE2(outputPtr, inputPtr, tapsPtr, decimation, tLen, workLen)
}
