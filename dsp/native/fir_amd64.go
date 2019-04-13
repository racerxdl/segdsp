package native

import "github.com/racerxdl/segdsp/dsp/native/amd64"

func FirFilter(input []complex64, output []complex64, taps []float32) {
	if nativeFirFilter == nil {
		nativeFirFilter = GetNativeFirFilter()
	}

	if nativeFirFilter == nil {
		panic("No native function available for arch")
	}
	nativeFirFilter(input, output, taps)
}

func FirFilterDecimate(decimation uint, input []complex64, output []complex64, taps []float32) {
	if nativeFirFilterDecimate == nil {
		nativeFirFilterDecimate = GetNativeFirFilterDecimate()
	}

	if nativeFirFilterDecimate == nil {
		panic("No native function available for arch")
	}
	nativeFirFilterDecimate(decimation, input, output, taps)
}

func GetNativeFirFilter() func(input []complex64, output []complex64, taps []float32) {
	if amd64.AVX {
		return amd64.FirFilterAVX
	}

	if amd64.SSE2 {
		return amd64.FirFilterSSE2
	}

	return nil
}

func GetNativeFirFilterDecimate() func(decimation uint, input []complex64, output []complex64, taps []float32) {
	if amd64.AVX {
		return amd64.FirFilterDecimateAVX
	}

	if amd64.SSE2 {
		return amd64.FirFilterDecimateSSE2
	}

	return nil
}
