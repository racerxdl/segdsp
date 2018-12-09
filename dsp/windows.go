package dsp

import "math"

func HammingWindow(nTaps int) []float64 {
	var taps = make([]float64, nTaps)

	var M = float64(nTaps - 1)

	for i := 0; i < nTaps; i++ {
		taps[i] = 0.5 - 0.5*math.Cos((2*math.Pi*float64(i))/M)
	}

	return taps
}

func cosWindow3(nTaps int, c0, c1, c2 float64) []float64 {
	taps := make([]float64, nTaps)
	M := float64(nTaps - 1)

	for i := 0; i < nTaps; i++ {
		taps[i] = c0 - c1*math.Cos((2*math.Pi*float64(i))/M) + c2*math.Cos((4*math.Pi*float64(i))/M)
	}

	return taps
}

func cosWindow4(nTaps int, c0, c1, c2, c3 float64) []float64 {
	taps := make([]float64, nTaps)
	M := float64(nTaps - 1)

	for i := 0; i < nTaps; i++ {
		var a = c0
		var b = c1 * math.Cos((2*math.Pi*float64(i))/M)
		var c = c2 * math.Cos((4*math.Pi*float64(i))/M)
		var d = c3 * math.Cos((6*math.Pi*float64(i))/M)
		taps[i] = a - b + c - d
	}

	return taps
}

//func cosWindow5(nTaps int, c0, c1, c2, c3, c4 float64) []float64 {
//	taps := make([]float64, nTaps)
//	M := float64(nTaps - 1)
//
//	for i := 0; i < nTaps; i++ {
//		var a = c0
//		var b = c1 * math.Cos((2 * math.Pi * float64(i)) / M)
//		var c = c2 * math.Cos((4 * math.Pi * float64(i)) / M)
//		var d = c3 * math.Cos((6 * math.Pi * float64(i)) / M)
//		var e = c3 * math.Cos((8 * math.Pi * float64(i)) / M)
//		taps[i] = a - b + c - d + e
//	}
//
//	return taps
//}

func BlackmanHarris(nTaps, atten int) []float64 {
	switch atten {
	case 61:
		return cosWindow3(nTaps, 0.42323, 0.49755, 0.07922)
	case 67:
		return cosWindow3(nTaps, 0.44959, 0.49364, 0.05677)
	case 74:
		return cosWindow4(nTaps, 0.40271, 0.49703, 0.09392, 0.00183)
	case 92:
		return cosWindow4(nTaps, 0.35875, 0.48829, 0.14128, 0.01168)
	default:
		panic("BlackmanHarris attenuation must be one of the following values: 61, 67, 74, or 92")
	}
}
