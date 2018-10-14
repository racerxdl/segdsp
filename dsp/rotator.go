package dsp

type Rotator struct {
	phaseIncrement complex64
	counter int
	lastPhase complex64
}

func MakeRotator() *Rotator {
	return &Rotator{
		counter: 0,
		lastPhase: 0,
	}
}

func (r *Rotator) SetPhase(p complex64) {
	r.lastPhase = p
}

func (r *Rotator) SetPhaseIncrement(increment complex64) {
	r.phaseIncrement = complex(real(increment) / ComplexAbs(increment), imag(increment) / ComplexAbs(increment))
}

func (r *Rotator) rotate(d complex64) complex64 {
	r.counter++

	var z = d * r.lastPhase

	r.lastPhase = r.lastPhase * r.phaseIncrement

	if r.counter % 512 == 0 {
		r.lastPhase = complex(real(r.lastPhase) / ComplexAbs(r.lastPhase), imag(r.lastPhase) / ComplexAbs(r.lastPhase))
	}

	return z
}

func (r *Rotator) Work(data []complex64) []complex64 {
	var out = make([]complex64, len(data))

	for i := 0; i < len(data); i++ {
		out[i] = r.rotate(data[i])
	}

	return out
}