package native

func FirFilter(input []complex64, output []complex64, taps []float32) {
	if nativeFirFilter == nil {
		nativeFirFilter = GetNativeFirFilter()
	}

	if nativeFirFilter == nil {
		panic("No native function available for arch")
	}
	nativeFirFilter(input, output, taps)
}

func FirFilterDecimate(decimation uint, input []complex64, output []complex64, taps []float32) {
	if nativeFirFilterDecimate == nil {
		nativeFirFilterDecimate = GetNativeFirFilterDecimate()
	}

	if nativeFirFilterDecimate == nil {
		panic("No native function available for arch")
	}
	nativeFirFilterDecimate(decimation, input, output, taps)
}

func GetNativeFirFilter() func(input []complex64, output []complex64, taps []float32) {
	// TODO
	return nil
}

func GetNativeFirFilterDecimate() func(decimation uint, input []complex64, output []complex64, taps []float32) {
	// TODO
	return nil
}
