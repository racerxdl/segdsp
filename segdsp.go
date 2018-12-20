package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/racerxdl/radioserver/client"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/segdsp/demodcore"
	"github.com/racerxdl/segdsp/dsp/fft"
	"github.com/racerxdl/segdsp/eventmanager"
	"github.com/racerxdl/segdsp/recorders"
	"github.com/racerxdl/segdsp/tools"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

var displayRange = 90
var displayOffset = -50
var lastFFT []float32
var fftSamples []uint8
var fftSamplesF []float32
var smartIqSamples []complex64

type segdspCallback struct {
	rs *client.RadioClient
}

func (cb *segdspCallback) OnData(dType int, data interface{}) {
	switch dType {
	case client.SamplesComplex32:
		go addS16Fifo(data.([]client.ComplexInt16))
	case client.SmartSamplesComplex32:
		go onSmartIQ(cb.rs, data.([]client.ComplexInt16))
	case client.DeviceSync:
		onDeviceSync(cb.rs)
	}
}

func onDeviceSync(rs *client.RadioClient) {
	var d = deviceMessage{
		DeviceName: rs.GetName(),

		DisplayCenterFrequency: rs.GetSmartCenterFrequency(),
		DisplayBandwidth:       rs.GetSmartSampleRate(),
		DisplayOffset:          int32(displayOffset),
		DisplayRange:           int32(displayRange),
		DisplayPixels:          uint32(displayPixels),

		CurrentSampleRate:      rs.GetSampleRate(),
		ChannelCenterFrequency: rs.GetCenterFrequency(),
		Gain:              rs.GetGain(),
		OutputRate:        uint32(outputRate),
		FilterBandwidth:   uint32(filterBandwidth),
		DemodulatorMode:   demodulatorMode,
		DemodulatorParams: nil,
		StationName:       stationName,
		WebCanControl:     webCanControl,
		TCPCanControl:     tcpCanControl,
		IsMuted:           false,
	}

	if demodulator != nil {
		d.DemodulatorParams = demodulator.GetDemodParams()
		d.IsMuted = demodulator.IsMuted()
	}

	currDevice = makeDeviceMessage(d)
	refreshDevice()
}

func refreshDevice() {
	sendPacket := currDevice.Gain != protocol.Invalid
	if sendPacket {
		m, err := json.Marshal(currDevice)
		if err != nil {
			log.Println("Error serializing JSON: ", err)
		}
		go broadcastMessage(string(m))
	}
}

func onSmartIQ(rs *client.RadioClient, data []client.ComplexInt16) {
	var scale = 256 / float32(displayRange)
	data = data[:displayPixels]
	if smartIqSamples == nil || len(smartIqSamples) != len(data) {
		smartIqSamples = make([]complex64, len(data))
	}

	for i, v := range data {
		var a = float32(v.Real)
		var b = float32(v.Imag)
		smartIqSamples[i] = complex(a/32768, b/32768)
	}

	fftCData := fft.FFT(smartIqSamples)

	var l = len(fftCData)
	var scaledV float32

	if fftSamples == nil || len(fftSamples) != len(fftCData) {
		fftSamples = make([]uint8, len(fftCData))
		fftSamplesF = make([]float32, len(fftCData))
	}

	if lastFFT == nil || len(lastFFT) != len(fftSamples) {
		lastFFT = make([]float32, len(fftCData))
	}

	for i, v := range fftCData {
		var m = float64(tools.ComplexAbsSquared(v) * (1.0 / float32(rs.GetSmartSampleRate())))
		var o = float32(10 * math.Log10(m))
		fftSamplesF[i] = o
	}

	for i := range fftSamplesF {
		fftSamplesF[i] = (lastFFT[i] + fftSamplesF[i]) / 2
	}

	copy(lastFFT, fftSamplesF)

	for i, v := range fftSamplesF {
		// FFT is symmetric
		var oI = (i + l/2) % l
		scaledV = 255 + ((v - float32(displayOffset)) * scale)
		if scaledV < 0 {
			scaledV = 0
		} else if scaledV > 255 {
			scaledV = 255
		}
		fftSamples[oI] = uint8(scaledV)
	}
	onFFT(fftSamples)
}

func onFFT(data []uint8) {
	//log.Println("Received FFT! ", len(data))
	var j = makeFFTMessage(data, demodulator.GetLevel())
	m, err := json.Marshal(j)
	if err != nil {
		log.Println("Error serializing JSON: ", err)
	}
	go broadcastMessage(string(m))
}

func sendData(data interface{}) {
	switch data.(type) {
	case demodcore.DemodData:
		var b = data.(demodcore.DemodData)
		go broadcastBMessage(b.Data.MarshalByteArray())
		go recordAudio(b.Data)
	default:
		var j = makeDataMessage(data)
		m, err := json.Marshal(j)
		if err != nil {
			log.Println("Error serializing JSON: ", err)
		}
		go broadcastMessage(string(m))
	}
}

func createServer() *http.Server {
	srv := &http.Server{Addr: httpAddr}

	fs := http.FileServer(http.Dir("./content/static"))

	http.HandleFunc("/ws", ws)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", content)

	go func() {
		log.Println(http.ListenAndServe(httpAddr, nil))
	}()

	return srv
}

var squelchOn chan interface{}
var squelchOff chan interface{}

func onSquelchOn(data eventmanager.SquelchEventData) {
	log.Println("Squelch ON", data.AvgValue, data.Threshold)
	currDevice.IsMuted = demodulator.IsMuted()
	stopRecording()
	refreshDevice()
}

func onSquelchOff(data eventmanager.SquelchEventData) {
	log.Println("Squelch OFF", data.AvgValue, data.Threshold)
	currDevice.IsMuted = demodulator.IsMuted()
	startRecording()
	refreshDevice()
}

func main() {
	var err error
	setEnv()
	log.SetFlags(0)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Starting CPU Profile")
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	squelchOn = make(chan interface{})
	squelchOff = make(chan interface{})

	var ev = eventmanager.EventManager{}

	ev.AddHandler(eventmanager.EvSquelchOn, squelchOn)
	ev.AddHandler(eventmanager.EvSquelchOff, squelchOff)

	go func() {
		log.Println("Starting Handler loop")
		for {
			select {
			case msg := <-squelchOn:
				onSquelchOn(msg.(eventmanager.SquelchEventData))
			case msg := <-squelchOff:
				onSquelchOff(msg.(eventmanager.SquelchEventData))
			}
		}
		//log.Println("Ending Handler loop")
	}()

	recorder = &recorders.FileRecorder{}
	recordingParams.recorderEnable = record

	if recordMethod != "file" {
		panic("Only\"file\" method is supported for recording.")
	}

	initDSP()
	var rs = client.MakeRadioClientByFullHS(radioserverhost)
	var cb = segdspCallback{
		rs: rs,
	}

	rs.SetCallback(&cb)

	rs.Connect()
	defer rs.Disconnect()

	log.Println(fmt.Sprintf("Device: %s", rs.GetName()))
	var srs = rs.GetAvailableSampleRates()

	log.Println("Available SampleRates:")
	for i := 0; i < len(srs); i++ {
		log.Println(fmt.Sprintf("		%f msps (dec stage %d)", float32(srs[i])/1e6, i))
	}

	rs.SetStreamingMode(protocol.TypeCombined)
	rs.SetSmartDecimation(uint32(displayDecimationStage))
	rs.SetSmartCenterFrequency(uint32(displayFrequency))

	if rs.SetDecimationStage(uint32(channelDecimationStage)) == protocol.Invalid {
		log.Println("Error setting sample rate.")
	}
	if rs.SetCenterFrequency(uint32(channelFrequency)) == protocol.Invalid {
		log.Println("Error setting center frequency.")
	}

	time.Sleep(10 * time.Millisecond)

	log.Println("IQ Sample Rate: ", rs.GetSampleRate())
	log.Println("IQ Center Frequency ", rs.GetCenterFrequency())
	log.Println("SmartIQ Center Frequency ", rs.GetSmartCenterFrequency())
	log.Println("SmartIQ Sample Rate: ", rs.GetSmartSampleRate())

	demodulator = buildDSP(rs.GetSampleRate())
	demodulator.SetEventManager(&ev)

	dspCb = sendData

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	var srv = createServer()

	startDSP()

	log.Println("Starting")
	rs.Start()

	<-done

	err = srv.Shutdown(context.TODO())
	if err != nil {
		log.Println(err)
	}

	log.Print("Stopping")
	rs.Stop()
	stopDSP()

	fmt.Println("Work Done")
}
