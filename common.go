package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// region Modes

const modeWBFM = "WBFM"
const modeNFM = "NFM"

var modes = []string {modeWBFM, modeNFM}

// endregion

// region Environment Variables
const envSpyserverAddr = "SPYSERVER"
const envCenterFrequency = "CENTER_FREQUENCY"
const envFFTFrequency = "FFT_FREQUENCY"
const envHTTPAddr = "HTTP_ADDRESS"
const envDisplayPixels = "DISPLAY_PIXELS"
const envDecimationStage = "DECIMATION_STAGE"
const envFFTDecimationStage = "FFT_DECIMATION_STAGE"
const envOutputRate = "OUTPUT_RATE"
const envMode = "DEMOD_MODE"

// region FM Demodulator Options
const envFMBW = "FM_BANDWIDTH"
const envFMDeviation = "FM_DEVIATION"
const envFMTau = "FM_TAU"
// endregion

// endregion
// region Arguments

var addrFlag = flag.String("httpAddr", "localhost:8080", "http service address")
var spyserverhostFlag = flag.String("spyserver", "localhost:5555", "spyserver address")
var displayPixelsFlag = flag.Uint("displayPixels", 512, "Width in pixels of the FFT")

var channelFrequencyFlag = flag.Uint("channelFrequency", 106300000, "Channel (IQ) Center Frequency")
var displayFrequencyFlag = flag.Uint("fftFrequency", 106300000, "FFT Center Frequency")

var channelDecimationStageFlag = flag.Uint("decimationStage", 4, "Channel (IQ) Decimation Stage (The actual decimation will be 2^d)")
var displayDecimationStageFlag = flag.Uint("fftDecimationStage", 0, "FFT Decimation Stage (The actual decimation will be 2^d)")

var demodulatorModeFlag = flag.String("demodMode", modeWBFM, fmt.Sprintf("Demodulator Mode: %s", modes))
var outputRateFlag = flag.Uint("outputRate", 48000, "Output Rate in Hertz")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// region FM Demodulator Flags
var fmBandwidthFlag = flag.Uint("fmBandwidth", 120e3, "FM Demodulator Filter Bandwidth in Hertz")
var fmDeviationFlag = flag.Uint("fmDeviation", 75e3, "FM Demodulator Max Deviation in Hertz")
var fmTauFlag = flag.Float64("fmTau", 75e-6, "FM Demodulator Tau in seconds (0 to disable)")
// endregion

// endregion
// region Variables
var httpAddr string
var spyserverhost string
var displayPixels uint

var channelFrequency uint
var displayFrequency uint

var channelDecimationStage uint
var displayDecimationStage uint

var demodulatorMode string
var outputRate uint

var fmBandwidth uint
var fmDeviation uint
var fmTau float32
// endregion

func SetEnv() {
	flag.Parse()
	// region Fill Environment
	if os.Getenv(envSpyserverAddr) == "" {
		os.Setenv(envSpyserverAddr, *spyserverhostFlag)
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

	if os.Getenv(envFMBW) == "" {
		os.Setenv(envFMBW, strconv.FormatUint(uint64(*fmBandwidthFlag), 10))
	}

	if os.Getenv(envFMDeviation) == "" {
		os.Setenv(envFMDeviation, strconv.FormatUint(uint64(*fmDeviationFlag), 10))
	}

	if os.Getenv(envFMTau) == "" {
		os.Setenv(envFMTau, strconv.FormatFloat(*fmTauFlag, 'E', -1, 32))
	}

	if os.Getenv(envFMBW) == "" {
		os.Setenv(envFMBW, strconv.FormatUint(uint64(*fmBandwidthFlag), 10))
	}
	// endregion
	// region Fill Variables
	httpAddr = os.Getenv(envHTTPAddr)
	spyserverhost = os.Getenv(envSpyserverAddr)
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
	fmbw, err := strconv.ParseUint(os.Getenv(envFMBW), 10, 32)
	if err != nil {
		panic(err)
	}
	fmBandwidth = uint(fmbw)
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
	// endregion
}