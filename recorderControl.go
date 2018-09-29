package main

import (
	"fmt"
	"github.com/racerxdl/segdsp/recorders"
	"sync"
	"time"
)

var recorder recorders.BaseRecorder
var recordMutex = sync.Mutex{}

var recordingParams = struct {
	baseFilename   string
	params         []interface{}
	recorderEnable bool
	recording      bool
}{
	baseFilename:   "%s-%s",
	params:         make([]interface{}, 0),
	recorderEnable: false,
	recording:      false,
}

type recordingMetadata struct {
	DemodParams  interface{}
	BaseFilename string
	Timestamp    time.Time
}

func startRecording() {
	recordMutex.Lock()
	if recordingParams.recorderEnable {
		if recorder != nil {
			recorder.Close()
		}
		var filename = fmt.Sprintf(recordingParams.baseFilename, stationName, time.Now().Local().Format("20060102_150405"))
		var newParams = []interface{}{
			filename,
			recordingMetadata{
				DemodParams:  demodulator.GetDemodParams(),
				BaseFilename: filename,
				Timestamp:    time.Now().Local(),
			},
		}

		recorder.Open(newParams)
		recordingParams.recording = true
	}
	recordMutex.Unlock()
}

func recordIQ(data []complex64) {
	recordMutex.Lock()
	if recordingParams.recorderEnable && recorder != nil && recordingParams.recording {
		go recorder.WriteIQ(data)
	}
	recordMutex.Unlock()
}

func recordAudio(data []float32) {
	recordMutex.Lock()
	if recordingParams.recorderEnable && recorder != nil && recordingParams.recording {
		go recorder.WriteAudio(data)
	}
	recordMutex.Unlock()
}

func recordData(data []byte) {
	recordMutex.Lock()
	if recordingParams.recorderEnable && recorder != nil && recordingParams.recording {
		go recorder.WriteData(data)
	}
	recordMutex.Unlock()
}

func stopRecording() {
	recordMutex.Lock()
	recorder.Close()
	recordingParams.recording = false
	recordMutex.Unlock()
}
