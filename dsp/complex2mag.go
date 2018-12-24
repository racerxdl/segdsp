package dsp

import (
	"math"
)

func Complex2Magnitude(data []complex64) []float32 {
	output := make([]float32, len(data))

	for i := 0; i < len(data); i++ {
		sample := data[i]
		output[i] = float32(math.Sqrt(float64(real(sample)*real(sample) + imag(sample)*imag(sample))))
	}

	return output
}

func Complex2MagnitudeBuffer(input []complex64, output []float32) int {
	if len(input) != len(output) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		sample := input[i]
		output[i] = float32(math.Sqrt(float64(real(sample)*real(sample) + imag(sample)*imag(sample))))
	}

	return len(output)
}
