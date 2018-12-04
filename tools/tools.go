package tools

import "math"

func NextPowerOf2(x int) int {
	if IsPowerOf2(x) {
		return x
	}

	return int(math.Pow(2, math.Ceil(math.Log2(float64(x)))))
}

// IsPowerOf2 returns true if x is a power of 2, else false.
func IsPowerOf2(x int) bool {
	return x&(x-1) == 0
}

// ZeroPadC64 returns x with zeros appended to the end to the specified length.
// If len(x) >= length, x is returned, otherwise a new array is returned.
func ZeroPadC64(x []complex64, length int) []complex64 {
	if len(x) >= length {
		return x
	}

	r := make([]complex64, length)
	copy(r, x)
	return r
}

// ToComplex64Array returns the complex equivalent of the real-valued slice.
func ToComplex64Array(x []float32) []complex64 {
	y := make([]complex64, len(x))
	for n, v := range x {
		y[n] = complex(v, 0)
	}
	return y
}

// ToComplex64Array2 returns the complex equivalent of the real-valued matrix.
func ToComplex64Array2(x [][]float32) [][]complex64 {
	y := make([][]complex64, len(x))
	for n, v := range x {
		y[n] = ToComplex64Array(v)
	}
	return y
}
