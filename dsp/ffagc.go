package dsp

import "math"

type FeedForwardAGC struct {
	sampleHistory []complex64
	reference     float64
	numSamples    int
}

const maxGain = float64(1e-4)

func MakeFeedForwardAGC(numSamples int, reference float32) *FeedForwardAGC {
	return &FeedForwardAGC{
		sampleHistory: make([]complex64, numSamples),
		reference:     float64(reference),
		numSamples:    numSamples,
	}
}

func (f *FeedForwardAGC) Work(input []complex64) []complex64 {
	var gain float64
	output := make([]complex64, len(input))
	samples := append(f.sampleHistory, input...)

	for i := 0; i < len(output); i++ {
		maxEnv := maxGain
		for j := 0; j < len(f.sampleHistory); j++ {
			maxEnv = math.Max(maxEnv, envelope(samples[i+j]))
		}

		gain = f.reference / maxEnv
		output[i] = complex(float32(gain)*real(samples[i]), float32(gain)*imag(samples[i]))
	}

	f.sampleHistory = samples[len(samples)-f.numSamples:]

	return output
}

func (f *FeedForwardAGC) WorkBuffer(input, output []complex64) int {
	var gain float64
	samples := append(f.sampleHistory, input...)

	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		maxEnv := maxGain
		for j := 0; j < len(f.sampleHistory); j++ {
			maxEnv = math.Max(maxEnv, envelope(samples[i+j]))
		}

		gain = f.reference / maxEnv
		output[i] = complex(float32(gain)*real(samples[i]), float32(gain)*imag(samples[i]))
	}

	f.sampleHistory = samples[len(samples)-f.numSamples:]

	return len(input)
}

func (dc *FeedForwardAGC) PredictOutputSize(inputLength int) int {
	return inputLength
}

func envelope(c complex64) float64 {
	realAbs := math.Abs(float64(real(c)))
	imagAbs := math.Abs(float64(real(c)))

	if realAbs > imagAbs {
		return realAbs + 0.4*imagAbs
	} else {
		return imagAbs + 0.4*realAbs
	}
}
