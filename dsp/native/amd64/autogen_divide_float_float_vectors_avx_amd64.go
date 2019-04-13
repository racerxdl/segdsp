//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _divideFloatFloatVectorsAVX(A, B unsafe.Pointer, length uint)

func DivideFloatFloatVectorsAVX(A, B []float32) {
	var aPtr = unsafe.Pointer(&A[0])
	var bPtr = unsafe.Pointer(&B[0])
	var cLen = uint(len(A))

	if cLen > uint(len(B)) {
		cLen = uint(len(B))
	}

	_divideFloatFloatVectorsAVX(aPtr, bPtr, cLen)
}
