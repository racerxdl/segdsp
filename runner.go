package main

import (
	"github.com/racerxdl/spy2go"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"os"
	"runtime/pprof"
	"os/signal"
	"syscall"
	"time"
	"github.com/racerxdl/segdsp/demodcore"
	"github.com/racerxdl/segdsp/eventmanager"
	"github.com/racerxdl/segdsp/recorders"
)

func OnInt16IQ(data []spy2go.ComplexInt16) {
	go AddS16Fifo(data)
}
func OnUInt8IQ(data []spy2go.ComplexUInt8) {
	go AddU8Fifo(data)
}

func OnDeviceSync(spyserver *spy2go.Spyserver) {
	var d = DeviceMessage{
		DeviceName: spyserver.GetName(),

		DisplayCenterFrequency: spyserver.GetDisplayCenterFrequency(),
		DisplayBandwidth: spyserver.GetDisplayBandwidth(),
		DisplayOffset: spyserver.GetDisplayOffset(),
		DisplayRange: spyserver.GetDisplayRange(),
		DisplayPixels: spyserver.GetDisplayPixels(),

		CurrentSampleRate: spyserver.GetSampleRate(),
		ChannelCenterFrequency: spyserver.GetCenterFrequency(),
		Gain: spyserver.GetGain(),
		OutputRate: uint32(outputRate),
		FilterBandwidth: uint32(filterBandwidth),
		DemodulatorMode: demodulatorMode,
		DemodulatorParams: nil,
		StationName: stationName,
		WebCanControl: webCanControl,
		TCPCanControl: tcpCanControl,
		IsMuted: false,
	}

	if demodulator != nil {
		d.DemodulatorParams = demodulator.GetDemodParams()
		d.IsMuted = demodulator.IsMuted()
	}

	currDevice = MakeDeviceMessage(d)
	RefreshDevice()
}

func RefreshDevice() {
	sendPacket := currDevice.Gain != spy2go.InvalidValue
	if sendPacket {
		m, err := json.Marshal(currDevice)
		if err != nil {
			log.Println("Error serializing JSON: ", err)
		}
		go broadcastMessage(string(m))
	}
}

func OnFFT(data []uint8) {
	//log.Println("Received FFT! ", len(data))
	var j = MakeFFTMessage(data)
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
		go RecordAudio(b.Data)
		break
	default:
		var j = MakeDataMessage(data)
		m, err := json.Marshal(j)
		if err != nil {
			log.Println("Error serializing JSON: ", err)
		}
		go broadcastMessage(string(m))
		break
	}
	//log.Println("Sending buffer")
}

func CreateServer() *http.Server {
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

func OnSquelchOn(data eventmanager.SquelchEventData) {
	log.Println("Squelch ON", data.AvgValue, data.Threshold)
	currDevice.IsMuted = demodulator.IsMuted()
	StopRecording()
	RefreshDevice()
}

func OnSquelchOff(data eventmanager.SquelchEventData) {
	log.Println("Squelch OFF", data.AvgValue, data.Threshold)
	currDevice.IsMuted = demodulator.IsMuted()
	StartRecording()
	RefreshDevice()
}

func main() {
	SetEnv()
	log.SetFlags(0)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Starting CPU Profile")
		pprof.StartCPUProfile(f)
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
				OnSquelchOn(msg.(eventmanager.SquelchEventData))
			case msg := <-squelchOff:
				OnSquelchOff(msg.(eventmanager.SquelchEventData))
			}
		}
		log.Println("Ending Handler loop")
	}()

	recorder = &recorders.FileRecorder{}
	recordingParams.recorderEnable = record

	if recordMethod != "file" {
		panic("Only\"file\" method is supported for recording.")
	}

	InitDSP()

	var spyserver = spy2go.MakeSpyserverByFullHS(spyserverhost)

	var cb = spy2go.CallbackBase{
		OnDeviceSync: func() { OnDeviceSync(spyserver) },
		OnUInt8IQ: OnUInt8IQ,
		OnInt16IQ: OnInt16IQ,
		OnFFT: OnFFT,
	}

	spyserver.SetCallback(&cb)

	spyserver.Connect()
	defer spyserver.Disconnect()

	log.Println(fmt.Sprintf("Device: %s", spyserver.GetName()))
	var srs = spyserver.GetAvailableSampleRates()

	log.Println("Available SampleRates:")
	for i := 0; i < len(srs); i++ {
		log.Println(fmt.Sprintf("		%f msps (dec stage %d)", float32(srs[i]) / 1e6, i))
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

	demodulator = BuildDSP(spyserver.GetSampleRate())
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

	var srv = CreateServer()

	StartDSP()

	log.Println("Starting")
	spyserver.Start()

	<-done

	srv.Shutdown(nil)

	log.Print("Stopping")
	spyserver.Stop()
	StopDSP()

	fmt.Println("Work Done")
}