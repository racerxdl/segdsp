package dsp

type FirFilter struct {
	taps []float32
	sampleHistory []complex64
	tapsLen int
}

func MakeFirFilter(taps []float32) *FirFilter {
	return &FirFilter{
		taps: taps,
		sampleHistory: make([]complex64, len(taps)),
		tapsLen: len(taps),
	}
}

func (f *FirFilter) Filter(data []complex64, length int) {
	var samples = append(f.sampleHistory, data...)
	for i := 0; i < length; i++ {
		DotProduct(&data[i], samples[i:i+f.tapsLen], f.taps, length)
	}
	f.sampleHistory = data[len(data) - f.tapsLen:]
}

func (f *FirFilter) FilterDecimate(data []complex64, decimate int, length int) {
	var samples = append(f.sampleHistory, data...)
	var j = 0
	for i := 0; i < length; i++ {
		DotProduct(&data[i], samples[j:j+f.tapsLen], f.taps, length)
		j += decimate
	}
	f.sampleHistory = data[len(data) - f.tapsLen:]
}
