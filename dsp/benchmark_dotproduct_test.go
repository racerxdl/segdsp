package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math/rand"
	"testing"
)

func BenchmarkDotProductGolang(b *testing.B) {
	var vecA = make([]complex64, 16384)
	var taps = make([]float32, 128)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		if i < len(taps) {
			taps[i] = rand.Float32()*2 - 1
		}
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericDotProductResult(vecA, taps)
	}
}

func BenchmarkDotProductNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeDotProductComplex() == nil {
		b.Logf("No Native SIMD Complex Dot Product to test")
		return
	}

	var vecA = make([]complex64, 16384)
	var taps = make([]float32, 128)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		if i < len(taps) {
			taps[i] = rand.Float32()*2 - 1
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.DotProductComplex(vecA, taps)
	}
}
