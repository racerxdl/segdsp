//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _subtractComplexComplexVectorsAVX(A, B unsafe.Pointer, length uint)

func SubtractComplexComplexVectorsAVX(A, B []complex64) {
	var aPtr = unsafe.Pointer(&A[0])
	var bPtr = unsafe.Pointer(&B[0])
	var cLen = uint(len(A))

	if cLen > uint(len(B)) {
		cLen = uint(len(B))
	}

	_subtractComplexComplexVectorsAVX(aPtr, bPtr, cLen)
}
