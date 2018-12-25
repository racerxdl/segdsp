package demodcore

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/eventmanager"
	"math"
)

type AMDemod struct {
	sampleRate   float64
	outputRate   uint32
	firstStage   *dsp.FirFilter
	signalBw     float64
	decimation   int
	resampler    *dsp.FloatResampler
	finalStage   *dsp.FloatFirFilter
	outFifo      *fifo.Queue
	sql          *dsp.Squelch
	packedParams AMDemodParams
	ev           *eventmanager.EventManager
	lastSquelch  bool
	ffAgc        *dsp.FeedForwardAGC
	c2m          *dsp.Complex2Magnitude
}

type AMDemodParams struct {
	SampleRate      uint32
	SignalBandwidth float64
	OutputRate      uint32
	Squelch         float32
	SquelchAlpha    float32
	AudioCut        float32
}

func MakeCustomAMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32, audioCut, squelch, squelchAlpha float32) *AMDemod {
	var decim = int(math.Floor(float64(sampleRate) / signalBw))

	if decim&1 == 1 {
		decim -= 1
	}

	if decim < 1 {
		decim = 1
	}

	var quadRate = float64(sampleRate) / float64(decim)
	var resampleRate = float32(float64(outputRate) / quadRate)

	//log.Println("Decimation:", decim)
	//log.Println("Quad Rate:", quadRate)
	//log.Println("Resample Rate:", resampleRate)

	var sql = dsp.MakeSquelch(squelch, squelchAlpha)
	var agc = dsp.MakeFeedForwardAGC(1024, 1)

	return &AMDemod{
		sampleRate: float64(sampleRate),
		firstStage: dsp.MakeFirFilter(
			dsp.MakeLowPassFixed(
				1,
				float64(sampleRate),
				signalBw/2,
				127,
			),
		),
		outputRate: outputRate,
		decimation: decim,
		resampler:  dsp.MakeFloatResampler(32, resampleRate),
		finalStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPassFixed(
				1,
				float64(outputRate),
				float64(audioCut),
				31,
			),
		),
		outFifo: fifo.NewQueue(),
		sql:     sql,
		packedParams: AMDemodParams{
			SampleRate:      sampleRate,
			SignalBandwidth: signalBw,
			OutputRate:      outputRate,
			Squelch:         squelch,
			SquelchAlpha:    squelchAlpha,
			AudioCut:        audioCut,
		},
		lastSquelch: true,
		ffAgc:       agc,
		signalBw:    signalBw,
		c2m:         dsp.MakeComplex2Magnitude(),
	}
}

func (f *AMDemod) GetDemodParams() interface{} {
	return f.packedParams
}

func (f *AMDemod) SetEventManager(ev *eventmanager.EventManager) {
	f.ev = ev
}

func (f *AMDemod) IsMuted() bool {
	return f.sql.IsMuted()
}

func (f *AMDemod) GetLevel() float32 {
	return f.sql.GetAvgLevel()
}

func (f *AMDemod) Work(data []complex64) interface{} {
	var filteredData = f.firstStage.FilterDecimateOut(data, f.decimation)
	filteredData = f.sql.Work(filteredData)
	filteredData = f.ffAgc.Work(filteredData)

	var amDemodData = f.c2m.Work(filteredData)

	amDemodData = f.resampler.Work(amDemodData)

	amDemodData = f.finalStage.FilterOut(amDemodData)

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

	for i := 0; i < len(amDemodData); i++ {
		f.outFifo.Add(amDemodData[i] - 1)
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
