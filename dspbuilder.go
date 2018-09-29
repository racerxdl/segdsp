package main

import (
	"fmt"
	"github.com/racerxdl/segdsp/demodcore"
)

func buildFM(sampleRate uint32) *demodcore.FMDemod {
	return demodcore.MakeCustomFMDemodulator(sampleRate, float64(filterBandwidth), uint32(outputRate), fmTau, squelch, squelchAlpha, float32(fmDeviation))
}
func buildAM(sampleRate uint32) *demodcore.AMDemod {
	return demodcore.MakeCustomAMDemodulator(sampleRate, float64(filterBandwidth), uint32(outputRate), amAudioCut, squelch, squelchAlpha)
}

func buildDSP(sampleRate uint32) demodcore.DemodCore {
	switch demodulatorMode {
	case modeFM:
		return buildFM(sampleRate)
	case modeAM:
		return buildAM(sampleRate)
	}

	panic(fmt.Sprintf("Unsupported Mode: %s", demodulatorMode))
}
