package dsp

// region Complex Fir Filter

type CTFirFilter struct {
	taps          []complex64
	sampleHistory []complex64
	tapsLen       int
	decimation    int
}

func MakeCTFirFilter(taps []complex64) *CTFirFilter {
	return &CTFirFilter{
		taps:          taps,
		sampleHistory: make([]complex64, len(taps)),
		tapsLen:       len(taps),
		decimation:    1,
	}
}

func MakeDecimationCTFirFilter(decimation int, taps []complex64) *CTFirFilter {
	return &CTFirFilter{
		taps:          taps,
		sampleHistory: make([]complex64, len(taps)),
		tapsLen:       len(taps),
		decimation:    decimation,
	}
}

func (f *CTFirFilter) Filter(data []complex64, length int) {
	var samples = append(f.sampleHistory, data...)
	for i := 0; i < length; i++ {
		ComplexDotProduct(&data[i], samples[i:i+f.tapsLen], f.taps)
	}
	f.sampleHistory = data[len(data)-f.tapsLen:]
}

func (f *CTFirFilter) FilterOut(data []complex64) []complex64 {
	var samples = append(f.sampleHistory, data...)
	var output = make([]complex64, len(data))
	var length = len(samples) - f.tapsLen
	for i := 0; i < length; i++ {
		output[i] = ComplexDotProductResult(samples[i:], f.taps)
	}
	f.sampleHistory = samples[length:]
	return output
}

func (f *CTFirFilter) FilterBuffer(input, output []complex64) int {
	var samples = append(f.sampleHistory, input...)
	var length = len(samples) - f.tapsLen

	if len(output) < length {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < length; i++ {
		output[i] = ComplexDotProductResult(samples[i:], f.taps)
	}
	f.sampleHistory = samples[length:]

	return length
}

func (f *CTFirFilter) Work(data []complex64) []complex64 {
	if f.decimation > 1 {
		return f.FilterDecimateOut(data, f.decimation)
	}
	return f.FilterOut(data)
}

func (f *CTFirFilter) WorkBuffer(input, output []complex64) int {
	if f.decimation > 1 {
		return f.FilterDecimateBuffer(input, output, f.decimation)
	}
	return f.FilterBuffer(input, output)
}

func (f *CTFirFilter) FilterSingle(data []complex64) complex64 {
	return ComplexDotProductResult(data, f.taps)
}

func (f *CTFirFilter) FilterDecimate(data []complex64, decimate int, length int) {
	var samples = append(f.sampleHistory, data...)
	var j = 0
	for i := 0; i < length; i++ {
		ComplexDotProduct(&data[i], samples[j:], f.taps)
		j += decimate
	}
	f.sampleHistory = data[len(data)-f.tapsLen:]
}

func (f *CTFirFilter) FilterDecimateOut(data []complex64, decimate int) []complex64 {
	var samples = append(f.sampleHistory, data...)
	var length = len(data) / decimate
	var output = make([]complex64, length)
	for i := 0; i < length; i++ {
		var srcIdx = decimate * i
		var sl = samples[srcIdx:]
		if len(sl) < len(f.taps) {
			break
		}
		output[i] = ComplexDotProductResult(sl, f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]
	return output
}

func (f *CTFirFilter) FilterDecimateBuffer(input, output []complex64, decimate int) int {
	var samples = append(f.sampleHistory, input...)
	var length = len(input) / decimate

	if len(output) < length {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < length; i++ {
		var srcIdx = decimate * i
		var sl = samples[srcIdx:]
		if len(sl) < len(f.taps) {
			break
		}
		output[i] = ComplexDotProductResult(sl, f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]
	return length
}

func (f *CTFirFilter) SetTaps(taps []complex64) {
	f.taps = taps
}

func (f *CTFirFilter) PredictOutputSize(inputLength int) int {
	return inputLength / f.decimation
}

// endregion
