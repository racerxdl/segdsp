package dsp

import "math"

type FeedForwardAGC struct {
	sampleHistory []complex64
	reference float64
	numSamples int
}

const maxGain = float64(1e-4)

func MakeFeedForwardAGC(numSamples int, reference float32) *FeedForwardAGC {
	return &FeedForwardAGC{
		sampleHistory: make([]complex64, numSamples),
		reference: float64(reference),
		numSamples: numSamples,
	}
}

func (f *FeedForwardAGC) Work(input []complex64) []complex64 {
	output := make([]complex64, len(input))
	gain := float64(0.0)
	samples := append(f.sampleHistory, input...)

	for i := 0; i < len(output); i++ {
		maxEnv := maxGain
		for j := 0; j < len(f.sampleHistory); j++ {
			maxEnv = math.Max(maxEnv, envelope(samples[i+j]))
		}

		gain = f.reference / maxEnv
		output[i] = complex(float32(gain) * real(samples[i]), float32(gain) * imag(samples[i]))
	}

	f.sampleHistory = samples[len(samples)-f.numSamples:]

	return output
}

func envelope(c complex64) float64 {
	realAbs := math.Abs(float64(real(c)))
	imagAbs := math.Abs(float64(real(c)))

	if realAbs > imagAbs {
		return realAbs + 0.4 * imagAbs
	} else {
		return imagAbs + 0.4 * realAbs
	}
}
