//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _multiplyConjugateAVX(vecA, vecB, output unsafe.Pointer, length uint)

func MultiplyConjugateAVX(vecA, vecB []complex64, length int) []complex64 {
	var output = make([]complex64, length)

	var vecAPtr = unsafe.Pointer(&vecA[0])
	var vecBPtr = unsafe.Pointer(&vecB[0])
	var outPtr = unsafe.Pointer(&output[0])

	_multiplyConjugateAVX(vecAPtr, vecBPtr, outPtr, uint(length))

	return output
}
