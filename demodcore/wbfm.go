package demodcore

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
	"strings"
	"fmt"
	"github.com/racerxdl/go.fifo"
	"os"
)

type JsonFloat32 []float32

func (u JsonFloat32) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%f", u)), ",")
	}
	return []byte(result), nil
}

type WBFMDemodData struct {
	OutputRate uint32
	Data JsonFloat32
}

type WBFMDemod struct {
	sampleRate float64
	outputRate uint32
	secondStage *dsp.FloatFirFilter
	signalBw float64
	deviation float64
	quadDemod *dsp.QuadDemod
	decimation int
	resampler *dsp.FloatResampler
	deemph *dsp.FMDeemph
	outFifo *fifo.Queue
	fileOut *os.File
}

func MakeWBFMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32) *WBFMDemod {
	var decim = math.Floor(float64(sampleRate) / signalBw)
	if (float64(sampleRate) / decim) <= float64(outputRate) {
		decim /= 4
	}

	if decim < 1 {
		decim = 1
	}

	decim = math.Floor(decim)

	var quadRate = sampleRate
	var maxDeviation = 75000.0

	var fmDemodGain = float64(quadRate) / ( 2 * math.Pi * maxDeviation )
	var resampleRate = float32(float64(outputRate) / (float64(quadRate) / decim))

	var stageCut = math.Min(float64(outputRate), float64(quadRate) / float64(decim)) / 2

	return &WBFMDemod{
		sampleRate: float64(sampleRate),
		secondStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPassFixed(
				float64(quadRate),
				stageCut,
				63,
			),
		),
		deviation: maxDeviation,
		quadDemod: dsp.MakeQuadDemod(float32(fmDemodGain)),
		decimation: int(decim),
		resampler: dsp.MakeFloatResampler(32, resampleRate),
		deemph: dsp.MakeFMDeemph(75e-6, float32(outputRate)),
		outputRate: outputRate,
		outFifo: fifo.NewQueue(),
	}
}

func (f *WBFMDemod) Work(data []complex64) interface{} {
	var fmDemodData = f.quadDemod.Work(data)

	fmDemodData = f.secondStage.FilterDecimateOut(fmDemodData, f.decimation)
	fmDemodData = f.resampler.Work(fmDemodData, len(fmDemodData))
	fmDemodData = f.deemph.Work(fmDemodData)


	for i := 0; i < len(fmDemodData); i++ {
		f.outFifo.Add(fmDemodData[i])
	}

	if f.outFifo.Len() >= 16384 {
		var outBuff = make([]float32, 16384)

		for i := 0; i < 16384; i++ {
			outBuff[i] = f.outFifo.Next().(float32)
		}

		return WBFMDemodData{
			OutputRate: f.outputRate,
			Data: outBuff,
		}
	}

	return nil
}