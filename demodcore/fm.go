package demodcore

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/eventmanager"
	"log"
	"math"
)

type FMDemod struct {
	sampleRate   float64
	outputRate   uint32
	firstStage   *dsp.FirFilter
	secondStage  *dsp.FloatFirFilter
	signalBw     float64
	deviation    float32
	quadDemod    *dsp.QuadDemod
	decimation   int
	resampler    *dsp.FloatResampler
	finalStage   *dsp.FloatFirFilter
	deemph       *dsp.FMDeemph
	outFifo      *fifo.Queue
	sql          *dsp.Squelch
	tau          float32
	packedParams FMDemodParams
	ev           *eventmanager.EventManager
	lastSquelch  bool
}

type FMDemodParams struct {
	SampleRate      uint32
	SignalBandwidth float64
	OutputRate      uint32
	Tau             float32
	Squelch         float32
	SquelchAlpha    float32
	MaxDeviation    float32
}

func MakeCustomFMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32, tau, squelch, squelchAlpha, maxDeviation float32) *FMDemod {
	var decim = int(math.Floor(float64(sampleRate) / signalBw / 2))

	if decim&1 == 1 {
		decim -= 1
	}

	if decim < 1 {
		decim = 1
	}

	var quadRate = float64(sampleRate) / float64(decim)

	log.Println("Decimation:", decim)
	log.Println("Quad Rate:", quadRate)

	var fmDemodGain = quadRate / (2 * math.Pi * float64(maxDeviation))
	var resampleRate = float32(float64(outputRate) / quadRate)

	var stageCut = math.Min(float64(outputRate), quadRate) / 2

	var sql = dsp.MakeSquelch(squelch, squelchAlpha)

	return &FMDemod{
		sampleRate: float64(sampleRate),
		firstStage: dsp.MakeFirFilter(
			dsp.MakeLowPassFixed(
				1,
				float64(sampleRate),
				signalBw/2,
				63,
			),
		),
		secondStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPassFixed(
				1,
				quadRate,
				stageCut,
				63,
			),
		),
		tau:        tau,
		deviation:  maxDeviation,
		quadDemod:  dsp.MakeQuadDemod(float32(fmDemodGain)),
		decimation: int(decim),
		resampler:  dsp.MakeFloatResampler(32, resampleRate),
		deemph:     dsp.MakeFMDeemph(tau, float32(outputRate)),
		finalStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPassFixed(
				0.25,
				float64(outputRate),
				float64(outputRate)/2-float64(outputRate)/32,
				63,
			),
		),
		outputRate: outputRate,
		outFifo:    fifo.NewQueue(),
		sql:        sql,
		packedParams: FMDemodParams{
			SampleRate:      sampleRate,
			SignalBandwidth: signalBw,
			OutputRate:      outputRate,
			Tau:             tau,
			Squelch:         squelch,
			SquelchAlpha:    squelchAlpha,
			MaxDeviation:    maxDeviation,
		},
		lastSquelch: true,
		signalBw:    signalBw,
	}
}

func MakeWBFMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32) *FMDemod {
	return MakeCustomFMDemodulator(sampleRate, signalBw, outputRate, 75e-6, -150, 0.01, 75000)
}

func (f *FMDemod) GetLevel() float32 {
	return f.sql.GetAvgLevel()
}

func (f *FMDemod) GetDemodParams() interface{} {
	return f.packedParams
}

func (f *FMDemod) SetEventManager(ev *eventmanager.EventManager) {
	f.ev = ev
}

func (f *FMDemod) IsMuted() bool {
	return f.sql.IsMuted()
}

func (f *FMDemod) Work(data []complex64) interface{} {
	var filteredData = f.firstStage.FilterDecimateOut(data, f.decimation)
	filteredData = f.sql.Work(filteredData)

	var fmDemodData = f.quadDemod.Work(filteredData)

	fmDemodData = f.secondStage.FilterOut(fmDemodData)
	fmDemodData = f.resampler.Work(fmDemodData)
	fmDemodData = f.finalStage.FilterOut(fmDemodData)
	if f.tau != 0 {
		fmDemodData = f.deemph.Work(fmDemodData)
	}

	if f.lastSquelch != f.sql.IsMuted() && f.ev != nil {
		var evName string
		if f.sql.IsMuted() {
			evName = eventmanager.EvSquelchOn
		} else {
			evName = eventmanager.EvSquelchOff
		}
		f.ev.Emit(evName, eventmanager.SquelchEventData{
			Threshold: f.sql.GetThreshold(),
			AvgValue:  f.sql.GetAvgLevel(),
		})
	}

	f.lastSquelch = f.sql.IsMuted()

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
			Level:      f.sql.GetAvgLevel(),
			Data:       outBuff,
		}
	}

	return nil
}
