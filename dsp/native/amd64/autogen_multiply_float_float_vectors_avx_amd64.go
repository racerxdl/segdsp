//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _multiplyFloatFloatVectorsAVX(A, B unsafe.Pointer, length uint)

func MultiplyFloatFloatVectorsAVX(A, B []float32) {
	var aPtr = unsafe.Pointer(&A[0])
	var bPtr = unsafe.Pointer(&B[0])
	var cLen = uint(len(A))

	if cLen > uint(len(B)) {
		cLen = uint(len(B))
	}

	_multiplyFloatFloatVectorsAVX(aPtr, bPtr, cLen)
}
