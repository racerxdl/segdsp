package dsp

import (
	"math"
	"math/cmplx"
)

type FrequencyTranslator struct {
	filter *CTFirFilter
	baseTaps []complex64
	centerFrequency float32
	sampleRate float32
	decimation int
	rotator *Rotator
}

func MakeFrequencyTranslatorComplexTaps(decimation int, centerFrequency, sampleRate float32, taps []complex64) *FrequencyTranslator {
	return &FrequencyTranslator{
		baseTaps: taps,
		sampleRate: sampleRate,
		centerFrequency: centerFrequency,
		decimation: decimation,
	}
}

func MakeFrequencyTranslator(decimation int, centerFrequency, sampleRate float32, taps []float32) *FrequencyTranslator {
	var baseTaps = make([]complex64, len(taps))

	for i := 0; i < len(taps); i++ {
		baseTaps[i] = complex(taps[i], 0)
	}

	return &FrequencyTranslator{
		baseTaps: baseTaps,
		sampleRate: sampleRate,
		centerFrequency: centerFrequency,
		decimation: decimation,
		rotator: MakeRotator(),
	}
}

func (ft *FrequencyTranslator) updateFilter() {
	var newTaps = make([]complex64, len(ft.baseTaps))
	var shift = 2 * math.Pi * (ft.centerFrequency / ft.sampleRate)

	for i := 0; i < len(newTaps); i++ {
		newTaps[i] = complex64(complex128(ft.baseTaps[i]) * cmplx.Exp(complex(0, float64(i) * float64(shift))))
	}

	ft.rotator.SetPhaseIncrement(cmplx.Exp(complex(0, float64(-shift * float32(ft.decimation)))))

	ft.filter = MakeCTFirFilter(newTaps)
}

func (ft *FrequencyTranslator) Work(data []complex64) []complex64 {
	var out = ft.filter.FilterDecimateOut(data, ft.decimation)
	out = ft.rotator.Work(out)

	return out
}