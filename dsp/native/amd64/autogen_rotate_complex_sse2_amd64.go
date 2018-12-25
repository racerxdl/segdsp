//+build !noasm
//+build !appengine

package amd64

import (
	"unsafe"
)

//go:noescape
func _rotateComplexSSE2(inputPtr, outPtr, phaseIncrementPtr, phasePtr unsafe.Pointer, length uint)

func RotateComplexSSE2(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	var output = make([]complex64, length)

	var inputPtr = unsafe.Pointer(&input[0])
	var phasePtr = unsafe.Pointer(phase)
	var phaseIncrementPtr = unsafe.Pointer(&phaseIncrement)
	var outPtr = unsafe.Pointer(&output[0])

	_rotateComplexSSE2(inputPtr, outPtr, phaseIncrementPtr, phasePtr, uint(length))

	return output
}

func RotateComplexBufferSSE2(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	var inputPtr = unsafe.Pointer(&input[0])
	var outPtr = unsafe.Pointer(&output[0])
	var phasePtr = unsafe.Pointer(phase)
	var phaseIncrementPtr = unsafe.Pointer(&phaseIncrement)

	_rotateComplexSSE2(inputPtr, outPtr, phaseIncrementPtr, phasePtr, uint(length))

	return length
}
