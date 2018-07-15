package main

import (
	"strings"
	"fmt"
)

type JsonUint8 []uint8

func (u JsonUint8) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
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
	Data JsonInt16
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
}

func MakeFFTMessage(data []uint8) FFTMessage {
	return FFTMessage{
		MessageType: "fft",
		FFTData: data,
	}
}

func MakeDataMessage(data []int16) DataMessage {
	return DataMessage{
		MessageType: "data",
		Data: data,
	}
}

func MakeDeviceMessage(d DeviceMessage) DeviceMessage {
	d.MessageType = "device"

	return d
}