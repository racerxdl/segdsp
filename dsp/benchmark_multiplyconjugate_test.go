package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math/rand"
	"testing"
)

const multiplyConjugateVecSize = 1 << 20

func BenchmarkMultiplyConjugateGolang(b *testing.B) {
	var vecA = make([]complex64, multiplyConjugateVecSize)
	var vecB = make([]complex64, multiplyConjugateVecSize)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericMultiplyConjugateInline(vecA, vecB, len(vecA))
	}
}

func BenchmarkMultiplyConjugateNative(b *testing.B) {
	if native.GetNativeDotProductComplex() == nil {
		b.Logf("No Native SIMD Complex Dot Product to test")
		return
	}
	var vecA = make([]complex64, multiplyConjugateVecSize)
	var vecB = make([]complex64, multiplyConjugateVecSize)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		native.MultiplyConjugateInline(vecA, vecB, len(vecA))
	}
}
