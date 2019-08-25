package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math/rand"
	"testing"
)

// region Complex-Complex Vector Operations
func BenchmarkAddComplexComplexVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericAddComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkAddComplexComplexVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeAddFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.AddComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkSubComplexComplexVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericSubtractComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkSubComplexComplexVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeSubtractFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.SubtractComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkMultiplyComplexComplexVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericMultiplyComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkMultiplyComplexComplexVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeMultiplyFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.MultiplyComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkDivideComplexComplexVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericDivideComplexComplexVectors(vecA, vecB)
	}
}

func BenchmarkDivideComplexComplexVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeDivideFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]complex64, 16384)
	var vecB = make([]complex64, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
		vecB[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.DivideComplexComplexVectors(vecA, vecB)
	}
}

// endregion
// region Complex-Complex Vector Operations
func BenchmarkAddFloatFloatVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericAddFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkAddFloatFloatVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeAddFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.AddFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkSubFloatFloatVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericSubtractFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkSubFloatFloatVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeSubtractFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.SubtractFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkMultiplyFloatFloatVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericMultiplyFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkMultiplyFloatFloatVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeMultiplyFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.MultiplyFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkDivideFloatFloatVectorsGolang(b *testing.B) {
	b.StopTimer()
	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericDivideFloatFloatVectors(vecA, vecB)
	}
}

func BenchmarkDivideFloatFloatVectorsNative(b *testing.B) {
	b.StopTimer()

	if native.GetNativeDivideFloatFloatVectors() == nil {
		b.Logf("No Native SIMD Operation to test")
		return
	}

	var vecA = make([]float32, 16384)
	var vecB = make([]float32, 16384)

	for i := 0; i < len(vecA); i++ {
		vecA[i] = rand.Float32()*2 - 1
		vecB[i] = rand.Float32()*2 - 1
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		native.DivideFloatFloatVectors(vecA, vecB)
	}
}

// endregion
