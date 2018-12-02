package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math/rand"
	"testing"
	"time"
)

const multiplyConjugateVecSize = 16384 * 16
const multiplyConjugateRuns = 8

func MultiplyConjugateBenchNative(vecA, vecB []complex64) time.Duration {
	var startTime = time.Now()
	for i := 0; i < multiplyConjugateRuns; i++ {
		native.MultiplyConjugateInline(vecA, vecB, len(vecA))
	}
	return time.Since(startTime)
}

func MultiplyConjugateBenchGolang(vecA, vecB []complex64) time.Duration {
	var startTime = time.Now()
	for i := 0; i < multiplyConjugateRuns; i++ {
		genericMultiplyConjugateInline(vecA, vecB, len(vecA))
	}
	return time.Since(startTime)
}

func TestBenchmarkMultiplyConjugateGolang(t *testing.T) {
	var vecA = make([]complex64, multiplyConjugateVecSize)
	var vecB = make([]complex64, multiplyConjugateVecSize)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}
	t.Log("Testing Multiply Conjugate\n")

	timing := MultiplyConjugateBenchGolang(vecA, vecB)
	t.Logf("Took %f s to complete\n", timing.Seconds())
}

func TestBenchmarMultiplyConjugateNative(t *testing.T) {
	if native.GetNativeDotProductComplex() == nil {
		t.Logf("No Native SIMD Complex Dot Product to test")
		return
	}
	var vecA = make([]complex64, multiplyConjugateVecSize)
	var vecB = make([]complex64, multiplyConjugateVecSize)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}
	t.Log("Testing Multiply Conjugate\n")

	timing := MultiplyConjugateBenchNative(vecA, vecB)
	t.Logf("Took %f s to complete\n", timing.Seconds())
}
