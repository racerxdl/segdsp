//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _multiplyConjugateInlineAVX(vecA, vecB unsafe.Pointer, length uint)

func MultiplyConjugateInlineAVX(vecA, vecB []complex64, length int) {
	var vecAPtr = unsafe.Pointer(&vecA[0])
	var vecBPtr = unsafe.Pointer(&vecB[0])

	_multiplyConjugateInlineAVX(vecAPtr, vecBPtr, uint(length))
}
