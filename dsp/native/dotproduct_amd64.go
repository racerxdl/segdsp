package native

import "github.com/racerxdl/segdsp/dsp/native/amd64"

var nativeDotProductFloat func(input []float32, taps []float32) float32
var nativeDotProductComplex func(input []complex64, taps []float32) complex64
var nativeDotProductComplexComplex func(input []complex64, taps []complex64) complex64

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
