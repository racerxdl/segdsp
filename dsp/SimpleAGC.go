package dsp

import (
	"github.com/racerxdl/segdsp/tools"
)

type SimpleAGC struct {
	rate      float32
	reference float32
	gain      float32
	maxGain   float32
}

func MakeSimpleAGCSimple() *SimpleAGC {
	return MakeSimpleAGC(1e-4, 1.0, 1.0, 0.0)
}

func MakeSimpleAGC(rate, reference, gain, maxGain float32) *SimpleAGC {
	return &SimpleAGC{
		rate:      rate,
		reference: reference,
		gain:      gain,
		maxGain:   maxGain,
	}
}

func (sa *SimpleAGC) scale(input complex64) complex64 {
	output := input * complex(sa.gain, 0)

	sa.gain += sa.rate * (sa.reference - tools.ComplexAbs(output))
	if sa.maxGain > 0 && sa.gain > sa.maxGain {
		sa.gain = sa.maxGain
	}

	return output
}

func (sa *SimpleAGC) Work(input []complex64) []complex64 {
	output := make([]complex64, len(input))

	for i := 0; i < len(output); i++ {
		output[i] = sa.scale(input[i])
	}

	return output
}

func (sa *SimpleAGC) WorkBuffer(input, output []complex64) int {
	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		output[i] = sa.scale(input[i])
	}

	return len(input)
}

func (sa *SimpleAGC) PredictOutputSize(inputLength int) int {
	return inputLength
}
