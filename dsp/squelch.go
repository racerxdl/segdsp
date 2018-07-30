package dsp

import (
	"math"
)

type Squelch struct {
	threshold float32
	muted bool
	filter *SinglePoleIIRFilter
}

func MakeSquelch(threshold, alpha float32) *Squelch {
	var s = Squelch{
		filter: MakeSinglePoleIIRFilter(alpha),
	}
	s.SetThreshold(threshold)
	return &s
}

func (f *Squelch) IsMuted() bool {
	return f.muted
}

func (f *Squelch) SetAlpha(alpha float32) {
	f.filter.SetTaps(alpha)
}

func (f *Squelch) SetThreshold(dB float32) {
	f.threshold = float32(math.Pow(10, float64(dB / 10.0)))
}

func (f *Squelch) Work(data []complex64) []complex64 {
	var out = make([]complex64, len(data))

	for i := 0; i < len(data); i++ {
		v := data[i]
		mag := real(v) * real(v) + imag(v) * imag(v)
		v2 := f.filter.Filter(mag)
		if v2 >= f.threshold {
			out[i] = v
		} else {
			out[i] = complex(0, 0)
		}
	}

	f.muted = f.filter.GetPreviousOutput() >= f.threshold

	return out
}
