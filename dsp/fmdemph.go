package dsp

import "math"

type FMDeemph struct {
	tau        float32
	sampleRate float32
	iir        *IIRFilter
}

func MakeFMDeemph(tau, sampleRate float32) *FMDeemph {
	var p = 1 / tau

	var ca = 2 * float64(sampleRate) * math.Tan(float64(p/(sampleRate*2)))

	var k = -ca / (2 * float64(sampleRate))

	var a1 = float32((1.0 + k) / (1.0 - k))
	var b0 = float32(-k / (1.0 - k))

	var btaps = []float32{b0, b0}
	var ataps = []float32{1, -a1}

	return &FMDeemph{
		tau:        tau,
		sampleRate: sampleRate,
		iir:        MakeIIRFilter(btaps, ataps),
	}
}

func (f *FMDeemph) Work(data []float32) []float32 {
	return f.iir.FilterArray(data)
}

func (f *FMDeemph) WorkBuffer(input, output []float32) int {
	return f.iir.FilterArrayBuffer(input, output)
}
