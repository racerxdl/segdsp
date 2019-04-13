package native

import "github.com/racerxdl/segdsp/dsp/native/amd64"

// region Dot Product
var nativeDotProductFloat func(input []float32, taps []float32) float32
var nativeDotProductComplex func(input []complex64, taps []float32) complex64
var nativeDotProductComplexComplex func(input []complex64, taps []complex64) complex64

// endregion
// region Multiply Conjugate
var nativeMultiplyConjugate func(vecA, vecB []complex64, length int) []complex64
var nativeMultiplyConjugateInline func(vecA, vecB []complex64, length int)

// endregion
// region Complex Rotation
var nativeComplexRotate func(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64
var nativeComplexBufferRotate func(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int

// endregion
// region FIR
var nativeFirFilter func(input []complex64, output []complex64, taps []float32)
var nativeFirFilterDecimate func(decimation uint, input []complex64, output []complex64, taps []float32)

// endregion
// Float-Float Vector Operations
var nativeMultiplyFloatFloatVectors func(A, B []float32)
var nativeDivideFloatFloatVectors func(A, B []float32)
var nativeAddFloatFloatVectors func(A, B []float32)
var nativeSubtractFloatFloatVectors func(A, B []float32)

// endregion
// Complex-Complex Vector Operations
var nativeMultiplyComplexComplexVectors func(A, B []complex64)
var nativeDivideComplexComplexVectors func(A, B []complex64)
var nativeAddComplexComplexVectors func(A, B []complex64)
var nativeSubtractComplexComplexVectors func(A, B []complex64)

// endregion

func RotateComplex(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	if nativeComplexRotate == nil {
		nativeComplexRotate = GetNativeRotateComplex()
	}

	if nativeComplexRotate == nil {
		panic("No native function available for arch")
	}
	return nativeComplexRotate(input, phase, phaseIncrement, length)
}

func RotateComplexBuffer(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	if nativeComplexBufferRotate == nil {
		nativeComplexBufferRotate = GetNativeRotateComplexBuffer()
	}

	if nativeComplexBufferRotate == nil {
		panic("No native function available for arch")
	}
	return nativeComplexBufferRotate(input, output, phase, phaseIncrement, length)
}

func DotProductComplex(input []complex64, taps []float32) complex64 {
	if nativeDotProductComplex == nil {
		nativeDotProductComplex = GetNativeDotProductComplex()
	}

	if nativeDotProductComplex == nil {
		panic("No native function available for arch")
	}
	return nativeDotProductComplex(input, taps)
}

func DotProductFloat(input []float32, taps []float32) float32 {
	if nativeDotProductFloat == nil {
		nativeDotProductFloat = GetNativeDotProductFloat()
	}

	if nativeDotProductFloat == nil {
		panic("No native function available for arch")
	}
	return nativeDotProductFloat(input, taps)
}

func DotProductComplexComplex(input []complex64, taps []complex64) complex64 {
	if nativeDotProductComplexComplex == nil {
		nativeDotProductComplexComplex = GetNativeDotProductComplexComplex()
	}

	if nativeDotProductComplexComplex == nil {
		panic("No native function available for arch")
	}
	return nativeDotProductComplexComplex(input, taps)
}

func MultiplyConjugate(vecA, vecB []complex64, length int) []complex64 {
	if nativeMultiplyConjugate == nil {
		nativeMultiplyConjugate = GetNativeMultiplyConjugate()
	}

	if nativeMultiplyConjugate == nil {
		panic("No native function available for arch")
	}
	return nativeMultiplyConjugate(vecA, vecB, length)
}

func MultiplyConjugateInline(vecA, vecB []complex64, length int) {
	if nativeMultiplyConjugateInline == nil {
		nativeMultiplyConjugateInline = GetMultiplyConjugateInline()
	}

	if nativeMultiplyConjugateInline == nil {
		panic("No native function available for arch")
	}
	nativeMultiplyConjugateInline(vecA, vecB, length)
}

func GetNativeDotProductComplex() func(input []complex64, taps []float32) complex64 {
	if amd64.AVX {
		return amd64.DotProductNativeComplexAVX
	}

	if amd64.SSE2 {
		return amd64.DotProductNativeComplexSSE2
	}

	return nil
}

func GetNativeDotProductComplexComplex() func(input []complex64, taps []complex64) complex64 {
	if amd64.AVX {
		return amd64.DotProductNativeComplexComplexAVX
	}

	if amd64.SSE2 {
		return amd64.DotProductNativeComplexComplexSSE2
	}

	return nil
}

func GetNativeDotProductFloat() func(input []float32, taps []float32) float32 {
	if amd64.AVX {
		return amd64.DotProductNativeFloatAVX
	}

	if amd64.SSE2 {
		return amd64.DotProductNativeFloatSSE2
	}

	return nil
}

func GetMultiplyConjugateInline() func(vecA, vecB []complex64, length int) {
	if amd64.AVX {
		return amd64.MultiplyConjugateInlineAVX
	}

	if amd64.SSE2 {
		return amd64.MultiplyConjugateInlineSSE2
	}

	return nil
}

func GetNativeMultiplyConjugate() func(vecA, vecB []complex64, length int) []complex64 {
	if amd64.AVX {
		return amd64.MultiplyConjugateAVX
	}

	if amd64.SSE2 {
		return amd64.MultiplyConjugateSSE2
	}

	return nil
}

func GetNativeRotateComplex() func(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	if amd64.AVX {
		return amd64.RotateComplexAVX
	}

	if amd64.SSE2 {
		return amd64.RotateComplexSSE2
	}

	return nil
}

func GetNativeRotateComplexBuffer() func(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	if amd64.AVX {
		return amd64.RotateComplexBufferAVX
	}

	if amd64.SSE2 {
		return amd64.RotateComplexBufferSSE2
	}

	return nil
}
