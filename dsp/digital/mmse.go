package digital

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
)

// region Complex MMSE Fir Interpolator
type ComplexMMSEFirInterpolator struct {
	filters []dsp.FirFilter
}

func MakeComplexMMSEFirInterpolator() *ComplexMMSEFirInterpolator {
	var filters = make([]dsp.FirFilter, iNSTEPS+1)
	for i := 0; i < iNSTEPS; i++ {
		var t = make([]float32, len(interpTaps[i]))
		copy(t, interpTaps[i])
		filters[i] = *dsp.MakeFirFilter(t)
	}

	return &ComplexMMSEFirInterpolator{
		filters: filters,
	}
}

func (mmse *ComplexMMSEFirInterpolator) GetNTaps() int {
	return iNTAPS
}

func (mmse *ComplexMMSEFirInterpolator) GetNSteps() int {
	return iNSTEPS
}

func (mmse *ComplexMMSEFirInterpolator) Interpolate(input []complex64, mu float32) complex64 {
	var imu = int(math.Round(float64(mu * iNSTEPS)))

	if imu < 0 || imu > iNSTEPS {
		panic("MMSE Fir Interpolator tried to filter with a unknown division value")
	}

	return mmse.filters[imu].FilterSingle(input)
}

// endregion
// region Float MMSE Fir Interpolator

type FloatMMSEFirInterpolator struct {
	filters []dsp.FloatFirFilter
}

func MakeFloatMMSEFirInterpolator() *FloatMMSEFirInterpolator {
	var filters = make([]dsp.FloatFirFilter, iNSTEPS+1)
	for i := 0; i < iNSTEPS; i++ {
		var t = make([]float32, len(interpTaps[i]))
		copy(t, interpTaps[i])
		filters[i] = *dsp.MakeFloatFirFilter(t)
	}

	return &FloatMMSEFirInterpolator{
		filters: filters,
	}
}

func (mmse *FloatMMSEFirInterpolator) GetNTaps() int {
	return iNTAPS
}

func (mmse *FloatMMSEFirInterpolator) GetNSteps() int {
	return iNSTEPS
}

func (mmse *FloatMMSEFirInterpolator) Interpolate(input []float32, mu float32) float32 {
	var imu = int(math.Round(float64(mu * iNSTEPS)))

	if imu < 0 || imu > iNSTEPS {
		panic("MMSE Fir Interpolator tried to filter with a unknown division value")
	}

	return mmse.filters[imu].FilterSingle(input)
}

// endregion
