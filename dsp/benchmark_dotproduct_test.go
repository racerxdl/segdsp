package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"log"
	"math/rand"
	"testing"
	"time"
)

func DotProductBenchNative(input []complex64, taps []float32) time.Duration {
	var sum = complex64(complex(0, 0))
	var startTime = time.Now()
	for i := 0; i < 10000; i++ {
		expected := native.DotProductComplex(input, taps)
		sum += expected
	}
	log.Printf("Native Sum Result: %f + %fi\n", real(sum), imag(sum))
	return time.Since(startTime)
}

func DotProductBenchGolang(input []complex64, taps []float32) time.Duration {
	var sum = complex64(complex(0, 0))
	var startTime = time.Now()
	for i := 0; i < 10000; i++ {
		expected := genericDotProductResult(input, taps)
		sum += expected
	}
	log.Printf("Golang Sum Result: %f + %fi\n", real(sum), imag(sum))
	return time.Since(startTime)
}

func TestBenchmarkDotProductGolang(t *testing.T) {
	var vecA = make([]complex64, 128)
	var taps = make([]float32, 16)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		if i < len(taps) {
			taps[i] = rand.Float32()*2 - 1
		}
	}
	t.Log("Testing DotProduct\n")

	timing := DotProductBenchGolang(vecA, taps)
	t.Logf("Took %f s to complete\n", timing.Seconds())
}

func TestBenchmarkDotProductNative(t *testing.T) {
	if native.GetNativeDotProductComplex() == nil {
		t.Logf("No Native SIMD Complex Dot Product to test")
		return
	}
	var vecA = make([]complex64, 128)
	var taps = make([]float32, 16)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		if i < len(taps) {
			taps[i] = rand.Float32()*2 - 1
		}
	}
	t.Log("Testing DotProduct\n")

	timing := DotProductBenchNative(vecA, taps)
	t.Logf("Took %f s to complete\n", timing.Seconds())
}
