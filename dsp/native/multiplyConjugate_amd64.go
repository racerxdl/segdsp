package native

import "github.com/racerxdl/segdsp/dsp/native/amd64"

var nativeMultiplyConjugate func(vecA, vecB []complex64, length int) []complex64
var nativeMultiplyConjugateInline func(vecA, vecB []complex64, length int)

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
