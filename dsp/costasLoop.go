package dsp

import (
	"github.com/racerxdl/segdsp/tools"
	"math"
)

type CostasLoop interface {
	ComplexWorker
	GetError() float32
	GetFrequency() float32
	GetFrequencyHz() float32
}

// region Base Costas Loop
type baseCostasLoop struct {
	error float32
}

func (cl *baseCostasLoop) GetError() float32 {
	return cl.error
}

func (cl *baseCostasLoop) PredictOutputSize(input int) int {
	return input
}

// endregion
// region Second Order Costas Loop
type CostasLoop2 struct {
	controlLoop
	baseCostasLoop
}

func MakeCostasLoop2WithFrequencyRange(loopBandwidth, minRelativeFrequency, maxRelativeFrequency float32) CostasLoop {
	cl := makeControlLoop(loopBandwidth, minRelativeFrequency, maxRelativeFrequency)
	cl2 := &CostasLoop2{controlLoop: *cl, baseCostasLoop: baseCostasLoop{error: 0}}

	return cl2
}

func MakeCostasLoop2(loopBandwidth float32) CostasLoop {
	return MakeCostasLoop2WithFrequencyRange(loopBandwidth, -1, 1)
}

func (cl *CostasLoop2) Work(input []complex64) []complex64 {
	output := make([]complex64, cl.PredictOutputSize(len(input)))
	cl.WorkBuffer(input, output)
	return output
}

func (cl *CostasLoop2) WorkBuffer(input, output []complex64) int {
	for i := 0; i < len(input); i++ {
		nr := tools.Cos(-cl.phase)
		ni := tools.Sin(-cl.phase)

		n := complex(nr, ni)
		output[i] = input[i] * n

		cl.error = real(output[i]) * imag(output[i])
		cl.error = tools.Clip(cl.error, 1)
		cl.AdvanceLoop(cl.error)
		cl.phaseWrap()
		cl.frequencyLimit()
	}

	return len(input)
}

// endregion
// region 4th Order Costas Loop
type CostasLoop4 struct {
	controlLoop
	baseCostasLoop
}

func MakeCostasLoop4WithFrequencyRange(loopBandwidth, minRelativeFrequency, maxRelativeFrequency float32) CostasLoop {
	cl := makeControlLoop(loopBandwidth, minRelativeFrequency, maxRelativeFrequency)
	cl2 := &CostasLoop4{controlLoop: *cl, baseCostasLoop: baseCostasLoop{error: 0}}

	return cl2
}

func MakeCostasLoop4(loopBandwidth float32) CostasLoop {
	return MakeCostasLoop4WithFrequencyRange(loopBandwidth, -1, 1)
}

func (cl *CostasLoop4) Work(input []complex64) []complex64 {
	output := make([]complex64, cl.PredictOutputSize(len(input)))
	cl.WorkBuffer(input, output)
	return output
}

func (cl *CostasLoop4) WorkBuffer(input, output []complex64) int {
	for i := 0; i < len(input); i++ {
		nr := tools.Cos(-cl.phase)
		ni := tools.Sin(-cl.phase)

		n := complex(nr, ni)
		output[i] = input[i] * n

		vr := float32(1)
		vi := float32(1)

		if real(output[i]) <= 0 {
			vr = -1
		}

		if imag(output[i]) <= 0 {
			vi = -1
		}

		cl.error = imag(output[i])*vr - real(output[i])*vi
		cl.error = tools.Clip(cl.error, 1)
		cl.AdvanceLoop(cl.error)
		cl.phaseWrap()
		cl.frequencyLimit()
	}

	return len(input)
}

// endregion
// region 8th Order Costas Loop
type CostasLoop8 struct {
	controlLoop
	baseCostasLoop
}

func MakeCostasLoop8WithFrequencyRange(loopBandwidth, minRelativeFrequency, maxRelativeFrequency float32) CostasLoop {
	cl := makeControlLoop(loopBandwidth, minRelativeFrequency, maxRelativeFrequency)
	cl2 := &CostasLoop8{controlLoop: *cl, baseCostasLoop: baseCostasLoop{error: 0}}

	return cl2
}

func MakeCostasLoop8(loopBandwidth float32) CostasLoop {
	return MakeCostasLoop2WithFrequencyRange(loopBandwidth, -1, 1)
}

func (cl *CostasLoop8) GetError() float32 {
	return cl.error
}

func (cl *CostasLoop8) Work(input []complex64) []complex64 {
	output := make([]complex64, cl.PredictOutputSize(len(input)))
	cl.WorkBuffer(input, output)
	return output
}

func (cl *CostasLoop8) WorkBuffer(input, output []complex64) int {
	K := float32(math.Sqrt(2) - 1)
	for i := 0; i < len(input); i++ {
		nr := tools.Cos(-cl.phase)
		ni := tools.Sin(-cl.phase)

		n := complex(nr, ni)
		output[i] = input[i] * n

		vr := float32(1)
		vi := float32(1)

		if real(output[i]) < 0 {
			vr = -1
		}

		if imag(output[i]) < 0 {
			vi = -1
		}

		if tools.Abs(real(output[i])) > tools.Abs(imag(output[i])) {
			cl.error = imag(output[i])*vr - real(output[i])*vi*K
		} else {
			cl.error = imag(output[i])*vr*K - real(output[i])*vi
		}

		cl.error = imag(output[i])*vr - real(output[i])*vi
		cl.error = tools.Clip(cl.error, 1)
		cl.AdvanceLoop(cl.error)
		cl.phaseWrap()
		cl.frequencyLimit()
	}

	return len(input)
}

// endregion
