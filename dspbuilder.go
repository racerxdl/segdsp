package main

import (
	"fmt"
	"github.com/racerxdl/segdsp/demodcore"
)

func BuildFM(sampleRate uint32) *demodcore.FMDemod {
	return demodcore.MakeCustomFMDemodulator(sampleRate, float64(filterBandwidth), uint32(outputRate), fmTau, squelch, squelchAlpha, float32(fmDeviation))
}
func BuildAM(sampleRate uint32) *demodcore.AMDemod {
	return demodcore.MakeCustomAMDemodulator(sampleRate, float64(filterBandwidth), uint32(outputRate), amAudioCut, squelch, squelchAlpha)
}

func BuildDSP(sampleRate uint32) demodcore.DemodCore {
	switch demodulatorMode {
	case modeFM:
		return BuildFM(sampleRate)
	case modeAM:
		return BuildAM(sampleRate)
	}

	panic(fmt.Sprintf("Unsupported Mode: %s", demodulatorMode))
}


