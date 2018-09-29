package main

type presetStruct struct {
	name            string
	demodMode       string
	outputRate      uint
	filterBandwidth float64
	demodOptions    map[string]interface{}
}

var presets = map[string]presetStruct{
	"am": {
		name:            "AM",
		demodMode:       modeAM,
		outputRate:      48000,
		filterBandwidth: 10e3,
		demodOptions: map[string]interface{}{
			"audioCut": 5e3,
		},
	},
	"nbfm": {
		name:            "Narrow Band FM",
		demodMode:       modeFM,
		outputRate:      48000,
		filterBandwidth: 10e3,
		demodOptions: map[string]interface{}{
			"deviation": 5e3,
			"tau":       75e-6,
		},
	},
	"wbfm": {
		name:            "Wide Band FM",
		demodMode:       modeFM,
		outputRate:      48000,
		filterBandwidth: 120e3,
		demodOptions: map[string]interface{}{
			"deviation": 75e3,
			"tau":       75e-6,
		},
	},
}
