package main

import (
	"flag"
	"fmt"
	"github.com/racerxdl/segdsp/recorders"
	"log"
	"os"
	"strconv"
)

const modeFM = "FM"
const modeAM = "AM"

var modes = []string{modeFM, modeAM}

var addrFlag = flag.String("httpAddr", "localhost:8080", "http service address")
var radioserverhostFlag = flag.String("radioserver", "localhost:4050", "radioserver address")
var displayPixelsFlag = flag.Uint("displayPixels", 512, "Width in pixels of the FFT")
var channelFrequencyFlag = flag.Uint("channelFrequency", 106.3e6, "Channel (IQ) Center Frequency")
var displayFrequencyFlag = flag.Uint("fftFrequency", 106.3e6, "FFT Center Frequency")
var channelDecimationStageFlag = flag.Uint("decimationStage", 3, "Channel (IQ) Decimation Stage")
var displayDecimationStageFlag = flag.Uint("fftDecimationStage", 1, "FFT Decimation Stage")
var demodulatorModeFlag = flag.String("demodMode", modeFM, fmt.Sprintf("Demodulator Mode: %s", modes))
var outputRateFlag = flag.Uint("outputRate", 48000, "Output Rate in Hertz")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var stationNameFlag = flag.String("stationName", "SegDSP", "Station name or callsign")
var webCanControlFlag = flag.Bool("webCanControl", false, "Web UI can control server")
var tcpCanControlFlag = flag.Bool("tcpCanControl", true, "TCP clients can control server")
var recordFlag = flag.Bool("record", false, "Record when not squelched")
var recordMethodFlag = flag.String("recordMethod", recorders.RecFile, "Recording method")
var presetFlag = flag.String("preset", "none", "Preset for Demodulator Params")
var squelchFlag = flag.Float64("squelch", -150, "Demodulator Squelch in dB")
var squelchAlphaFlag = flag.Float64("squelchAlpha", 0.001, "Demodulator Squelch Filter Alpha")
var filterBandwidthFlag = flag.Uint("filterBandwidth", 120e3, "First Stage Filter Bandwidth in Hz")
var fmDeviationFlag = flag.Uint("fmDeviation", 75e3, "FM Max Deviation in Hz")
var fmTauFlag = flag.Float64("fmTau", 75e-6, "FM Tau in seconds (0 to disable)")
var amAudioCutFlag = flag.Float64("amAudioCut", 5000, "AM Low Pass Filter Cut")

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

func envString(envKey, fallback string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return fallback
}

func envUint(envKey string, fallback uint) uint {
	if v := os.Getenv(envKey); v != "" {
		n, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			panic(fmt.Sprintf("invalid value for %s: %v", envKey, err))
		}
		return uint(n)
	}
	return fallback
}

func envFloat32(envKey string, fallback float32) float32 {
	if v := os.Getenv(envKey); v != "" {
		n, err := strconv.ParseFloat(v, 32)
		if err != nil {
			panic(fmt.Sprintf("invalid value for %s: %v", envKey, err))
		}
		return float32(n)
	}
	return fallback
}

func envBool(envKey string, fallback bool) bool {
	if v := os.Getenv(envKey); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("invalid value for %s: %v", envKey, err))
		}
		return b
	}
	return fallback
}

func applyPreset(p presetStruct) {
	log.Printf("PRESET: %s — outputRate=%d, mode=%s, fsbw=%.0f\n", p.name, p.outputRate, p.demodMode, p.filterBandwidth)
	_ = os.Setenv("OUTPUT_RATE", strconv.FormatUint(uint64(p.outputRate), 10))
	_ = os.Setenv("DEMOD_MODE", p.demodMode)
	_ = os.Setenv("FS_BANDWIDTH", strconv.FormatFloat(p.filterBandwidth, 'E', -1, 32))

	switch p.demodMode {
	case modeFM:
		_ = os.Setenv("FM_TAU", strconv.FormatFloat(p.demodOptions["tau"].(float64), 'E', -1, 32))
		_ = os.Setenv("FM_DEVIATION", strconv.FormatFloat(p.demodOptions["deviation"].(float64), 'E', -1, 32))
	case modeAM:
		_ = os.Setenv("AM_AUDIO_CUT", strconv.FormatFloat(p.demodOptions["audioCut"].(float64), 'E', -1, 32))
	}
}

func setEnv() {
	flag.Parse()

	if val, ok := presets[*presetFlag]; ok {
		log.Printf("Selected %s preset.\n", val.name)
		applyPreset(val)
	}

	httpAddr = envString("HTTP_ADDRESS", *addrFlag)
	radioserverhost = envString("RADIOSERVER", *radioserverhostFlag)
	displayPixels = envUint("DISPLAY_PIXELS", *displayPixelsFlag)
	channelFrequency = envUint("CENTER_FREQUENCY", *channelFrequencyFlag)
	displayFrequency = envUint("FFT_FREQUENCY", *displayFrequencyFlag)
	channelDecimationStage = envUint("DECIMATION_STAGE", *channelDecimationStageFlag)
	displayDecimationStage = envUint("FFT_DECIMATION_STAGE", *displayDecimationStageFlag)
	demodulatorMode = envString("DEMOD_MODE", *demodulatorModeFlag)
	outputRate = envUint("OUTPUT_RATE", *outputRateFlag)
	filterBandwidth = envUint("FS_BANDWIDTH", *filterBandwidthFlag)
	fmDeviation = envUint("FM_DEVIATION", *fmDeviationFlag)
	fmTau = envFloat32("FM_TAU", float32(*fmTauFlag))
	squelch = envFloat32("SQUELCH", float32(*squelchFlag))
	squelchAlpha = envFloat32("SQUELCH_ALPHA", float32(*squelchAlphaFlag))
	amAudioCut = envFloat32("AM_AUDIO_CUT", float32(*amAudioCutFlag))
	stationName = envString("STATION_NAME", *stationNameFlag)
	webCanControl = envBool("WEB_CAN_CONTROL", *webCanControlFlag)
	tcpCanControl = envBool("TCP_CAN_CONTROL", *tcpCanControlFlag)
	record = envBool("RECORD", *recordFlag)
	recordMethod = envString("RECORD_METHOD", *recordMethodFlag)
}
