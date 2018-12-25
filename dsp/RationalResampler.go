package dsp

type RationalResampler struct {
	decimator    *Decimator
	interpolator *Interpolator
}

func MakeRationalResampler(interpolationRatio, decimationRatio int) *RationalResampler {
	return &RationalResampler{
		decimator:    MakeDecimator(decimationRatio),
		interpolator: MakeInterpolator(interpolationRatio),
	}
}

func (f *RationalResampler) Work(data []complex64) []complex64 {
	var interpolated = f.interpolator.Work(data)
	return f.decimator.Work(interpolated)
}

func (f *RationalResampler) WorkBuffer(input, output []complex64) int {
	tmp := f.interpolator.Work(input)
	return f.decimator.WorkBuffer(tmp, output)
}

func (f *RationalResampler) PredictOutputSize(inputLength int) int {
	return f.decimator.PredictOutputSize(f.interpolator.PredictOutputSize(inputLength))
}

type FloatRationalResampler struct {
	decimator    *FloatDecimator
	interpolator *FloatInterpolator
}

func MakeFloatRationalResampler(interpolationRatio, decimationRatio int) *FloatRationalResampler {
	return &FloatRationalResampler{
		decimator:    MakeFloatDecimator(decimationRatio),
		interpolator: MakeFloatInterpolator(interpolationRatio),
	}
}

func (f *FloatRationalResampler) Work(data []float32) []float32 {
	var interpolated = f.interpolator.Work(data)
	return f.decimator.Work(interpolated)
}

func (f *FloatRationalResampler) WorkBuffer(input, output []float32) int {
	tmp := f.interpolator.Work(input)
	return f.decimator.WorkBuffer(tmp, output)
}

func (f *FloatRationalResampler) PredictOutputSize(inputLength int) int {
	return f.decimator.PredictOutputSize(f.interpolator.PredictOutputSize(inputLength))
}
