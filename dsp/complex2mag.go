package dsp

import (
	"math"
)

func Complex2Magnitude(data []complex64) []float32 {
	output := make([]float32, len(data))

	for i := 0; i < len(data); i++ {
		sample := data[i]
		output[i] = float32(math.Sqrt(float64(real(sample)*real(sample) + imag(sample)*imag(sample))))
		//output[i] = float32(cmplx.Abs(complex128(sample)))
	}

	return output
}
