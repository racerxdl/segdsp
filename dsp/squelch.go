package dsp

import (
	"math"
)

type Squelch struct {
	threshold    float32
	thresholddB  float32
	muted        bool
	filter       *SinglePoleIIRFilter
	avgThreshold float32
}

func MakeSquelch(threshold, alpha float32) *Squelch {
	var s = Squelch{
		filter: MakeSinglePoleIIRFilter(alpha),
	}
	s.SetThreshold(threshold)
	return &s
}

func (f *Squelch) GetAvgLevel() float32 {
	return float32(10 * math.Log10(float64(f.avgThreshold)))
}

func (f *Squelch) IsMuted() bool {
	return f.muted
}

func (f *Squelch) SetAlpha(alpha float32) {
	f.filter.SetTaps(alpha)
}

func (f *Squelch) SetThreshold(dB float32) {
	f.thresholddB = dB
	f.threshold = float32(math.Pow(10, float64(dB/10.0)))
}

func (f *Squelch) GetThreshold() float32 {
	return f.thresholddB
}

func (f *Squelch) Work(data []complex64) []complex64 {
	var out = make([]complex64, len(data))

	var avg = float32(0)
	for i := 0; i < len(data); i++ {
		v := data[i]
		mag := real(v)*real(v) + imag(v)*imag(v)
		v2 := f.filter.Filter(mag)
		avg += v2
		out[i] = complex(0, 0)
	}
	avg /= float32(len(data))
	f.avgThreshold = avg
	f.muted = avg <= f.threshold

	if avg >= f.threshold {
		return data
	} else {
		return out
	}
}

func (f *Squelch) WorkBuffer(input, output []complex64) int {
	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	var avg = float32(0)
	for i := 0; i < len(input); i++ {
		v := input[i]
		mag := real(v)*real(v) + imag(v)*imag(v)
		v2 := f.filter.Filter(mag)
		avg += v2
		output[i] = complex(0, 0)
	}
	avg /= float32(len(input))
	f.avgThreshold = avg
	f.muted = avg <= f.threshold

	if avg >= f.threshold {
		copy(output, input)
	}

	return len(input)
}

func (f *Squelch) PredictOutputSize(inputLength int) int {
	return inputLength
}
