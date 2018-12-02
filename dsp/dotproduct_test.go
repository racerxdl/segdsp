package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math/rand"
	"testing"
)

func TestDotProductComplex(t *testing.T) {
	if native.GetNativeDotProductComplex() == nil {
		t.Logf("No Native SIMD Complex Dot Product to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

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

	var expected = genericDotProductResult(vecA, taps)
	var got = native.DotProductComplex(vecA, taps)

	if expected != got {
		t.Errorf("Expected %f + %fi got %f + %fi\n", real(expected), imag(expected), real(got), imag(got))
	}
}

func TestDotProductFloat(t *testing.T) {
	if native.GetNativeDotProductFloat() == nil {
		t.Logf("No Native SIMD Float Dot Product to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]float32, 128)
	var taps = make([]float32, 16)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		if i < len(taps) {
			taps[i] = rand.Float32()*2 - 1
		}
	}
	t.Log("Testing DotProduct\n")

	var expected = genericDotProductFloatResult(vecA, taps)
	var got = native.DotProductFloat(vecA, taps)

	if expected != got {
		t.Errorf("Expected %f got %f\n", expected, got)
	}
}
