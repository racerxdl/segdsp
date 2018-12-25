package dsp

import (
	"math"
	"math/cmplx"
)

type FrequencyTranslator struct {
	filter          *CTFirFilter
	sampleHistory   []complex64
	baseTaps        []complex64
	centerFrequency float32
	sampleRate      float32
	decimation      int
	rotator         *Rotator
	needsUpdate     bool
	tapsLen         int
}

func MakeFrequencyTranslator(decimation int, centerFrequency, sampleRate float32, taps []float32) *FrequencyTranslator {
	var baseTaps = make([]complex64, len(taps))

	for i := 0; i < len(taps); i++ {
		baseTaps[i] = complex(taps[i], 0)
	}

	var ft = FrequencyTranslator{
		baseTaps:        baseTaps,
		sampleRate:      sampleRate,
		centerFrequency: centerFrequency,
		decimation:      decimation,
		rotator:         MakeRotator(),
		needsUpdate:     true,
		sampleHistory:   make([]complex64, len(taps)),
		tapsLen:         len(taps),
	}

	ft.updateFilter()

	return &ft
}

func (ft *FrequencyTranslator) updateFilter() {
	var newTaps = make([]complex64, len(ft.baseTaps))

	var fDecimation = float64(ft.decimation)
	var shift = float64(2 * math.Pi * ft.centerFrequency / ft.sampleRate)

	ft.rotator.SetPhaseIncrement(complex64(cmplx.Exp(complex(0, -shift*fDecimation))))

	for i := 0; i < len(newTaps); i++ {
		var fi = float64(i)
		newTaps[i] = complex64(complex128(ft.baseTaps[i]) * cmplx.Exp(complex(0, fi*shift)))
	}

	ft.filter = MakeDecimationCTFirFilter(ft.decimation, newTaps)

	ft.needsUpdate = false
}

func (ft *FrequencyTranslator) Work(data []complex64) []complex64 {
	if ft.needsUpdate {
		ft.updateFilter()
	}
	var outLength = len(data) / ft.decimation
	var samples = append(ft.sampleHistory, data...)
	var output []complex64
	if ft.decimation > 1 {
		output = ft.filter.FilterDecimateOut(samples, ft.decimation)
	} else {
		output = ft.filter.FilterOut(samples)
	}
	output = ft.rotator.Work(output)
	ft.sampleHistory = samples[len(samples)-ft.tapsLen:]

	return output[:outLength]
}

func (ft *FrequencyTranslator) WorkBuffer(input, output []complex64) int {
	if ft.needsUpdate {
		ft.updateFilter()
	}

	var outLength = len(input) / ft.decimation
	var samples = append(ft.sampleHistory, input...)

	if len(output) < outLength {
		panic("There is not enough space in output buffer")
	}

	ft.filter.WorkBuffer(input, output)
	ft.rotator.WorkInline(output)
	ft.sampleHistory = samples[len(samples)-ft.tapsLen:]

	return outLength
}

func (ft *FrequencyTranslator) PredictOutputSize(inputLength int) int {
	return ft.filter.PredictOutputSize(inputLength)
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
