package dsp

import "github.com/racerxdl/segdsp/tools"

type QuadDemod struct {
	gain    float32
	history []complex64
}

func MakeQuadDemod(gain float32) *QuadDemod {
	return &QuadDemod{
		gain:    gain,
		history: make([]complex64, 2),
	}
}

func (f *QuadDemod) Work(data []complex64) []float32 {
	out := make([]float32, f.PredictOutputSize(len(data)))

	f.WorkBuffer(data, out)

	return out
}

func (f *QuadDemod) WorkBuffer(input []complex64, output []float32) int {
	var samples = append(f.history, input...)
	var tmp = MultiplyConjugate(samples[1:], samples, len(samples)-2)

	for i := 0; i < len(input); i++ {
		output[i] = f.gain * tools.ComplexPhase(tmp[i])
	}

	f.history = samples[len(samples)-2:]
	return len(input)
}

func (f *QuadDemod) PredictOutputSize(inputLength int) int {
	return inputLength
}
