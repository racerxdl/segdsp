package main


import (
	"github.com/racerxdl/segdsp/demodcore"
	"fmt"
)

func BuildFM(sampleRate uint32) *demodcore.FMDemod {
	return demodcore.MakeCustomFMDemodulator(sampleRate, float64(filterBandwidth), uint32(outputRate), fmTau, fmSquelch, fmSquelchAlpha, float32(fmDeviation))
}

func BuildDSP(sampleRate uint32) demodcore.DemodCore {
	switch demodulatorMode {
	case modeFM:
		return BuildFM(sampleRate)
	}

	panic(fmt.Sprintf("Unsupported Mode: %s", demodulatorMode))
}


