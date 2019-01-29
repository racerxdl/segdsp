package dsp

import (
	"fmt"
	"math"
)

const TwoPi = float32(math.Pi * 2)
const MinusTwoPi = -TwoPi
const OneOverTwoPi = float32(1 / (2 * math.Pi))

type controlLoop struct {
	phase         float32
	freq          float32
	maxRelFreq    float32
	minRelFreq    float32
	dampingFactor float32
	loopBandwidth float32
	alpha         float32
	beta          float32
}

func makeControlLoop(loopBw, minFreq, maxFreq float32) *controlLoop {
	cl := &controlLoop{
		dampingFactor: float32(math.Sqrt(2) / 2),
		phase:         0,
		freq:          0,
		maxRelFreq:    maxFreq,
		minRelFreq:    minFreq,
	}

	_ = cl.SetLoopBandwidth(loopBw)
	return cl
}

func (cl *controlLoop) SetLoopBandwidth(bw float32) error {
	if bw < 0 {
		return fmt.Errorf("bandwidth should be higher or equal to 0. Got %f", bw)
	}

	cl.loopBandwidth = bw
	cl.UpdateGains()

	return nil
}

func (cl *controlLoop) SetDampingFactor(df float32) error {
	if df <= 0 {
		return fmt.Errorf("damping factor should be higher than 0. Got %f", df)
	}

	cl.dampingFactor = df
	cl.UpdateGains()

	return nil
}

func (cl *controlLoop) SetAlpha(alpha float32) error {
	if alpha < 0 || alpha > 1 {
		return fmt.Errorf("alpha needs to be between 0 and 1. Got %f", alpha)
	}

	cl.alpha = alpha

	return nil
}

func (cl *controlLoop) SetBeta(beta float32) error {
	if beta < 0 || beta > 1 {
		return fmt.Errorf("beta needs to be between 0 and 1. Got %f", beta)
	}

	cl.beta = beta

	return nil
}

func (cl *controlLoop) SetFrequency(freq float32) {
	if freq > cl.maxRelFreq {
		cl.freq = cl.minRelFreq
	} else if freq < cl.minRelFreq {
		cl.freq = cl.maxRelFreq
	} else {
		cl.freq = freq
	}
}

func (cl *controlLoop) Reset() {
	cl.phase = 0
	cl.freq = 0
	cl.dampingFactor = float32(math.Sqrt(2.0) / 2.0)
	_ = cl.SetLoopBandwidth(cl.loopBandwidth)
}

func (cl *controlLoop) UpdateGains() {
	denom := 1.0 + 2.0*cl.dampingFactor*cl.loopBandwidth + cl.loopBandwidth*cl.loopBandwidth
	cl.alpha = (4 * cl.dampingFactor * cl.loopBandwidth) / denom
	cl.beta = (4 * cl.loopBandwidth * cl.loopBandwidth) / denom
}

func (cl *controlLoop) SetPhase(phase float32) {
	cl.phase = phase
}

func (cl *controlLoop) phaseWrap() {
	if cl.phase > TwoPi || cl.phase < MinusTwoPi {
		cl.phase = cl.phase*OneOverTwoPi - float32(int(cl.phase*OneOverTwoPi))
		cl.phase = cl.phase * TwoPi
	}
}

func (cl *controlLoop) AdvanceLoop(err float32) {
	cl.freq = cl.beta*err + cl.freq
	cl.phase = cl.phase + cl.alpha*err + cl.freq
}

func (cl *controlLoop) SetRelativeMaxFrequency(freq float32) {
	cl.maxRelFreq = freq
}

func (cl *controlLoop) SetRelativeMinFrequency(freq float32) {
	cl.minRelFreq = freq
}

func (cl *controlLoop) GetLoopBandwidth() float32 {
	return cl.loopBandwidth
}

func (cl *controlLoop) GetDampingFactor() float32 {
	return cl.dampingFactor
}

func (cl *controlLoop) GetAlpha() float32 {
	return cl.alpha
}

func (cl *controlLoop) GetBeta() float32 {
	return cl.beta
}

func (cl *controlLoop) GetFrequency() float32 {
	return cl.freq
}

func (cl *controlLoop) GetFrequencyHz() float32 {
	return cl.freq / TwoPi
}

func (cl *controlLoop) GetPhase() float32 {
	return cl.phase
}

func (cl *controlLoop) GetMaxRelativeFrequency() float32 {
	return cl.maxRelFreq
}

func (cl *controlLoop) GetMinRelativeFrequency() float32 {
	return cl.minRelFreq
}

func (cl *controlLoop) frequencyLimit() {
	if cl.freq > cl.maxRelFreq {
		cl.freq = cl.maxRelFreq
	} else if cl.freq < cl.minRelFreq {
		cl.freq = cl.minRelFreq
	}
}
