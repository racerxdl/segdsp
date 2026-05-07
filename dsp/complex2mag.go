package dsp

import (
	"github.com/racerxdl/segdsp/tools"
)

type Complex2Magnitude struct {
	outBuf []float32
}

func MakeComplex2Magnitude() *Complex2Magnitude {
	return &Complex2Magnitude{}
}

func (cm *Complex2Magnitude) Work(data []complex64) []float32 {
	if cap(cm.outBuf) < len(data) {
		cm.outBuf = make([]float32, len(data))
	}
	cm.WorkBuffer(data, cm.outBuf[:len(data)])
	return cm.outBuf[:len(data)]
}

func (cm *Complex2Magnitude) WorkBuffer(input []complex64, output []float32) int {
	if len(input) != len(output) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		sample := input[i]
		output[i] = tools.ComplexAbs(sample)
		//output[i] = float32(math.Sqrt(float64(real(sample)*real(sample) + imag(sample)*imag(sample))))
	}

	return len(output)
}

func (cm *Complex2Magnitude) PredictOutputSize(inputLength int) int {
	return inputLength
}
