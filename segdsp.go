package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/racerxdl/segdsp/demodcore"
	"github.com/racerxdl/segdsp/eventmanager"
	"github.com/racerxdl/segdsp/recorders"
	"github.com/racerxdl/spy2go/spyserver"
	"github.com/racerxdl/spy2go/spytypes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

type segdspCallback struct {
	ss *spyserver.Spyserver
}

func (cb *segdspCallback) OnData(dType int, data interface{}) {
	if dType == spytypes.SamplesComplex32 {
		go addS16Fifo(data.([]spytypes.ComplexInt16))
	} else if dType == spytypes.SamplesComplexUInt8 {
		go addU8Fifo(data.([]spytypes.ComplexUInt8))
	} else if dType == spytypes.DeviceSync {
		onDeviceSync(cb.ss)
	} else if dType == spytypes.FFTUInt8 {
		onFFT(data.([]byte))
	}
}

func onDeviceSync(ss *spyserver.Spyserver) {
	var d = deviceMessage{
		DeviceName: ss.GetName(),

		DisplayCenterFrequency: ss.GetDisplayCenterFrequency(),
		DisplayBandwidth:       ss.GetDisplayBandwidth(),
		DisplayOffset:          ss.GetDisplayOffset(),
		DisplayRange:           ss.GetDisplayRange(),
		DisplayPixels:          ss.GetDisplayPixels(),

		CurrentSampleRate:      ss.GetSampleRate(),
		ChannelCenterFrequency: ss.GetCenterFrequency(),
		Gain:              ss.GetGain(),
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
	sendPacket := currDevice.Gain != spyserver.InvalidValue
	if sendPacket {
		m, err := json.Marshal(currDevice)
		if err != nil {
			log.Println("Error serializing JSON: ", err)
		}
		go broadcastMessage(string(m))
	}
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
	//log.Println("Sending buffer")
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

	var ss = spyserver.MakeSpyserverByFullHS(spyserverhost)
	var cb = segdspCallback{
		ss: ss,
	}

	ss.SetCallback(&cb)

	ss.Connect()
	defer ss.Disconnect()

	log.Println(fmt.Sprintf("Device: %s", ss.GetName()))
	var srs = ss.GetAvailableSampleRates()

	log.Println("Available SampleRates:")
	for i := 0; i < len(srs); i++ {
		log.Println(fmt.Sprintf("		%f msps (dec stage %d)", float32(srs[i])/1e6, i))
	}

	ss.SetStreamingMode(spyserver.StreamModeFFTIQ)
	//ss.SetStreamingMode(spy2go.StreamModeIQOnly)
	ss.SetDisplayPixels(uint32(displayPixels))
	ss.SetDisplayDecimationStage(uint32(displayDecimationStage))
	ss.SetDisplayCenterFrequency(uint32(displayFrequency))
	ss.SetDisplayRange(90)
	ss.SetDisplayOffset(0)

	if ss.SetDecimationStage(uint32(channelDecimationStage)) == spyserver.InvalidValue {
		log.Println("Error setting sample rate.")
	}
	if ss.SetCenterFrequency(uint32(channelFrequency)) == spyserver.InvalidValue {
		log.Println("Error setting center frequency.")
	}

	time.Sleep(10 * time.Millisecond)

	log.Println("IQ Sample Rate: ", ss.GetSampleRate())
	log.Println("IQ Center Frequency ", ss.GetCenterFrequency())
	log.Println("FFT Center Frequency ", ss.GetDisplayCenterFrequency())
	log.Println("FFT Sample Rate: ", ss.GetDisplaySampleRate())

	demodulator = buildDSP(ss.GetSampleRate())
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
	ss.Start()

	<-done

	err = srv.Shutdown(context.TODO())
	if err != nil {
		log.Println(err)
	}

	log.Print("Stopping")
	ss.Stop()
	stopDSP()

	fmt.Println("Work Done")
}
