package dsp

type Decimator struct {
	fir             *FirFilter
	decimationRatio int
}

func MakeDecimator(decimationRatio int) *Decimator {
	return &Decimator{
		fir:             MakeDecimationFirFilter(decimationRatio, MakeLowPassFixed(1, 1, 1/(2*float64(decimationRatio)), 127)),
		decimationRatio: decimationRatio,
	}
}

func (f *Decimator) Work(data []complex64) []complex64 {
	return f.fir.FilterDecimateOut(data, f.decimationRatio)
}

func (f *Decimator) WorkBuffer(input, output []complex64) int {
	return f.fir.FilterDecimateBuffer(input, output, f.decimationRatio)
}

func (f *Decimator) PredictOutputSize(inputLength int) int {
	return f.fir.PredictOutputSize(inputLength)
}

type FloatDecimator struct {
	fir             *FloatFirFilter
	decimationRatio int
}

func MakeFloatDecimator(decimationRatio int) *FloatDecimator {
	return &FloatDecimator{
		fir:             MakeDecimationFloatFirFilter(decimationRatio, MakeLowPassFixed(1, 1, 1/(2*float64(decimationRatio)), 127)),
		decimationRatio: decimationRatio,
	}
}

func (f *FloatDecimator) Work(data []float32) []float32 {
	return f.fir.Work(data)
}

func (f *FloatDecimator) WorkBuffer(input, output []float32) int {
	return f.fir.WorkBuffer(input, output)
}

func (f *FloatDecimator) PredictOutputSize(inputLength int) int {
	return f.fir.PredictOutputSize(inputLength)
}
