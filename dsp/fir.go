package dsp

// region Complex Fir Filter

type FirFilter struct {
	taps          []float32
	sampleHistory []complex64
	tapsLen       int
	decimation    int
}

func MakeFirFilter(taps []float32) *FirFilter {
	return &FirFilter{
		taps:          taps,
		sampleHistory: make([]complex64, len(taps)),
		tapsLen:       len(taps),
		decimation:    1,
	}
}

func MakeDecimationFirFilter(decimation int, taps []float32) *FirFilter {
	return &FirFilter{
		taps:          taps,
		sampleHistory: make([]complex64, len(taps)),
		tapsLen:       len(taps),
		decimation:    decimation,
	}
}

func (f *FirFilter) Filter(data []complex64, length int) {
	var samples = append(f.sampleHistory, data...)
	for i := 0; i < length; i++ {
		DotProduct(&data[i], samples[i:i+f.tapsLen], f.taps)
	}
	f.sampleHistory = data[len(data)-f.tapsLen:]
}

func (f *FirFilter) FilterOut(data []complex64) []complex64 {
	var samples = append(f.sampleHistory, data...)
	var output = make([]complex64, len(data))
	var length = len(samples) - f.tapsLen
	for i := 0; i < length; i++ {
		output[i] = DotProductResult(samples[i:], f.taps)
	}
	f.sampleHistory = samples[length:]
	return output
}

func (f *FirFilter) FilterBuffer(input, output []complex64) int {
	var samples = append(f.sampleHistory, input...)
	var length = len(samples) - f.tapsLen

	if len(output) < length {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < length; i++ {
		output[i] = DotProductResult(samples[i:], f.taps)
	}
	f.sampleHistory = samples[length:]

	return length
}

func (f *FirFilter) FilterSingle(data []complex64) complex64 {
	return DotProductResult(data, f.taps)
}

func (f *FirFilter) FilterDecimate(data []complex64, decimate int, length int) {
	var samples = append(f.sampleHistory, data...)
	var j = 0
	for i := 0; i < length; i++ {
		DotProduct(&data[i], samples[j:], f.taps)
		j += decimate
	}
	f.sampleHistory = data[len(data)-f.tapsLen:]
}

func (f *FirFilter) FilterDecimateBuffer(input, output []complex64, decimate int) int {
	var samples = append(f.sampleHistory, input...)
	var length = len(input) / decimate

	if len(output) < length {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < length; i++ {
		var srcIdx = decimate * i
		output[i] = DotProductResult(samples[srcIdx:], f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]
	return length
}

func (f *FirFilter) FilterDecimateOut(data []complex64, decimate int) []complex64 {
	var samples = append(f.sampleHistory, data...)
	var length = len(data) / decimate
	var output = make([]complex64, length)
	for i := 0; i < length; i++ {
		var srcIdx = decimate * i
		var sl = samples[srcIdx:]
		if len(sl) < len(f.taps) {
			break
		}
		output[i] = DotProductResult(sl, f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]
	return output
}

func (f *FirFilter) SetTaps(taps []float32) {
	f.taps = taps
}

func (f *FirFilter) Work(data []complex64) []complex64 {
	if f.decimation > 1 {
		return f.FilterDecimateOut(data, f.decimation)
	}
	return f.FilterOut(data)
}

func (f *FirFilter) WorkBuffer(input, output []complex64) int {
	if f.decimation > 1 {
		return f.FilterDecimateBuffer(input, output, f.decimation)
	}
	return f.FilterBuffer(input, output)
}

func (f *FirFilter) PredictOutputSize(inputLength int) int {
	return inputLength / f.decimation
}

// endregion
// region Float Fir Filter

type FloatFirFilter struct {
	taps          []float32
	sampleHistory []float32
	tapsLen       int
	decimation    int
}

func MakeFloatFirFilter(taps []float32) *FloatFirFilter {
	return &FloatFirFilter{
		taps:          taps,
		sampleHistory: make([]float32, len(taps)),
		tapsLen:       len(taps),
	}
}

func MakeDecimationFloatFirFilter(decimation int, taps []float32) *FloatFirFilter {
	return &FloatFirFilter{
		taps:          taps,
		sampleHistory: make([]float32, len(taps)),
		tapsLen:       len(taps),
		decimation:    decimation,
	}
}

func (f *FloatFirFilter) Filter(data []float32, length int) {
	var samples = append(f.sampleHistory, data...)
	for i := 0; i < length; i++ {
		DotProductFloat(&data[i], samples[i:], f.taps)
	}
	f.sampleHistory = data[len(data)-f.tapsLen:]
}

func (f *FloatFirFilter) FilterBuffer(input, output []float32) int {
	var samples = append(f.sampleHistory, input...)
	var length = len(samples) - f.tapsLen

	if len(output) < length {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < length; i++ {
		output[i] = DotProductFloatResult(samples[i:], f.taps)
	}
	f.sampleHistory = samples[length:]

	return length
}

func (f *FloatFirFilter) FilterSingle(data []float32) float32 {
	return DotProductFloatResult(data, f.taps)
}

func (f *FloatFirFilter) FilterDecimate(data []float32, decimate int, length int) {
	var samples = append(f.sampleHistory, data...)
	var j = 0
	for i := 0; i < length; i++ {
		DotProductFloat(&data[i], samples[j:], f.taps)
		j += decimate
		if j >= len(samples) {
			break
		}
	}
	f.sampleHistory = data[len(data)-f.tapsLen:]
}

func (f *FloatFirFilter) FilterDecimateOut(data []float32, decimate int) []float32 {
	var samples = append(f.sampleHistory, data...)
	var length = len(data) / decimate
	var output = make([]float32, length)
	for i := 0; i < length; i++ {
		var srcIdx = decimate * i
		output[i] = DotProductFloatResult(samples[srcIdx:], f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]
	return output
}

func (f *FloatFirFilter) FilterDecimateBuffer(input, output []float32, decimate int) int {
	var samples = append(f.sampleHistory, input...)
	var length = len(input) / decimate

	if len(output) < length {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < length; i++ {
		var srcIdx = decimate * i
		output[i] = DotProductFloatResult(samples[srcIdx:], f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]

	return length
}

func (f *FloatFirFilter) FilterOut(data []float32) []float32 {
	var samples = append(f.sampleHistory, data...)
	var length = len(data)
	var output = make([]float32, length)
	for i := 0; i < length; i++ {
		output[i] = DotProductFloatResult(samples[i:], f.taps)
	}
	f.sampleHistory = samples[len(samples)-f.tapsLen:]
	return output
}

func (f *FloatFirFilter) Work(data []float32) []float32 {
	return f.FilterOut(data)
}

func (f *FloatFirFilter) WorkBuffer(input, output []float32) int {
	return f.FilterBuffer(input, output)
}

func (f *FloatFirFilter) SetTaps(taps []float32) {
	f.taps = taps
}

func (f *FloatFirFilter) PredictOutputSize(inputLength int) int {
	return inputLength / f.decimation
}

// endregion
