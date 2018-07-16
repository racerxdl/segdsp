package demodcore

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
)

type WBFMDemod struct {
	sampleRate float64
	firstStage *dsp.FirFilter
	signalBw float64
	deviation float64
}

func MakeWBFMDemodulator(sampleRate uint32, signalBw float64) *WBFMDemod {
	return &WBFMDemod{
		sampleRate: float64(sampleRate),
		firstStage: dsp.MakeFirFilter(dsp.MakeLowPass(1, float64(sampleRate), signalBw, signalBw / 100)),
		deviation: 75000,
	}
}

func (f *WBFMDemod) Work(data []complex64) interface{} {
	var decim = math.Floor(f.sampleRate / f.signalBw)
	decim /= 4

	if decim < 1 {
		decim = 1
	}

	if decim == 1 {
		f.firstStage.Filter(data, len(data))
	} else {
		f.firstStage.FilterDecimate(data, int(decim), len(data))
	}
	return nil
}