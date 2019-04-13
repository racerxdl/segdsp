package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"github.com/racerxdl/segdsp/tools"
	"math/rand"
	"testing"
)

func NonZeroRand() float32 {
	v := rand.Float32()*2 - 1
	for v == 0 {
		v = rand.Float32()*2 - 1
	}
	return v
}

// region FloatFloat
func TestAddFloatFloatVectors(t *testing.T) {
	if native.GetNativeAddFloatFloatVectors() == nil {
		t.Logf("No Native SIMD AddFloatFloatVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]float32, 64)
	var vecA2 = make([]float32, 64)
	var vecB = make([]float32, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericAddFloatFloatVectors(vecA, vecB)
	native.AddFloatFloatVectors(vecA2, vecB)

	for i := range vecA {
		if vecA[i] != vecA2[i] {
			t.Errorf("Expected %f got %f at %d\n", vecA[i], vecA2[i], i)
		}
	}
}

func TestSubFloatFloatVectors(t *testing.T) {
	if native.GetNativeAddFloatFloatVectors() == nil {
		t.Logf("No Native SIMD SubFloatFloatVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]float32, 64)
	var vecA2 = make([]float32, 64)
	var vecB = make([]float32, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericSubtractFloatFloatVectors(vecA, vecB)
	native.SubtractFloatFloatVectors(vecA2, vecB)

	for i := range vecA {
		if vecA[i] != vecA2[i] {
			t.Errorf("Expected %f got %f at %d\n", vecA[i], vecA2[i], i)
		}
	}
}

func TestMultiplyFloatFloatVectors(t *testing.T) {
	if native.GetNativeMultiplyFloatFloatVectors() == nil {
		t.Logf("No Native SIMD MultiplyFloatFloatVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]float32, 64)
	var vecA2 = make([]float32, 64)
	var vecB = make([]float32, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericMultiplyFloatFloatVectors(vecA, vecB)
	native.MultiplyFloatFloatVectors(vecA2, vecB)

	for i := range vecA {
		if vecA[i] != vecA2[i] {
			t.Errorf("Expected %f got %f at %d\n", vecA[i], vecA2[i], i)
		}
	}
}

func TestDivideFloatFloatVectors(t *testing.T) {
	if native.GetNativeDivideFloatFloatVectors() == nil {
		t.Logf("No Native SIMD DivideloatFloatVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]float32, 64)
	var vecA2 = make([]float32, 64)
	var vecB = make([]float32, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = NonZeroRand()
		vecB[i] = NonZeroRand()
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericDivideFloatFloatVectors(vecA, vecB)
	native.DivideFloatFloatVectors(vecA2, vecB)

	for i := range vecA {
		if vecA[i] != vecA2[i] {
			t.Errorf("Expected %f got %f at %d\n", vecA[i], vecA2[i], i)
		}
	}
}

// endregion
// region ComplexComplex
func TestAddComplexComplexVectors(t *testing.T) {
	if native.GetNativeAddComplexComplexVectors() == nil {
		t.Logf("No Native SIMD AddComplexComplexVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]complex64, 64)
	var vecA2 = make([]complex64, 64)
	var vecB = make([]complex64, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericAddComplexComplexVectors(vecA, vecB)
	native.AddComplexComplexVectors(vecA2, vecB)

	if !tools.Complex64ArrayEqual(vecA, vecA2) {
		t.Errorf("Expected \n%f\n got \n%f\n\n", vecA, vecA2)
	}
}

func TestSubComplexComplexVectors(t *testing.T) {
	if native.GetNativeAddComplexComplexVectors() == nil {
		t.Logf("No Native SIMD SubComplexComplexVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]complex64, 64)
	var vecA2 = make([]complex64, 64)
	var vecB = make([]complex64, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericSubtractComplexComplexVectors(vecA, vecB)
	native.SubtractComplexComplexVectors(vecA2, vecB)

	if !tools.Complex64ArrayEqual(vecA, vecA2) {
		t.Errorf("Expected \n%f\n got \n%f\n\n", vecA, vecA2)
	}
}

func TestMultiplyComplexComplexVectors(t *testing.T) {
	if native.GetNativeMultiplyComplexComplexVectors() == nil {
		t.Logf("No Native SIMD MultiplyComplexComplexVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]complex64, 64)
	var vecA2 = make([]complex64, 64)
	var vecB = make([]complex64, 64)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericMultiplyComplexComplexVectors(vecA, vecB)
	native.MultiplyComplexComplexVectors(vecA2, vecB)

	if !tools.Complex64ArrayEqual(vecA, vecA2) {
		t.Errorf("Expected \n%f\n got \n%f\n\n", vecA, vecA2)
	}
}

func TestDivideComplexComplexVectors(t *testing.T) {
	if native.GetNativeDivideComplexComplexVectors() == nil {
		t.Logf("No Native SIMD DivideloatFloatVectors to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var vecA = make([]complex64, 64)
	var vecA2 = make([]complex64, 64)
	var vecB = make([]complex64, 64)

	t.Logf("Initializing Vectors\n")

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(NonZeroRand(), NonZeroRand())
		vecB[i] = complex(NonZeroRand(), NonZeroRand())
	}

	copy(vecA2, vecA)
	t.Log("Testing Operation\n")

	genericDivideComplexComplexVectors(vecA, vecB)
	native.DivideComplexComplexVectors(vecA2, vecB)

	if !tools.Complex64ArrayEqual(vecA, vecA2) {
		t.Errorf("Expected \n%f\n got \n%f\n\n", vecA, vecA2)
	}
}

// endregion
