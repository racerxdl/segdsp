package digital

import (
	"github.com/racerxdl/segdsp/tools"
)

const ccHistoryLength = 3
const fudgeFactor = 16

func complexSlicer(sample complex64) complex64 {
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

func floatSlicer(sample float32) float32 {
	if sample > 0 {
		return 1
	} else {
		return -1
	}
}

// region Complex Clock Recovery
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

	interp *ComplexMMSEFirInterpolator
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
		interp:             MakeComplexMMSEFirInterpolator(),
		p0T:                0,
		p1T:                0,
		p2T:                0,
		c0T:                0,
		c1T:                0,
		c2T:                0,
		consumed:           0,
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

func (ccr *ComplexClockRecovery) internalWorkBuffer(input, output []complex64) int {
	var inputIndex = 0
	var nInput = len(input) - ccr.interp.GetNTaps() - fudgeFactor
	var nOutput = len(output)
	var mmVal float32
	var u, x, y complex64
	var outputIndex = 0

	for inputIndex < nInput && outputIndex < nOutput {
		ccr.p2T = ccr.p1T
		ccr.p1T = ccr.p0T
		ccr.p0T = ccr.interp.Interpolate(input[inputIndex:], ccr.mu)

		ccr.c2T = ccr.c1T
		ccr.c1T = ccr.c0T
		ccr.c0T = complexSlicer(ccr.p0T)

		x = (ccr.c0T - ccr.c2T) * tools.Conj(ccr.p1T)
		y = (ccr.p0T - ccr.p2T) * tools.Conj(ccr.c1T)
		u = y - x

		mmVal = real(u)
		output[outputIndex] = ccr.p0T
		outputIndex++

		mmVal = tools.Clip(mmVal, 1.0)

		ccr.omega = ccr.omega + ccr.gainOmega*mmVal
		ccr.omega = ccr.omegaMidValue + tools.Clip(ccr.omega-ccr.omegaMidValue, ccr.omegaLimit)

		ccr.mu = ccr.mu + ccr.omega + ccr.gainMu*mmVal
		inputIndex += int(tools.Floor(ccr.mu))
		ccr.mu -= tools.Floor(ccr.mu)

		if inputIndex < 0 {
			inputIndex = 0
		}
	}

	ccr.consumed = inputIndex

	if ccr.consumed > len(input) {
		panic("Consumed more samples than input!")
	}

	return outputIndex
}

func (ccr *ComplexClockRecovery) Work(input []complex64) []complex64 {
	var buff = make([]complex64, len(input)+ccr.sampleHistoryCount)
	l := ccr.WorkBuffer(input, buff)
	return buff[:l]
}

func (ccr *ComplexClockRecovery) WorkBuffer(input, output []complex64) int {
	var s = append(ccr.sampleHistory, input...)
	var symbols = ccr.internalWorkBuffer(s, output)

	ccr.sampleHistoryCount = len(s) - ccr.consumed

	if ccr.sampleHistoryCount < ccHistoryLength {
		ccr.sampleHistoryCount = ccHistoryLength
	}

	if len(ccr.sampleHistory) != ccr.sampleHistoryCount {
		ccr.sampleHistory = make([]complex64, ccr.sampleHistoryCount)
	}

	copy(ccr.sampleHistory, s[len(s)-ccr.sampleHistoryCount:])

	return symbols
}

func (ccr *ComplexClockRecovery) PredictOutputSize(inputLength int) int {
	return inputLength + ccr.sampleHistoryCount
}

// endregion
// region Float Clock Recovery

type FloatClockRecovery struct {
	lastSample         float32
	consumed           int
	sampleHistoryCount int
	sampleHistory      []float32

	mu                 float32
	omega              float32
	gainOmega          float32
	omegaRelativeLimit float32
	omegaMidValue      float32
	omegaLimit         float32
	gainMu             float32

	interp *FloatMMSEFirInterpolator
}

func (ccr *FloatClockRecovery) GetMu() float32 {
	return ccr.mu
}

func (ccr *FloatClockRecovery) GetOmega() float32 {
	return ccr.omega
}

func (ccr *FloatClockRecovery) GetGainMu() float32 {
	return ccr.gainMu
}

func (ccr *FloatClockRecovery) GetGainOmega() float32 {
	return ccr.gainOmega
}

func (ccr *FloatClockRecovery) SetGainMu(gain float32) {
	ccr.gainMu = gain
}

func (ccr *FloatClockRecovery) SetGainOmega(gain float32) {
	ccr.gainOmega = gain
}

func (ccr *FloatClockRecovery) SetMu(mu float32) {
	ccr.mu = mu
}

func (ccr *FloatClockRecovery) SetOmega(omega float32) {
	ccr.omega = omega
	ccr.omegaMidValue = omega
	ccr.omegaLimit = ccr.omegaRelativeLimit * omega
}

func NewFloatClockRecovery(omega, gainOmega, mu, gainMu, omegaRelativeLimit float32) *FloatClockRecovery {
	var ccr = &FloatClockRecovery{
		lastSample:         0,
		omega:              omega,
		gainOmega:          gainOmega,
		mu:                 mu,
		gainMu:             gainMu,
		omegaRelativeLimit: omegaRelativeLimit,
		interp:             MakeFloatMMSEFirInterpolator(),
	}

	if omega <= 0 {
		panic("Clock Recovery Rate (omega) must be higher than 0")
	}

	if gainMu < 0 || gainOmega < 0 {
		panic("Clock Recovery gains should be positive.")
	}

	ccr.SetOmega(omega)

	ccr.sampleHistoryCount = ccHistoryLength
	ccr.sampleHistory = make([]float32, ccHistoryLength)

	return ccr
}

func (ccr *FloatClockRecovery) internalWorkBuffer(input, output []float32) int {
	var inputIndex = 0
	var nInput = len(input) - ccr.interp.GetNTaps() - fudgeFactor
	var mmVal float32
	var nOutput = len(output)
	var outputIndex = 0

	for inputIndex < nInput && outputIndex < nOutput {
		var o = ccr.interp.Interpolate(input[inputIndex:], ccr.mu)
		output[outputIndex] = o
		outputIndex++

		mmVal = floatSlicer(ccr.lastSample)*o - floatSlicer(o)*ccr.lastSample
		ccr.lastSample = o

		ccr.omega = ccr.omega + ccr.gainOmega*mmVal
		ccr.omega = ccr.omegaMidValue + tools.Clip(ccr.omega-ccr.omegaMidValue, ccr.omegaLimit)

		ccr.mu = ccr.mu + ccr.omega + ccr.gainMu*mmVal
		inputIndex += int(tools.Floor(ccr.mu))
		ccr.mu -= tools.Floor(ccr.mu)

		if inputIndex < 0 {
			inputIndex = 0
		}
	}

	ccr.consumed = inputIndex

	if ccr.consumed > len(input) {
		panic("Consumed more samples than input!")
	}

	return outputIndex
}

func (ccr *FloatClockRecovery) Work(input []float32) []float32 {
	var buff = make([]float32, len(input)+ccr.sampleHistoryCount)
	l := ccr.WorkBuffer(input, buff)
	return buff[:l]
}

func (ccr *FloatClockRecovery) WorkBuffer(input, output []float32) int {
	var s = append(ccr.sampleHistory, input...)
	var symbols = ccr.internalWorkBuffer(s, output)

	ccr.sampleHistoryCount = len(s) - ccr.consumed

	if ccr.sampleHistoryCount < ccHistoryLength {
		ccr.sampleHistoryCount = ccHistoryLength
	}

	if len(ccr.sampleHistory) != ccr.sampleHistoryCount {
		ccr.sampleHistory = make([]float32, ccr.sampleHistoryCount)
	}

	copy(ccr.sampleHistory, s[len(s)-ccr.sampleHistoryCount:])

	return symbols
}

func (ccr *FloatClockRecovery) PredictOutputSize(inputLength int) int {
	return inputLength + ccr.sampleHistoryCount
}

// endregion
