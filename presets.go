package main

type presetStruct struct {
	name            string
	demodMode       string
	outputRate      uint
	filterBandwidth float64
	demodOptions    demodOptions
}

type demodOptions struct {
	fmDeviation float64
	fmTau       float64
	amAudioCut  float64
}

var presets = map[string]presetStruct{
	"am": {
		name:            "AM",
		demodMode:       modeAM,
		outputRate:      48000,
		filterBandwidth: 10e3,
		demodOptions: demodOptions{
			amAudioCut: 5e3,
		},
	},
	"nbfm": {
		name:            "Narrow Band FM",
		demodMode:       modeFM,
		outputRate:      48000,
		filterBandwidth: 10e3,
		demodOptions: demodOptions{
			fmDeviation: 5e3,
			fmTau:       75e-6,
		},
	},
	"wbfm": {
		name:            "Wide Band FM",
		demodMode:       modeFM,
		outputRate:      48000,
		filterBandwidth: 120e3,
		demodOptions: demodOptions{
			fmDeviation: 75e3,
			fmTau:       75e-6,
		},
	},
}
