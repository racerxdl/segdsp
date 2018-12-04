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
	var samples = append(f.history, data...)
	var tmp = MultiplyConjugate(samples[1:], samples, len(samples)-2)
	var out = make([]float32, len(samples)-2)

	for i := 0; i < len(out); i++ {
		out[i] = f.gain * tools.ComplexPhase(tmp[i])
	}

	f.history = samples[len(samples)-2:]
	return out
}
