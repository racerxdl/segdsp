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

// WorkBuffer Resamples the input. The output is in buffA, buffB is used as transient state for Interpolator output
// Returns number of samples in buffA
func (f *RationalResampler) WorkBuffer(buffA, buffB []complex64) int {
	l := f.interpolator.WorkBuffer(buffA, buffB)
	return f.decimator.WorkBuffer(buffB[:l], buffA)
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

// WorkBuffer Resamples the input. The output is in buffA, buffB is used as transient state for Interpolator output
// Returns number of samples in buffA
func (f *FloatRationalResampler) WorkBuffer(buffA, buffB []float32) int {
	l := f.interpolator.WorkBuffer(buffA, buffB)
	return f.decimator.WorkBuffer(buffB[:l], buffA)
}
