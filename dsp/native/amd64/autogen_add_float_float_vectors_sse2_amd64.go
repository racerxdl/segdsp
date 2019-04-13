//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _addFloatFloatVectorsSSE2(A, B unsafe.Pointer, length uint)

func AddFloatFloatVectorsSSE2(A, B []float32) {
	var aPtr = unsafe.Pointer(&A[0])
	var bPtr = unsafe.Pointer(&B[0])
	var cLen = uint(len(A))

	if cLen > uint(len(B)) {
		cLen = uint(len(B))
	}

	_addFloatFloatVectorsSSE2(aPtr, bPtr, cLen)
}
