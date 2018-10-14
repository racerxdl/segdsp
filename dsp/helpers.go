package dsp

import "math"

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ComplexDotProduct(result *complex64, input []complex64, taps []complex64) {
	var length = Min(len(taps), len(input))

	var res = complex64(complex(0, 0))

	for i := 0; i < length; i++ {
		var r = real(input[i]) * real(taps[i]) - imag(input[i]) * imag(taps[i])
		var i = real(input[i]) * imag(taps[i]) + imag(input[i]) * real(taps[i])

		res += complex(r, i)
	}

	*result = res
}

func ComplexDotProductResult(input []complex64, taps []complex64) complex64 {
	var length = Min(len(taps), len(input))

	var res = complex64(complex(0, 0))

	for i := 0; i < length; i++ {
		var r = real(input[i]) * real(taps[i]) - imag(input[i]) * imag(taps[i])
		var i = real(input[i]) * imag(taps[i]) + imag(input[i]) * real(taps[i])

		res += complex(r, i)
	}

	return res
}

func DotProduct(result *complex64, input []complex64, taps []float32) {
	var length = Min(len(taps), len(input))
	var res [2]float32

	for i := 0; i < length; i++ {
		res[0] += real(input[i]) * taps[i]
		res[1] += imag(input[i]) * taps[i]
	}

	*result = complex(res[0], res[1])
}

func DotProductResult(input []complex64, taps []float32) complex64 {
	var length = Min(len(taps), len(input))
	var res [2]float32

	for i := 0; i < length; i++ {
		res[0] += real(input[i]) * taps[i]
		res[1] += imag(input[i]) * taps[i]
	}

	return complex(res[0], res[1])
}

func DotProductFloat(result *float32, input []float32, taps []float32) {
	var res = float32(0.0)
	var length = Min(len(taps), len(input))

	for i := 0; i < length; i++ {
		res += input[i] * taps[i]
	}

	*result = res
}

func DotProductFloatResult(input []float32, taps []float32) float32 {
	var res = float32(0.0)
	var length = Min(len(taps), len(input))

	for i := 0; i < length; i++ {
		res += input[i] * taps[i]
	}

	return res
}

func MultiplyConjugate(vecA, vecB []complex64, length int) []complex64 {
	var output = make([]complex64, length)
	for i := 0; i < length; i++ {
		output[i] = vecA[i] * Conj(vecB[i])
	}

	return output
}

func MultiplyConjungateInline(vecA, vecB []complex64, length int) {
	for i := 0; i < length; i++ {
		vecA[i] = vecA[i] * Conj(vecB[i])
	}
}

func Modulus(c complex64) float32 {
	return float32(math.Sqrt(float64(real(c)*real(c) + imag(c)*imag(c))))
}

func Divide(c complex64, f float32) complex64 {
	var b = 1 / f
	return complex(real(c)*b, imag(c)*b)
}

func Argument(c complex64) float32 {
	return float32(math.Atan2(float64(imag(c)), float64(real(c))))
}
