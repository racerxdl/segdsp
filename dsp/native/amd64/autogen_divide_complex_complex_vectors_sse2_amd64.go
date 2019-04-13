//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _divideComplexComplexVectorsSSE2(A, B unsafe.Pointer, length uint)

func DivideComplexComplexVectorsSSE2(A, B []complex64) {
	var aPtr = unsafe.Pointer(&A[0])
	var bPtr = unsafe.Pointer(&B[0])
	var cLen = uint(len(A))

	if cLen > uint(len(B)) {
		cLen = uint(len(B))
	}

	_divideComplexComplexVectorsSSE2(aPtr, bPtr, cLen)
}
