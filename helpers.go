package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"math"
)

type jsonUint8s []uint8

func (u jsonUint8s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("")
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	_, err := enc.Write(u)
	if err != nil {
		panic(err)
	}
	enc.Close()
	return json.Marshal(buf.String())
}

//type jsonInt16s []int16
//
//func (u jsonInt16s) MarshalJSON() ([]byte, error) {
//	var result string
//	if u == nil {
//		result = "null"
//	} else {
//		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
//	}
//	return []byte(result), nil
//}

type fftMessage struct {
	MessageType      string
	DemodOutputLevel float32
	FFTData          jsonUint8s
}

type dataMessage struct {
	MessageType string
	Data        interface{}
}

type deviceMessage struct {
	MessageType string

	DeviceName             string
	DisplayBandwidth       uint32
	DisplayCenterFrequency uint32
	DisplayOffset          int32
	DisplayRange           int32
	DisplayPixels          uint32

	ChannelCenterFrequency uint32
	CurrentSampleRate      uint32
	Gain                   uint32
	OutputRate             uint32

	FilterBandwidth   uint32
	DemodulatorMode   string
	DemodulatorParams interface{}

	StationName   string
	WebCanControl bool
	TCPCanControl bool
	IsMuted       bool
}

func makeFFTMessage(data []uint8, level float32) fftMessage {
	if math.IsInf(float64(level), 0) {
		level = 0
	}
	return fftMessage{
		MessageType:      "fft",
		DemodOutputLevel: level,
		FFTData:          data,
	}
}

func makeDataMessage(data interface{}) dataMessage {
	return dataMessage{
		MessageType: "data",
		Data:        data,
	}
}

func makeDeviceMessage(d deviceMessage) deviceMessage {
	d.MessageType = "device"

	return d
}
