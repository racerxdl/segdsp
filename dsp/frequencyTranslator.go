package dsp

import (
	"math"
	"math/cmplx"
)

type FrequencyTranslator struct {
	filter          *FirFilter
	baseTaps        []float32
	centerFrequency float32
	sampleRate      float32
	decimation      int
	rotator         *Rotator
	needsUpdate     bool
}

func MakeFrequencyTranslator(decimation int, centerFrequency, sampleRate float32, taps []float32) *FrequencyTranslator {
	var ft = FrequencyTranslator{
		baseTaps:        taps,
		sampleRate:      sampleRate,
		centerFrequency: centerFrequency,
		decimation:      decimation,
		rotator:         MakeRotator(),
		needsUpdate:     true,
	}

	ft.updateFilter()

	return &ft
}

func (ft *FrequencyTranslator) updateFilter() {
	var shift = float64(2 * math.Pi * (ft.centerFrequency / ft.sampleRate))

	ft.rotator.SetPhaseIncrement(complex64(cmplx.Exp(complex(0, -shift*float64(ft.decimation)))))

	ft.filter = MakeFirFilter(ft.baseTaps)
	ft.needsUpdate = false
}

func (ft *FrequencyTranslator) Work(data []complex64) []complex64 {
	if ft.needsUpdate {
		ft.updateFilter()
	}

	var out = ft.rotator.Work(data)
	if ft.decimation != 1 {
		out = ft.filter.FilterOut(out)
	} else {
		out = ft.filter.FilterDecimateOut(out, ft.decimation)
	}

	return out
}

func (ft *FrequencyTranslator) SetFrequency(frequency float32) {
	ft.centerFrequency = frequency
	ft.needsUpdate = true
}

func (ft *FrequencyTranslator) SetDecimation(decimation int) {
	ft.decimation = decimation
	ft.needsUpdate = true
}

func (ft *FrequencyTranslator) GetDecimation() int {
	return ft.decimation
}

func (ft *FrequencyTranslator) GetFrequency() float32 {
	return ft.centerFrequency
}
