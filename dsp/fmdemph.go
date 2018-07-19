package dsp

import "math"

type FMDeemph struct {
	tau float32
	sampleRate float32
	iir *IIRFilter
}

func MakeFMDeemph(tau, sampleRate float32) *FMDeemph {
	var p = 1 / tau
	var pp = math.Tan(float64(p / (sampleRate * 2)))
	var a1 = float32((pp - 1) / (pp + 1))
	var b0 = float32(pp / (1 + pp))
	var b1 = float32(b0)

	var btaps = []float32{b0, b1}
	var ataps = []float32{1, a1}

	return &FMDeemph{
		tau: tau,
		sampleRate: sampleRate,
		iir: MakeIIRFilter(btaps, ataps),
	}
}

func (f *FMDeemph) Work(data []float32) []float32 {
	return f.iir.FilterArray(data)
}