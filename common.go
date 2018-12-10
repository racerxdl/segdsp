package main

import (
	"flag"
	"fmt"
	"github.com/racerxdl/segdsp/recorders"
	"log"
	"os"
	"strconv"
)

// region Modes

const modeFM = "FM"
const modeAM = "AM"

var modes = []string{modeFM, modeAM}

// endregion

// region Environment Variables
const envRadioServerAddr = "RADIOSERVER"
const envCenterFrequency = "CENTER_FREQUENCY"
const envFFTFrequency = "FFT_FREQUENCY"
const envHTTPAddr = "HTTP_ADDRESS"
const envDisplayPixels = "DISPLAY_PIXELS"
const envDecimationStage = "DECIMATION_STAGE"
const envFFTDecimationStage = "FFT_DECIMATION_STAGE"
const envOutputRate = "OUTPUT_RATE"
const envMode = "DEMOD_MODE"
const envFSBW = "FS_BANDWIDTH"
const envStationName = "STATION_NAME"
const envWebCanControl = "WEB_CAN_CONTROL"
const envTCPCanControl = "TCP_CAN_CONTROL"
const envRecord = "RECORD"
const envRecordMethod = "RECORD_METHOD"
const envPreset = "PRESET"

const envSquelch = "SQUELCH"
const envSquelchAlpha = "SQUELCH_ALPHA"

// region FM Demodulator Options
const envFMDeviation = "FM_DEVIATION"
const envFMTau = "FM_TAU"

// endregion

// region AM Demodulator Options
const envAMAudioCut = "AM_AUDIO_CUT"

// endregion

// endregion
// region Arguments

var addrFlag = flag.String("httpAddr", "localhost:8080", "http service address")
var radioserverhostFlag = flag.String("radioserver", "localhost:4050", "radioserver address")
var displayPixelsFlag = flag.Uint("displayPixels", 512, "Width in pixels of the FFT")

var channelFrequencyFlag = flag.Uint("channelFrequency", 106.3e6, "Channel (IQ) Center Frequency")
var displayFrequencyFlag = flag.Uint("fftFrequency", 106e6, "FFT Center Frequency")

var channelDecimationStageFlag = flag.Uint("decimationStage", 3, "Channel (IQ) Decimation Stage (The actual decimation will be 2^d)")
var displayDecimationStageFlag = flag.Uint("fftDecimationStage", 1, "FFT Decimation Stage (The actual decimation will be 2^d)")

var demodulatorModeFlag = flag.String("demodMode", modeFM, fmt.Sprintf("Demodulator Mode: %s", modes))
var outputRateFlag = flag.Uint("outputRate", 48000, "Output Rate in Hertz")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var stationNameFlag = flag.String("stationName", "SegDSP", "Your station name or callsign (it identifies this instance)")
var webCanControlFlag = flag.Bool("webCanControl", false, "If Web UI Clients can control this server")
var tcpCanControlFlag = flag.Bool("tcpCanControl", true, "If TCP Clients can control this server")

var recordFlag = flag.Bool("record", false, "If it should record when not squelched")
var recordMethodFlag = flag.String("recordMethod", recorders.RecFile, "Method to use when recording")

var presetFlag = flag.String("preset", "none", "presetStruct for Demodulator Params")

var squelchFlag = flag.Float64("squelch", -72, "Demodulator Squelch in dB")
var squelchAlphaFlag = flag.Float64("squelchAlpha", 0.001, "Demodulator Squelch Filter Alpha")

// region FM Demodulator Flags
var filterBandwidthFlag = flag.Uint("filterBandwidth", 120e3, "First Stage Filter Bandwidth in Hertz")
var fmDeviationFlag = flag.Uint("fmDeviation", 75e3, "FM Demodulator Max Deviation in Hertz")
var fmTauFlag = flag.Float64("fmTau", 75e-6, "FM Demodulator Tau in seconds (0 to disable)")

// endregion

// region AM Demodulator Flags
var amAudioCutFlag = flag.Float64("amAudioCut", 5000, "AM Low Pass Filter Cut")

// endregion

// endregion
// region Variables
var httpAddr string
var radioserverhost string
var displayPixels uint

var channelFrequency uint
var displayFrequency uint

var channelDecimationStage uint
var displayDecimationStage uint

var demodulatorMode string
var outputRate uint
var filterBandwidth uint
var squelch float32
var squelchAlpha float32

var fmDeviation uint
var fmTau float32

var amAudioCut float32

var stationName string
var webCanControl bool
var tcpCanControl bool
var record bool
var recordMethod string
var preset string

// endregion

func applyPreset(preset presetStruct) {
	log.Printf("PRESET: Setting Output Rate to %d Hz\n", preset.outputRate)
	log.Printf("PRESET: Setting Demod Mode to %s\n", preset.demodMode)
	log.Printf("PRESET: Setting First Stage Filter to %f Hz\n", preset.filterBandwidth)

	os.Setenv(envOutputRate, strconv.FormatUint(uint64(preset.outputRate), 10))
	os.Setenv(envMode, preset.demodMode)
	os.Setenv(envFSBW, strconv.FormatFloat(preset.filterBandwidth, 'E', -1, 32))

	switch preset.demodMode {
	case modeFM:
		applyFMPreset(preset)
	case modeAM:
		applyAMPreset(preset)
	}
}

func applyFMPreset(preset presetStruct) {
	log.Printf("PRESET: Setting FM Tau to %f\n", preset.demodOptions["tau"].(float64))
	log.Printf("PRESET: Setting FM Devation to %f Hz\n", preset.demodOptions["devation"].(float64))
	os.Setenv(envFMTau, strconv.FormatFloat(preset.demodOptions["tau"].(float64), 'E', -1, 32))
	os.Setenv(envFMDeviation, strconv.FormatFloat(preset.demodOptions["deviation"].(float64), 'E', -1, 32))
}

func applyAMPreset(preset presetStruct) {
	log.Printf("PRESET: Setting AM Audio Cut to %f\n", preset.demodOptions["audioCut"].(float64))
	os.Setenv(envAMAudioCut, strconv.FormatFloat(preset.demodOptions["audioCut"].(float64), 'E', -1, 32))
}

func setEnv() {
	flag.Parse()
	// region Parse presetStruct
	if val, ok := presets[preset]; ok {
		log.Printf("Selected %s preset.\n", val.name)
		applyPreset(val)
	}
	// endregion
	// region Fill Environment
	if os.Getenv(envRadioServerAddr) == "" {
		os.Setenv(envRadioServerAddr, *radioserverhostFlag)
	}

	if os.Getenv(envCenterFrequency) == "" {
		os.Setenv(envCenterFrequency, strconv.FormatUint(uint64(*channelFrequencyFlag), 10))
	}

	if os.Getenv(envFFTFrequency) == "" {
		os.Setenv(envFFTFrequency, strconv.FormatUint(uint64(*displayFrequencyFlag), 10))
	}

	if os.Getenv(envHTTPAddr) == "" {
		os.Setenv(envHTTPAddr, *addrFlag)
	}

	if os.Getenv(envFFTFrequency) == "" {
		os.Setenv(envFFTFrequency, strconv.FormatUint(uint64(*displayFrequencyFlag), 10))
	}

	if os.Getenv(envDisplayPixels) == "" {
		os.Setenv(envDisplayPixels, strconv.FormatUint(uint64(*displayPixelsFlag), 10))
	}

	if os.Getenv(envDecimationStage) == "" {
		os.Setenv(envDecimationStage, strconv.FormatUint(uint64(*channelDecimationStageFlag), 10))
	}

	if os.Getenv(envFFTDecimationStage) == "" {
		os.Setenv(envFFTDecimationStage, strconv.FormatUint(uint64(*displayDecimationStageFlag), 10))
	}

	if os.Getenv(envMode) == "" {
		os.Setenv(envMode, *demodulatorModeFlag)
	}

	if os.Getenv(envOutputRate) == "" {
		os.Setenv(envOutputRate, strconv.FormatUint(uint64(*outputRateFlag), 10))
	}

	if os.Getenv(envFSBW) == "" {
		os.Setenv(envFSBW, strconv.FormatUint(uint64(*filterBandwidthFlag), 10))
	}

	if os.Getenv(envFMDeviation) == "" {
		os.Setenv(envFMDeviation, strconv.FormatUint(uint64(*fmDeviationFlag), 10))
	}

	if os.Getenv(envFMTau) == "" {
		os.Setenv(envFMTau, strconv.FormatFloat(*fmTauFlag, 'E', -1, 32))
	}

	if os.Getenv(envSquelch) == "" {
		os.Setenv(envSquelch, strconv.FormatFloat(*squelchFlag, 'E', -1, 32))
	}

	if os.Getenv(envSquelchAlpha) == "" {
		os.Setenv(envSquelchAlpha, strconv.FormatFloat(*squelchAlphaFlag, 'E', -1, 32))
	}

	if os.Getenv(envAMAudioCut) == "" {
		os.Setenv(envAMAudioCut, strconv.FormatFloat(*amAudioCutFlag, 'E', -1, 32))
	}

	if os.Getenv(envStationName) == "" {
		os.Setenv(envStationName, *stationNameFlag)
	}

	if os.Getenv(envWebCanControl) == "" {
		os.Setenv(envWebCanControl, strconv.FormatBool(*webCanControlFlag))
	}

	if os.Getenv(envTCPCanControl) == "" {
		os.Setenv(envTCPCanControl, strconv.FormatBool(*tcpCanControlFlag))
	}

	if os.Getenv(envRecord) == "" {
		os.Setenv(envRecord, strconv.FormatBool(*recordFlag))
	}

	if os.Getenv(envRecordMethod) == "" {
		os.Setenv(envRecordMethod, *recordMethodFlag)
	}

	if os.Getenv(envPreset) == "" {
		os.Setenv(envPreset, *presetFlag)
	}

	// endregion
	// region Fill Variables
	httpAddr = os.Getenv(envHTTPAddr)
	radioserverhost = os.Getenv(envRadioServerAddr)
	dp, err := strconv.ParseUint(os.Getenv(envDisplayPixels), 10, 16)
	if err != nil {
		panic(err)
	}
	displayPixels = uint(dp)
	cf, err := strconv.ParseUint(os.Getenv(envCenterFrequency), 10, 32)
	if err != nil {
		panic(err)
	}
	channelFrequency = uint(cf)
	df, err := strconv.ParseUint(os.Getenv(envFFTFrequency), 10, 32)
	if err != nil {
		panic(err)
	}
	displayFrequency = uint(df)
	cds, err := strconv.ParseUint(os.Getenv(envDecimationStage), 10, 8)
	if err != nil {
		panic(err)
	}
	channelDecimationStage = uint(cds)
	dds, err := strconv.ParseUint(os.Getenv(envFFTDecimationStage), 10, 8)
	if err != nil {
		panic(err)
	}
	displayDecimationStage = uint(dds)
	demodulatorMode = os.Getenv(envMode)
	or, err := strconv.ParseUint(os.Getenv(envOutputRate), 10, 32)
	if err != nil {
		panic(err)
	}
	outputRate = uint(or)
	fsbw, err := strconv.ParseUint(os.Getenv(envFSBW), 10, 32)
	if err != nil {
		panic(err)
	}
	filterBandwidth = uint(fsbw)
	fmdev, err := strconv.ParseUint(os.Getenv(envFMDeviation), 10, 32)
	if err != nil {
		panic(err)
	}
	fmDeviation = uint(fmdev)
	fmtau, err := strconv.ParseFloat(os.Getenv(envFMTau), 32)
	if err != nil {
		panic(err)
	}
	fmTau = float32(fmtau)
	squelchx, err := strconv.ParseFloat(os.Getenv(envSquelch), 32)
	if err != nil {
		panic(err)
	}
	squelch = float32(squelchx)
	squelchalpha, err := strconv.ParseFloat(os.Getenv(envSquelchAlpha), 32)
	if err != nil {
		panic(err)
	}
	squelchAlpha = float32(squelchalpha)

	amaudiocut, err := strconv.ParseFloat(os.Getenv(envAMAudioCut), 32)
	if err != nil {
		panic(err)
	}
	amAudioCut = float32(amaudiocut)

	stationName = os.Getenv(envStationName)

	webcancontrol, err := strconv.ParseBool(os.Getenv(envWebCanControl))
	if err != nil {
		panic(err)
	}
	webCanControl = webcancontrol

	tcpcancontrol, err := strconv.ParseBool(os.Getenv(envTCPCanControl))
	if err != nil {
		panic(err)
	}
	tcpCanControl = tcpcancontrol

	recordf, err := strconv.ParseBool(os.Getenv(envRecord))
	if err != nil {
		panic(err)
	}
	record = recordf

	recordMethod = os.Getenv(envRecordMethod)

	preset = os.Getenv(envPreset)
	// endregion
}
