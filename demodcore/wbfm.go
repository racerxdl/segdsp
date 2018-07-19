package demodcore

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
	"strings"
	"fmt"
	"github.com/racerxdl/go.fifo"
	"os"
	"bytes"
	"encoding/binary"
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
	rresampler *dsp.FloatRationalResampler
	deemph *dsp.FMDeemph
	outFifo *fifo.Queue
	fileOut *os.File
}

func MakeWBFMDemodulator(sampleRate uint32, signalBw float64, outputRate uint32) *WBFMDemod {
	var decim = math.Floor(float64(sampleRate) / signalBw)
	//if (float64(sampleRate) / decim) <= float64(outputRate) {
	//	decim /= 4
	//}
	//
	//if decim < 1 {
	//	decim = 1
	//}
	//
	//decim = math.Floor(decim)

	decim = 20

	var quadRate = sampleRate
	var maxDeviation = 75000.0

	var fmDemodGain = float64(quadRate) / ( 2 * math.Pi * maxDeviation )

	fmt.Println("Demod Gain: ", fmDemodGain)

	var transitionBand = outputRate / 32
	var resampleRate = float32(float64(outputRate) / (float64(quadRate) / decim))

	file, _ := os.Create("test.raw")

	return &WBFMDemod{
		sampleRate: float64(sampleRate),
		secondStage: dsp.MakeFloatFirFilter(
			dsp.MakeLowPass(
				1,
				float64(quadRate),
				float64(outputRate / 2 - transitionBand),
				float64(transitionBand),
			),
		),
		//secondStage: dsp.MakeFloatFirFilter(
		//	dsp.MakeLowPassFixed(
		//		float64(quadRate),
		//		float64(float64(quadRate) / (2 * float64(decim))),
		//		63,
		//	),
		//),
		deviation: maxDeviation,
		quadDemod: dsp.MakeQuadDemod(float32(fmDemodGain)),
		decimation: int(decim),
		resampler: dsp.MakeFloatResampler(16, resampleRate),
		rresampler: dsp.MakeFloatRationalResampler(32, 25),
		deemph: dsp.MakeFMDeemph(75e-6, float32(outputRate)),
		outputRate: outputRate,
		outFifo: fifo.NewQueue(),
		fileOut: file,
	}
}

func (f *WBFMDemod) Work(data []complex64) interface{} {
	var fmDemodData = f.quadDemod.Work(data)
	fmDemodData = f.secondStage.FilterDecimateOut(fmDemodData, f.decimation)

	//fmDemodData = f.rresampler.Work(fmDemodData)
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

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, outBuff)

		f.fileOut.Write(buf.Bytes())
		go f.fileOut.Sync()

		return WBFMDemodData{
			OutputRate: f.outputRate,
			Data: outBuff,
		}
	}

	return nil
}