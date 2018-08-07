package main

import (
	"strings"
	"fmt"
	"encoding/base64"
	"bytes"
	"encoding/json"
)

type JsonUint8 []uint8

func (u JsonUint8) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("")
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	enc.Write(u)
	enc.Close()
	return json.Marshal(string(buf.Bytes()))
}

type JsonInt16 []int16

func (u JsonInt16) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}


type FFTMessage struct {
	MessageType string
	FFTData JsonUint8
}

type DataMessage struct {
	MessageType string
	Data interface{}
}

type DeviceMessage struct {
	MessageType 			string

	DeviceName				string
	DisplayBandwidth		uint32
	DisplayCenterFrequency  uint32
	DisplayOffset			int32
	DisplayRange			int32
	DisplayPixels			uint32

	ChannelCenterFrequency  uint32
	CurrentSampleRate		uint32
	Gain					uint32
	OutputRate				uint32

	FilterBandwidth			uint32
	DemodulatorMode			string
	DemodulatorParams		interface{}

	StationName				string
	WebCanControl			bool
	TCPCanControl			bool
	IsMuted					bool
}

func MakeFFTMessage(data []uint8) FFTMessage {
	return FFTMessage{
		MessageType: "fft",
		FFTData: data,
	}
}

func MakeDataMessage(data interface{}) DataMessage {
	return DataMessage{
		MessageType: "data",
		Data: data,
	}
}

func MakeDeviceMessage(d DeviceMessage) DeviceMessage {
	d.MessageType = "device"

	return d
}