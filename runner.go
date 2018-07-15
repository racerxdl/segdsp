package main

import (
	"github.com/racerxdl/spy2go"
	"fmt"
	"log"
	"flag"
	"net/http"
	"encoding/json"
)

func OnInt16IQ(data []spy2go.ComplexInt16) {
	//log.Println("Received Complex 16 Data! ", len(data))
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

	}
	currDevice = MakeDeviceMessage(d)
	m, err := json.Marshal(currDevice)
	if err != nil {
		log.Println("Error serializing JSON: ", err)
	}
	broadcastMessage(string(m))
}
func OnFFT(data []uint8) {
	//log.Println("Received FFT! ", len(data))
	var j = MakeFFTMessage(data)
	m, err := json.Marshal(j)
	if err != nil {
		log.Println("Error serializing JSON: ", err)
	}
	broadcastMessage(string(m))
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var spyserverhost = flag.String("spyserver", "localhost:5555", "spyserver address")
var displayPixels = flag.Uint("displayPixels", 512, "Width in pixels of the FFT")

var channelFrequency = flag.Uint("channelFrequency", 106300000, "Channel (IQ) Center Frequency")
var displayFrequency = flag.Uint("fftFrequency", 106000000, "FFT Center Frequency")

var channelDecimationStage = flag.Uint("decimationStage", 3, "Channel (IQ) Decimation Stage (The actual decimation will be 2^d)")
var displayDecimationStage = flag.Uint("fftDecimationStage", 1, "FFT Decimation Stage (The actual decimation will be 2^d)")

func main() {

	flag.Parse()
	log.SetFlags(0)

	fs := http.FileServer(http.Dir("./content/static"))

	http.HandleFunc("/ws", ws)
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", content)

	var spyserver = spy2go.MakeSpyserverByFullHS(*spyserverhost)

	var cb = spy2go.CallbackBase{
		OnDeviceSync: func() { OnDeviceSync(spyserver) },
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
		log.Println(fmt.Sprintf("		%f msps", float32(srs[i]) / 1e6))
	}
	if spyserver.SetDecimationStage(uint32(*channelDecimationStage)) == 0xFFFFFFFF {
		log.Println("Error setting sample rate.")
	}
	if spyserver.SetCenterFrequency(uint32(*channelFrequency)) == 0xFFFFFFFF {
		log.Println("Error setting center frequency.")
	}

	spyserver.SetStreamingMode(spy2go.StreamModeFFTIQ)
	spyserver.SetDisplayPixels(uint32(*displayPixels))
	spyserver.SetDisplayRange(90)
	spyserver.SetDisplayOffset(10)
	spyserver.SetDisplayDecimationStage(uint32(*displayDecimationStage))
	spyserver.SetDisplayCenterFrequency(uint32(*displayFrequency))

	log.Println("Starting")
	spyserver.Start()

	log.Print(http.ListenAndServe(*addr, nil))

	log.Print("Stopping")
	spyserver.Stop()
}