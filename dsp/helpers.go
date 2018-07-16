package dsp

func DotProduct(result *complex64, input []complex64, taps []float32, length int) {
	var res [2]float32

	for i := 0; i < length; i++ {
		res[0] += real(input[i]) * taps[i]
		res[1] += imag(input[i]) * taps[i]
	}

	*result = complex(res[0], res[1])
}

func MultiplyConjugate(vecA, vecB []complex64, length int) []complex64 {
	var output = make([]complex64, length)
	for i := 0; i < length; i++ {
		output[i] = vecA[i] * Conj(vecB[i])
	}

	return output
}