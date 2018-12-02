package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math/rand"
	"testing"
)

const epsilon = 10e-6

func CompareComplex(a, b complex64) bool {
	// Safe Compare two complexes within
	// This is needed here because multiplying using SIMD might generate slightly different value
	return Abs(real(a)-real(b)) < epsilon && Abs(imag(a)-imag(b)) < epsilon
}

func Equal(a, b []complex64, t *testing.T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !CompareComplex(v, b[i]) {
			t.Errorf("Difference at %d: Expected %v got %v\n", i, v, b[i])
			// return false
		}
	}
	return true
}

func TestMultiplyConjugate(t *testing.T) {
	if native.GetMultiplyConjugate() == nil {
		t.Logf("No Native SIMD Multiply Conjugate to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]complex64, 128)
	var vecB = make([]complex64, 128)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}
	t.Log("Testing Multiply Conjugate\n")

	var expected = genericMultiplyConjugate(vecA, vecB, 128)
	var got = native.MultiplyConjugate(vecA, vecB, 128)

	if !Equal(expected, got, t) {
		t.Error("Expected != got")
	}
}

func TestMultiplyConjugateInline(t *testing.T) {
	if native.GetMultiplyConjugate() == nil {
		t.Logf("No Native SIMD Multiply Conjugate to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]complex64, 128)
	var vecB = make([]complex64, 128)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	var vecA1 = make([]complex64, 128)
	var vecB1 = make([]complex64, 128)

	copy(vecA1, vecA)
	copy(vecB1, vecB)

	t.Log("Testing Multiply Conjugate Inline\n")

	genericMultiplyConjugateInline(vecA, vecB, 128)
	native.MultiplyConjugateInline(vecA1, vecB1, 128)

	if !Equal(vecA, vecA1, t) {
		t.Error("Expected != got")
	}
}
