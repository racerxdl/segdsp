package dsp

type SinglePoleIIRFilter struct {
	alpha          float32
	alpham1        float32
	previousOutput float32
}

func MakeSinglePoleIIRFilter(alpha float32) *SinglePoleIIRFilter {
	if alpha < 0 || alpha > 1 {
		panic("Alpha should be between 0 and 1")
	}

	return &SinglePoleIIRFilter{
		alpha:          alpha,
		alpham1:        1.0 - alpha,
		previousOutput: 0,
	}
}

func (f *SinglePoleIIRFilter) Filter(input float32) float32 {
	output := f.alpha*input + f.alpham1*f.previousOutput
	f.previousOutput = output
	return output
}

func (f *SinglePoleIIRFilter) FilterArray(input []float32) []float32 {
	var out = make([]float32, len(input))
	for i := 0; i < len(input); i++ {
		out[i] = f.Filter(input[i])
	}
	return out
}

func (f *SinglePoleIIRFilter) FilterArrayBuffer(input, output []float32) int {
	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		output[i] = f.Filter(input[i])
	}

	return len(input)
}

func (f *SinglePoleIIRFilter) SetTaps(alpha float32) {
	f.alpha = alpha
	f.alpham1 = 1.0 - alpha
}

func (f *SinglePoleIIRFilter) GetPreviousOutput() float32 {
	return f.previousOutput
}

type IIRFilter struct {
	fftaps []float32
	fbtaps []float32

	latestN    int
	latestM    int
	prevOutput []float32
	prevInput  []float32
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
		fftaps:     ataps,
		fbtaps:     btaps,
		latestM:    0,
		latestN:    0,
		prevInput:  make([]float32, 2*n),
		prevOutput: make([]float32, 2*m),
	}
}

func (f *IIRFilter) FilterArray(input []float32) []float32 {
	var out = make([]float32, len(input))
	for i := 0; i < len(input); i++ {
		out[i] = f.Filter(input[i])
	}
	return out
}

func (f *IIRFilter) FilterArrayBuffer(input, output []float32) int {
	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		output[i] = f.Filter(input[i])
	}

	return len(input)
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
		acc += f.fftaps[i] * f.prevInput[latestN+i]
	}
	for i := 1; i < m; i++ {
		acc += f.fbtaps[i] * f.prevOutput[latestM+i]
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
