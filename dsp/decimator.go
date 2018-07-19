package dsp


type Decimator struct {
	fir *FirFilter
	decimationRatio int
}

func MakeDecimator(decimationRatio int) *Decimator {
	return &Decimator{
		fir: MakeFirFilter(MakeLowPassFixed(1, 1 / (2 * float64(decimationRatio)), 63)),
		decimationRatio: decimationRatio,
	}
}

func (f *Decimator) Work(data []complex64) []complex64 {
	return f.fir.FilterDecimateOut(data, f.decimationRatio)
}


type FloatDecimator struct {
	fir *FloatFirFilter
	decimationRatio int
}

func MakeFloatDecimator(decimationRatio int) *FloatDecimator {
	return &FloatDecimator{
		fir: MakeFloatFirFilter(MakeLowPassFixed(1, 1 / (2 * float64(decimationRatio)), 63)),
		decimationRatio: decimationRatio,
	}
}

func (f *FloatDecimator) Work(data []float32) []float32 {
	return f.fir.FilterDecimateOut(data, f.decimationRatio)
}