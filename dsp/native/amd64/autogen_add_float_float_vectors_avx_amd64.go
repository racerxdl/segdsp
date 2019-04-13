//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _addFloatFloatVectorsAVX(A, B unsafe.Pointer, length uint)

func AddFloatFloatVectorsAVX(A, B []float32) {
	var aPtr = unsafe.Pointer(&A[0])
	var bPtr = unsafe.Pointer(&B[0])
	var cLen = uint(len(A))

	if cLen > uint(len(B)) {
		cLen = uint(len(B))
	}

	_addFloatFloatVectorsAVX(aPtr, bPtr, cLen)
}
