package dsp

import (
	"github.com/racerxdl/segdsp/dsp/native"
	"github.com/racerxdl/segdsp/tools"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
)

func TestRotator(t *testing.T) {
	if native.GetNativeRotateComplex() == nil {
		t.Logf("No Native SIMD Rotate Complex to test")
		return
	}

	t.Logf("SIMD Mode: %s", native.GetSIMDMode())

	var phaseIncrement = complex64(cmplx.Exp(complex(0, -math.Pi/8)))

	var vecA = make([]complex64, 32)

	t.Logf("Initializing Vectors\n")
	for i := 0; i < len(vecA); i++ {
		vecA[i] = complex(rand.Float32()*2-1, rand.Float32()*2-1)
	}
	t.Log("Testing Rotate Complex\n")

	var phase = complex64(complex(1, 0))

	var expected = genericRotateComplex(vecA, &phase, phaseIncrement, len(vecA))

	phase = complex64(complex(1, 0))
	var got = native.RotateComplex(vecA, &phase, phaseIncrement, len(vecA))

	if !tools.Complex64ArrayEqual(got, expected) {
		t.Errorf("Expected \n%f\n got \n%f\n\n", expected, got)
	}
}
