package digital

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
)

type MMSEFirInterpolator struct {
	filters []dsp.FirFilter
}

func MakeMMSEFirInterpolator() *MMSEFirInterpolator {
	var filters = make([]dsp.FirFilter, iNSTEPS+1)
	for i := 0; i < iNSTEPS; i++ {
		var t = make([]float32, len(interpTaps[i]))
		copy(t, interpTaps[i])
		filters[i] = *dsp.MakeFirFilter(t)
	}

	return &MMSEFirInterpolator{
		filters: filters,
	}
}

func (mmse *MMSEFirInterpolator) GetNTaps() int {
	return iNTAPS
}

func (mmse *MMSEFirInterpolator) GetNSteps() int {
	return iNSTEPS
}

func (mmse *MMSEFirInterpolator) Interpolate(input []complex64, mu float32) complex64 {
	var imu = int(math.Round(float64(mu * iNSTEPS)))

	if imu < 0 || imu > iNSTEPS {
		panic("MMSE Fir Interpolator tried to filter with a unknown division value")
	}

	return mmse.filters[imu].FilterSingle(input)
}
