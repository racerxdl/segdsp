package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// region Modes

const modeFM = "FM"

var modes = []string {modeFM}

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
const envFSBW = "FS_BANDWIDTH"
const envStationName = "STATION_NAME"
const envWebCanControl = "WEB_CAN_CONTROL"
const envTCPCanControl = "TCP_CAN_CONTROL"

// region FM Demodulator Options
const envFMDeviation = "FM_DEVIATION"
const envFMTau = "FM_TAU"
const envFMSquelch = "FM_SQUELCH"
const envFMSquelchAlpha = "FM_SQUELCH_ALPHA"
// endregion

// endregion
// region Arguments

var addrFlag = flag.String("httpAddr", "localhost:8080", "http service address")
var spyserverhostFlag = flag.String("spyserver", "localhost:5555", "spyserver address")
var displayPixelsFlag = flag.Uint("displayPixels", 512, "Width in pixels of the FFT")

var channelFrequencyFlag = flag.Uint("channelFrequency", 106.3e6, "Channel (IQ) Center Frequency")
var displayFrequencyFlag = flag.Uint("fftFrequency", 106.3e6, "FFT Center Frequency")

var channelDecimationStageFlag = flag.Uint("decimationStage", 2, "Channel (IQ) Decimation Stage (The actual decimation will be 2^d)")
var displayDecimationStageFlag = flag.Uint("fftDecimationStage", 1, "FFT Decimation Stage (The actual decimation will be 2^d)")

var demodulatorModeFlag = flag.String("demodMode", modeFM, fmt.Sprintf("Demodulator Mode: %s", modes))
var outputRateFlag = flag.Uint("outputRate", 48000, "Output Rate in Hertz")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var stationNameFlag = flag.String("stationName", "SegDSP", "Your station name or callsign (it identifies this instance)")
var webCanControlFlag = flag.Bool("webCanControl", false, "If Web UI Clients can control this server")
var tcpCanControlFlag = flag.Bool("tcpCanControl", false, "If TCP Clients can control this server")

// region FM Demodulator Flags
var filterBandwidthFlag = flag.Uint("filterBandwidth", 120e3, "First Stage Filter Bandwidth in Hertz")
var fmDeviationFlag = flag.Uint("fmDeviation", 75e3, "FM Demodulator Max Deviation in Hertz")
var fmTauFlag = flag.Float64("fmTau", 75e-6, "FM Demodulator Tau in seconds (0 to disable)")
var fmSquelchFlag = flag.Float64("fmSquelch", -65, "FM Demodulator Squelch in dB")
var fmSquelchAlphaFlag = flag.Float64("fmSquelchAlpha", 0.001, "FM Demodulator Squelch Filter Alpha")
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
var filterBandwidth uint

var fmDeviation uint
var fmTau float32
var fmSquelch float32
var fmSquelchAlpha float32

var stationName string
var webCanControl bool
var tcpCanControl bool
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

	if os.Getenv(envFSBW) == "" {
		os.Setenv(envFSBW, strconv.FormatUint(uint64(*filterBandwidthFlag), 10))
	}

	if os.Getenv(envFMDeviation) == "" {
		os.Setenv(envFMDeviation, strconv.FormatUint(uint64(*fmDeviationFlag), 10))
	}

	if os.Getenv(envFMTau) == "" {
		os.Setenv(envFMTau, strconv.FormatFloat(*fmTauFlag, 'E', -1, 32))
	}

	if os.Getenv(envFMSquelch) == "" {
		os.Setenv(envFMSquelch, strconv.FormatFloat(*fmSquelchFlag, 'E', -1, 32))
	}

	if os.Getenv(envFMSquelchAlpha) == "" {
		os.Setenv(envFMSquelchAlpha, strconv.FormatFloat(*fmSquelchAlphaFlag, 'E', -1, 32))
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
	fmsquelch, err := strconv.ParseFloat(os.Getenv(envFMSquelch), 32)
	if err != nil {
		panic(err)
	}
	fmSquelch = float32(fmsquelch)
	fmsquelchalpha, err := strconv.ParseFloat(os.Getenv(envFMSquelchAlpha), 32)
	if err != nil {
		panic(err)
	}
	fmSquelchAlpha = float32(fmsquelchalpha)

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
	// endregion
}