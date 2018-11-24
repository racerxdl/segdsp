package dsp

import "math"

// NCO is a Numeric Controlled Oscillator
// Based on GNURadio Implementation
type NCO struct {
	phase          float32
	phaseIncrement float32
}

func MakeNCO() *NCO {
	return &NCO{
		phase:          0,
		phaseIncrement: 0,
	}
}

// SetPhase in Radians
func (nco *NCO) SetPhase(angle float32) {
	nco.phase = angle
}

// AdjustPhase Increments / decrements current phase. In radians
func (nco *NCO) AdjustPhase(deltaAngle float32) {
	nco.phase += deltaAngle
}

// SetFrequency Sets the Phase Increment in Radians / step
func (nco *NCO) SetFrequency(rate float32) {
	nco.phaseIncrement = rate
}

// AdjustFrequency Increments / Decrements the phase increment. In radians / step
func (nco *NCO) AdjustFrequency(deltaRate float32) {
	nco.phaseIncrement += deltaRate
}

// Step makes a single step in NCO
func (nco *NCO) Step() {
	nco.phase += nco.phaseIncrement
}

// StepN makes N steps in NCO
func (nco *NCO) StepN(n int) {
	nco.phase += nco.phaseIncrement * float32(n)
}

// GetPhase returns the current phase value in radians
func (nco *NCO) GetPhase() float32 {
	return nco.phase
}

// GetPhaseIncrement returns the phase increment value in radians / step
func (nco *NCO) GetPhaseIncrement() float32 {
	return nco.phaseIncrement
}

// Float32Sin Compute N elements for a float32 array sine wave
func (nco *NCO) Float32Sin(n int, amplitude float32) []float32 {
	var d = make([]float32, n)
	for i := 0; i < n; i++ {
		d[i] = float32(math.Sin(float64(nco.phase))) * amplitude
	}
	return d
}

// Float32Cos Compute N elements for a float32 array cosine wave
func (nco *NCO) Float32Cos(n int, amplitude float32) []float32 {
	var d = make([]float32, n)
	for i := 0; i < n; i++ {
		d[i] = float32(math.Cos(float64(nco.phase))) * amplitude
	}
	return d
}

// Float32Sin Compute N elements for a float32 array sine wave
func (nco *NCO) Complex64SinCos(n int, amplitude float32) []complex64 {
	var d = make([]complex64, n)
	for i := 0; i < n; i++ {
		a, b := math.Sincos(float64(nco.phase))
		d[i] = complex(float32(a)*amplitude, float32(b)*amplitude)
	}
	return d
}
