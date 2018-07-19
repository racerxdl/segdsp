package dsp



type IIRFilter struct {
	fftaps []float32
	fbtaps []float32

	latestN int
	latestM int
	prevOutput []float32
	prevInput []float32
}

func MakeIIRFilter(ataps, btaps []float32) *IIRFilter {
	var fbtaps = make([]float32, len(btaps))
	copy(fbtaps, btaps)

	var n = len(ataps)
	var m = len(btaps)

	for i := 1; i < len(btaps); i++ {
		btaps[i] = -btaps[i]
	}

	return &IIRFilter{
		fftaps: ataps,
		fbtaps: btaps,
		prevInput: make([]float32, 2 * n),
		prevOutput: make([]float32, 2 * m),
	}
}

func (f *IIRFilter) FilterArray(input []float32) []float32 {
	var out = make([]float32, len(input))
	for i := 0; i < len(input); i++ {
		out[i] = f.Filter(input[i])
	}
	return out
}

func (f *IIRFilter) Filter(input float32) float32 {
	var n = len(f.fftaps)
	var m = len(f.fbtaps)

	if n == 0 {
		return 0
	}

	var latestN = f.latestN
	var latestM = f.latestM

	var acc = f.fftaps[0] * input

	for i := 1; i < n; i++ {
		acc += f.fftaps[i] * f.prevInput[latestN + i]
	}
	for i := 1; i < m; i++ {
		acc += f.fbtaps[i] * f.prevOutput[latestM + i]
	}

	f.prevOutput[latestM] = acc
	f.prevOutput[latestM+m] = acc
	f.prevInput[latestN] = input
	f.prevInput[latestN+n] = input

	latestM--
	latestN--

	if latestM < 0 {
		latestM += m
	}

	if latestN < 0 {
		latestN += n
	}

	f.latestN = latestN
	f.latestM = latestM

	return acc
}