package dsp

type Interpolator struct {
	fir                *FirFilter
	interpolationRatio int
}

func MakeInterpolator(interpolationRatio int) *Interpolator {
	return &Interpolator{
		fir:                MakeFirFilter(MakeLowPassFixed(1, 1, 1/float64(interpolationRatio*2), 63)),
		interpolationRatio: interpolationRatio,
	}
}

func (f *Interpolator) Work(data []complex64) []complex64 {
	var output = make([]complex64, len(data)*f.interpolationRatio)

	for i := 0; i < len(data); i++ {
		var idx = i * f.interpolationRatio
		output[idx] = data[i]
		for j := 1; j < f.interpolationRatio; j++ {
			output[idx+j] = complex(0, 0)
		}
	}

	f.fir.Filter(output, len(output))
	return output
}

type FloatInterpolator struct {
	fir                *FloatFirFilter
	interpolationRatio int
}

func MakeFloatInterpolator(interpolationRatio int) *FloatInterpolator {
	return &FloatInterpolator{
		fir:                MakeFloatFirFilter(MakeLowPassFixed(1, 1, 1/float64(interpolationRatio*2), 63)),
		interpolationRatio: interpolationRatio,
	}
}

func (f *FloatInterpolator) Work(data []float32) []float32 {
	var output = make([]float32, len(data)*f.interpolationRatio)

	for i := 0; i < len(data); i++ {
		var idx = i * f.interpolationRatio
		output[idx] = data[i]
		for j := 1; j < f.interpolationRatio; j++ {
			output[idx+j] = 0
		}
	}

	f.fir.Filter(output, len(output))
	return output
}
