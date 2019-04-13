package native

import "github.com/racerxdl/segdsp/dsp/native/amd64"

var nativeDotProductFloat func(input []float32, taps []float32) float32
var nativeDotProductComplex func(input []complex64, taps []float32) complex64
var nativeDotProductComplexComplex func(input []complex64, taps []complex64) complex64
var nativeMultiplyConjugate func(vecA, vecB []complex64, length int) []complex64
var nativeMultiplyConjugateInline func(vecA, vecB []complex64, length int)
var nativeComplexRotate func(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64
var nativeComplexBufferRotate func(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int
var nativeFirFilter func(input []complex64, output []complex64, taps []float32)
var nativeFirFilterDecimate func(decimation uint, input []complex64, output []complex64, taps []float32)
var nativeMultiplyFloatFloatVectors func(A, B []float32)
var nativeDivideFloatFloatVectors func(A, B []float32)
var nativeAddFloatFloatVectors func(A, B []float32)
var nativeSubtractFloatFloatVectors func(A, B []float32)

func RotateComplex(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	if nativeComplexRotate == nil {
		nativeComplexRotate = GetRotateComplex()
	}

	if nativeComplexRotate == nil {
		panic("No native function available for arch")
	}
	return nativeComplexRotate(input, phase, phaseIncrement, length)
}

func RotateComplexBuffer(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	if nativeComplexBufferRotate == nil {
		nativeComplexBufferRotate = GetRotateComplexBuffer()
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
		nativeMultiplyConjugate = GetMultiplyConjugate()
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

func FirFilter(input []complex64, output []complex64, taps []float32) {
	if nativeFirFilter == nil {
		nativeFirFilter = GetFirFilter()
	}

	if nativeFirFilter == nil {
		panic("No native function available for arch")
	}
	nativeFirFilter(input, output, taps)
}

func FirFilterDecimate(decimation uint, input []complex64, output []complex64, taps []float32) {
	if nativeFirFilterDecimate == nil {
		nativeFirFilterDecimate = GetFirFilterDecimate()
	}

	if nativeFirFilterDecimate == nil {
		panic("No native function available for arch")
	}
	nativeFirFilterDecimate(decimation, input, output, taps)
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

func GetMultiplyConjugate() func(vecA, vecB []complex64, length int) []complex64 {
	if amd64.AVX {
		return amd64.MultiplyConjugateAVX
	}

	if amd64.SSE2 {
		return amd64.MultiplyConjugateSSE2
	}

	return nil
}

func GetRotateComplex() func(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	if amd64.AVX {
		return amd64.RotateComplexAVX
	}

	if amd64.SSE2 {
		return amd64.RotateComplexSSE2
	}

	return nil
}

func GetRotateComplexBuffer() func(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	if amd64.AVX {
		return amd64.RotateComplexBufferAVX
	}

	if amd64.SSE2 {
		return amd64.RotateComplexBufferSSE2
	}

	return nil
}

func GetFirFilter() func(input []complex64, output []complex64, taps []float32) {
	if amd64.AVX {
		return amd64.FirFilterAVX
	}

	if amd64.SSE2 {
		return amd64.FirFilterSSE2
	}

	return nil
}

func GetFirFilterDecimate() func(decimation uint, input []complex64, output []complex64, taps []float32) {
	if amd64.AVX {
		return amd64.FirFilterDecimateAVX
	}

	if amd64.SSE2 {
		return amd64.FirFilterDecimateSSE2
	}

	return nil
}

func MultiplyFloatFloatVectors(A, B []float32) {
	if nativeMultiplyFloatFloatVectors == nil {
		nativeMultiplyFloatFloatVectors = GetMultiplyFloatFloatVectors()
	}

	if nativeMultiplyFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeMultiplyFloatFloatVectors(A, B)
}

func DivideFloatFloatVectors(A, B []float32) {
	if nativeDivideFloatFloatVectors == nil {
		nativeDivideFloatFloatVectors = GetDivideFloatFloatVectors()
	}

	if nativeDivideFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeDivideFloatFloatVectors(A, B)
}

func AddFloatFloatVectors(A, B []float32) {
	if nativeAddFloatFloatVectors == nil {
		nativeAddFloatFloatVectors = GetAddFloatFloatVectors()
	}

	if nativeAddFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeAddFloatFloatVectors(A, B)
}

func SubtractFloatFloatVectors(A, B []float32) {
	if nativeSubtractFloatFloatVectors == nil {
		nativeSubtractFloatFloatVectors = GetSubtractFloatFloatVectors()
	}

	if nativeSubtractFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeSubtractFloatFloatVectors(A, B)
}

func GetMultiplyFloatFloatVectors() func(A, B []float32) {
	if amd64.AVX {
		return amd64.MultiplyFloatFloatVectorsAVX
	}

	if amd64.SSE2 {
		return amd64.MultiplyFloatFloatVectorsSSE2
	}

	return nil
}

func GetDivideFloatFloatVectors() func(A, B []float32) {
	if amd64.AVX {
		return amd64.DivideFloatFloatVectorsAVX
	}

	if amd64.SSE2 {
		return amd64.DivideFloatFloatVectorsSSE2
	}

	return nil
}

func GetAddFloatFloatVectors() func(A, B []float32) {
	if amd64.AVX {
		return amd64.AddFloatFloatVectorsAVX
	}

	if amd64.SSE2 {
		return amd64.AddFloatFloatVectorsSSE2
	}

	return nil
}

func GetSubtractFloatFloatVectors() func(A, B []float32) {
	if amd64.AVX {
		return amd64.SubtractFloatFloatVectorsAVX
	}

	if amd64.SSE2 {
		return amd64.SubtractFloatFloatVectorsSSE2
	}

	return nil
}
