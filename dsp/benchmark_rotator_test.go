package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
)

const rotateComplexVecSize = 1 << 20
const phaseShift = (2 * math.Pi) * 1.38

func BenchmarkRotateComplexGolang(b *testing.B) {
	var vecA = make([]complex64, rotateComplexVecSize)
	var phaseIncrement = complex64(cmplx.Exp(complex(0, -phaseShift)))
	var phase = complex64(complex(1, 0))

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		genericRotateComplex(vecA, &phase, phaseIncrement, len(vecA))
	}
}

func BenchmarkRotateComplexNative(b *testing.B) {
	if native.GetNativeRotateComplex() == nil {
		b.Logf("No Native SIMD Rotate Complex to test")
		return
	}
	var vecA = make([]complex64, rotateComplexVecSize)
	var phaseIncrement = complex64(cmplx.Exp(complex(0, -phaseShift)))
	var phase = complex64(complex(1, 0))

	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		native.RotateComplex(vecA, &phase, phaseIncrement, len(vecA))
	}
}
