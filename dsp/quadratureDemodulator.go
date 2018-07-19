package dsp


type QuadDemod struct {
	gain float32
	history []complex64
}

func MakeQuadDemod(gain float32) *QuadDemod {
	return &QuadDemod{
		gain: gain,
		history: make([]complex64, 2),
	}
}

func (f *QuadDemod) Work(data []complex64) []float32 {
	var samples = append(f.history, data...)
	var tmp = MultiplyConjugate(samples[1:], samples, len(samples) - 2)
	var out = make([]float32, len(samples) - 2)

	for i := 0; i < len(out); i++ {
		out[i] = f.gain * Atan2(imag(tmp[i]), real(tmp[i]))
	}

	f.history = samples[len(samples)-2:]
	//var out = make([]float32, len(data))
	//
	//for i := 0; i < len(data); i++ {
	//	var c = data[i] * Conj(f.history[0])
	//	//var mod = Modulus(c)
	//	//if mod > 0 {
	//	//	mod = 1 / mod
	//	//	c = complex(real(c) * mod, imag(c) * mod)
	//	//}
	//
	//	//var argument = Argument(c)
	//	out[i] = Atan2(imag(c), real(c)) * f.gain
	//	f.history[0] = data[i]
	//}

	return out
}