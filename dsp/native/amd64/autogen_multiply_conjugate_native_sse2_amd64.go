//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _multiplyConjugateSSE2(vecA, vecB, output unsafe.Pointer, length uint)

func MultiplyConjugateSSE2(vecA, vecB []complex64, length int) []complex64 {
	var output = make([]complex64, length)

	var vecAPtr = unsafe.Pointer(&vecA[0])
	var vecBPtr = unsafe.Pointer(&vecB[0])
	var outPtr = unsafe.Pointer(&output[0])

	_multiplyConjugateSSE2(vecAPtr, vecBPtr, outPtr, uint(length))

	return output
}
