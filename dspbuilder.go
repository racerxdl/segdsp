package main


import (
	"github.com/racerxdl/segdsp/demodcore"
	"fmt"
)

func BuildFM(sampleRate uint32) *demodcore.FMDemod {
	return demodcore.MakeCustomFMDemodulator(sampleRate, float64(fmBandwidth), uint32(outputRate), fmTau, float32(fmDeviation))
}

func BuildDSP(sampleRate uint32) demodcore.DemodCore {
	switch demodulatorMode {
	case modeNFM:
	case modeWBFM:
		return BuildFM(sampleRate)
	}

	panic(fmt.Sprintf("Unsupported Mode: %s", demodulatorMode))
}


