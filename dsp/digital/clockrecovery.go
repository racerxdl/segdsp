package digital

import (
	"github.com/racerxdl/segdsp/dsp"
)

const ccHistoryLength = 3
const fudgeFactor = 16

type ComplexClockRecovery struct {
	consumed           int
	sampleHistoryCount int
	sampleHistory      []complex64

	p2T complex64
	p1T complex64
	p0T complex64

	c2T complex64
	c1T complex64
	c0T complex64

	mu                 float32
	omega              float32
	gainOmega          float32
	omegaRelativeLimit float32
	omegaMidValue      float32
	omegaLimit         float32
	gainMu             float32

	interp *MMSEFirInterpolator
}

func (ccr *ComplexClockRecovery) GetMu() float32 {
	return ccr.mu
}

func (ccr *ComplexClockRecovery) GetOmega() float32 {
	return ccr.omega
}

func (ccr *ComplexClockRecovery) GetGainMu() float32 {
	return ccr.gainMu
}

func (ccr *ComplexClockRecovery) GetGainOmega() float32 {
	return ccr.gainOmega
}

func (ccr *ComplexClockRecovery) SetGainMu(gain float32) {
	ccr.gainMu = gain
}

func (ccr *ComplexClockRecovery) SetGainOmega(gain float32) {
	ccr.gainOmega = gain
}

func (ccr *ComplexClockRecovery) SetMu(mu float32) {
	ccr.mu = mu
}

func (ccr *ComplexClockRecovery) SetOmega(omega float32) {
	ccr.omega = omega
	ccr.omegaMidValue = omega
	ccr.omegaLimit = ccr.omegaRelativeLimit * omega
}

func NewComplexClockRecovery(omega, gainOmega, mu, gainMu, omegaRelativeLimit float32) *ComplexClockRecovery {
	var ccr = &ComplexClockRecovery{
		omega:              omega,
		gainOmega:          gainOmega,
		mu:                 mu,
		gainMu:             gainMu,
		omegaRelativeLimit: omegaRelativeLimit,
		interp:             MakeMMSEFirInterpolator(),
	}

	if omega <= 0 {
		panic("Clock Recovery Rate (omega) must be higher than 0")
	}

	if gainMu < 0 || gainOmega < 0 {
		panic("Clock Recovery gains should be positive.")
	}

	ccr.SetOmega(omega)

	ccr.sampleHistoryCount = ccHistoryLength
	ccr.sampleHistory = make([]complex64, ccHistoryLength)

	return ccr
}

func slicer(sample complex64) complex64 {
	var r = 0.0
	var i = 0.0

	if real(sample) > 0 {
		r = 1
	}

	if imag(sample) > 0 {
		i = 1
	}

	return complex64(complex(r, i))
}

func (ccr *ComplexClockRecovery) internalWork(input []complex64) []complex64 {
	var output = make([]complex64, 0)
	var inputIndex = 0
	var nInput = len(input) - ccr.interp.GetNTaps() - fudgeFactor
	var mmVal float32
	var u, x, y complex64

	for inputIndex < nInput {
		ccr.p2T = ccr.p1T
		ccr.p1T = ccr.p0T
		ccr.p0T = ccr.interp.Interpolate(input[inputIndex:], ccr.mu)

		ccr.c2T = ccr.c1T
		ccr.c1T = ccr.c0T
		ccr.c0T = slicer(ccr.p0T)

		x = (ccr.c0T - ccr.c2T) * dsp.Conj(ccr.p1T)
		y = (ccr.p0T - ccr.p2T) * dsp.Conj(ccr.c1T)
		u = y - x

		mmVal = real(u)
		output = append(output, ccr.p0T)

		mmVal = dsp.Clip(mmVal, 1.0)

		ccr.omega = ccr.omega + ccr.gainOmega*mmVal
		ccr.omega = ccr.omegaMidValue + dsp.Clip(ccr.omega-ccr.omegaMidValue, ccr.omegaLimit)

		ccr.mu = ccr.mu + ccr.omega + ccr.gainMu*mmVal
		inputIndex += int(dsp.Floor(ccr.mu))
		ccr.mu -= dsp.Floor(ccr.mu)

		if inputIndex < 0 {
			inputIndex = 0
		}
	}

	ccr.consumed = inputIndex

	if ccr.consumed > len(input) {
		panic("Consumed more samples than input!")
	}

	return output
}

func (ccr *ComplexClockRecovery) Work(input []complex64) []complex64 {
	var buff = make([]complex64, len(input)+ccr.sampleHistoryCount)
	copy(buff, ccr.sampleHistory[:ccr.sampleHistoryCount])
	copy(buff[ccr.sampleHistoryCount:], input)

	var symbols = ccr.internalWork(buff)

	ccr.sampleHistoryCount = len(buff) - ccr.consumed

	if ccr.sampleHistoryCount < ccHistoryLength {
		ccr.sampleHistoryCount = ccHistoryLength
	}

	if len(ccr.sampleHistory) != ccr.sampleHistoryCount {
		ccr.sampleHistory = make([]complex64, ccr.sampleHistoryCount)
	}

	copy(ccr.sampleHistory, buff[len(buff)-ccr.sampleHistoryCount:])

	return symbols
}
