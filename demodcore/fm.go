package demodcore

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
	"github.com/racerxdl/go.fifo"
)

type FMDemod struct {
	sampleRate float64
	outputRate uint32
	firstStage *dsp.FirFilter
	secondStage *dsp.FloatFirFilter
	signalBw float64
	deviation float32
	quadDemod *dsp.QuadDemod
	decimation int
	resampler *dsp.FloatResampler
	finalStage *dsp.FloatFirFilter
	deemph *dsp.FMDeemph
	outFifo *fifo.Queue
	sql *dsp.Squelch
	tau float32
}

func MakeCustomFMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32, tau, squelch, squelchAlpha, maxDeviation float32) *FMDemod {
	var decim = math.Floor(float64(sampleRate) / signalBw)
	if (float64(sampleRate) / decim) <= float64(outputRate) {
		decim /= 4
	}

	if decim < 1 {
		decim = 1
	}

	decim = math.Floor(decim)

	var quadRate = sampleRate

	var fmDemodGain = float64(quadRate) / ( 2 * math.Pi * float64(maxDeviation) )
	var intermediateRate = float64(quadRate) / decim
	var resampleRate = float32(float64(outputRate) / intermediateRate)

	var stageCut = math.Min(float64(outputRate), intermediateRate) / 2

	return &FMDemod{
		sampleRate: float64(sampleRate),
		firstStage: dsp.MakeFirFilter(
			dsp.MakeLowPassFixed(
				1,
				float64(sampleRate),
				signalBw,
				63,
			),
		),
		secondStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPassFixed(
				1,
				float64(quadRate),
				stageCut,
				63,
			),
		),
		tau: tau,
		deviation: maxDeviation,
		quadDemod: dsp.MakeQuadDemod(float32(fmDemodGain)),
		decimation: int(decim),
		resampler: dsp.MakeFloatResampler(32, resampleRate),
		deemph: dsp.MakeFMDeemph(tau, float32(outputRate)),
		finalStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPassFixed(
				0.25,
				float64(outputRate),
				float64(outputRate) / 2 - float64(outputRate) / 32,
				63,
			),
		),
		outputRate: outputRate,
		outFifo: fifo.NewQueue(),
		sql: dsp.MakeSquelch(squelch, squelchAlpha),
	}
}

func MakeWBFMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32) *FMDemod {
	return MakeCustomFMDemodulator(sampleRate, signalBw, outputRate, 75e-6, 75000, -150, 0.001)
}

func (f *FMDemod) Work(data []complex64) interface{} {
	var filteredData = f.firstStage.FilterOut(data)
	filteredData = f.sql.Work(data)

	var fmDemodData = f.quadDemod.Work(filteredData)

	fmDemodData = f.secondStage.FilterDecimateOut(fmDemodData, f.decimation)
	fmDemodData = f.resampler.Work(fmDemodData)
	if f.tau != 0 {
		fmDemodData = f.deemph.Work(fmDemodData)
	}
	fmDemodData = f.finalStage.FilterOut(fmDemodData)


	for i := 0; i < len(fmDemodData); i++ {
		f.outFifo.Add(fmDemodData[i])
	}

	if f.outFifo.Len() >= 16384 {
		var outBuff = make([]float32, 16384)

		for i := 0; i < 16384; i++ {
			outBuff[i] = f.outFifo.Next().(float32)
		}

		return DemodData{
			OutputRate: f.outputRate,
			Data: outBuff,
		}
	}

	return nil
}