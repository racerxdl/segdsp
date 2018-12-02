//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _multiplyConjugateInlineSSE2(vecA, vecB unsafe.Pointer, length uint)

func MultiplyConjugateInlineSSE2(vecA, vecB []complex64, length int) {
	var vecAPtr = unsafe.Pointer(&vecA[0])
	var vecBPtr = unsafe.Pointer(&vecB[0])

	_multiplyConjugateInlineSSE2(vecAPtr, vecBPtr, uint(length))
}
