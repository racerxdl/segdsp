package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/racerxdl/segdsp/demodcore"
	"github.com/racerxdl/segdsp/eventmanager"
	"github.com/racerxdl/segdsp/recorders"
	"github.com/racerxdl/spy2go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

func onInt16IQ(data []spy2go.ComplexInt16) {
	go addS16Fifo(data)
}
func onUInt8IQ(data []spy2go.ComplexUInt8) {
	go addU8Fifo(data)
}

func onDeviceSync(spyserver *spy2go.Spyserver) {
	var d = deviceMessage{
		DeviceName: spyserver.GetName(),

		DisplayCenterFrequency: spyserver.GetDisplayCenterFrequency(),
		DisplayBandwidth:       spyserver.GetDisplayBandwidth(),
		DisplayOffset:          spyserver.GetDisplayOffset(),
		DisplayRange:           spyserver.GetDisplayRange(),
		DisplayPixels:          spyserver.GetDisplayPixels(),

		CurrentSampleRate:      spyserver.GetSampleRate(),
		ChannelCenterFrequency: spyserver.GetCenterFrequency(),
		Gain:              spyserver.GetGain(),
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
	sendPacket := currDevice.Gain != spy2go.InvalidValue
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

	var spyserver = spy2go.MakeSpyserverByFullHS(spyserverhost)

	var cb = spy2go.CallbackBase{
		OnDeviceSync: func() { onDeviceSync(spyserver) },
		OnUInt8IQ:    onUInt8IQ,
		OnInt16IQ:    onInt16IQ,
		OnFFT:        onFFT,
	}

	spyserver.SetCallback(&cb)

	spyserver.Connect()
	defer spyserver.Disconnect()

	log.Println(fmt.Sprintf("Device: %s", spyserver.GetName()))
	var srs = spyserver.GetAvailableSampleRates()

	log.Println("Available SampleRates:")
	for i := 0; i < len(srs); i++ {
		log.Println(fmt.Sprintf("		%f msps (dec stage %d)", float32(srs[i])/1e6, i))
	}

	spyserver.SetStreamingMode(spy2go.StreamModeFFTIQ)
	//spyserver.SetStreamingMode(spy2go.StreamModeIQOnly)
	spyserver.SetDisplayPixels(uint32(displayPixels))
	spyserver.SetDisplayDecimationStage(uint32(displayDecimationStage))
	spyserver.SetDisplayCenterFrequency(uint32(displayFrequency))
	spyserver.SetDisplayRange(90)
	spyserver.SetDisplayOffset(0)

	if spyserver.SetDecimationStage(uint32(channelDecimationStage)) == spy2go.InvalidValue {
		log.Println("Error setting sample rate.")
	}
	if spyserver.SetCenterFrequency(uint32(channelFrequency)) == spy2go.InvalidValue {
		log.Println("Error setting center frequency.")
	}

	time.Sleep(10 * time.Millisecond)

	log.Println("IQ Sample Rate: ", spyserver.GetSampleRate())
	log.Println("IQ Center Frequency ", spyserver.GetCenterFrequency())
	log.Println("FFT Center Frequency ", spyserver.GetDisplayCenterFrequency())
	log.Println("FFT Sample Rate: ", spyserver.GetDisplaySampleRate())

	demodulator = buildDSP(spyserver.GetSampleRate())
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
	spyserver.Start()

	<-done

	err = srv.Shutdown(context.TODO())
	if err != nil {
		log.Println(err)
	}

	log.Print("Stopping")
	spyserver.Stop()
	stopDSP()

	fmt.Println("Work Done")
}
