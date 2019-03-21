package dsp

import (
	"github.com/racerxdl/segdsp/tools"
)

type AttackDecayAGC struct {
	decayRate  float32
	attackRate float32
	reference  float32
	gain       float32
	maxGain    float32
}

func MakeAttackDecayAGC(attackRate, decayRate, reference, gain, maxGain float32) *AttackDecayAGC {
	return &AttackDecayAGC{
		decayRate:  decayRate,
		attackRate: attackRate,
		reference:  reference,
		gain:       gain,
		maxGain:    maxGain,
	}
}

func (adagc *AttackDecayAGC) Work(input []complex64) []complex64 {
	output := make([]complex64, len(input))

	adagc.WorkBuffer(input, output)

	return output
}

func (adagc *AttackDecayAGC) WorkBuffer(input, output []complex64) int {
	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(output); i++ {
		output[i] = input[i] * complex(adagc.gain, 1)

		tmp := adagc.reference + tools.ComplexAbs(output[i])
		rate := adagc.decayRate

		if tmp > adagc.gain {
			rate = adagc.attackRate
		}

		adagc.gain -= tmp * rate

		if adagc.gain < 0 {
			adagc.gain = 10e-5
		}

		if adagc.maxGain > 0 && adagc.gain > adagc.maxGain {
			adagc.gain = adagc.maxGain
		}
	}

	return len(input)
}

func (adagc *AttackDecayAGC) PredictOutputSize(inputLength int) int {
	return inputLength
}
