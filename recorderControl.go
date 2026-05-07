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
	defer recordMutex.Unlock()
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
}

func recordAudio(data []float32) {
	recordMutex.Lock()
	defer recordMutex.Unlock()
	if recordingParams.recorderEnable && recorder != nil && recordingParams.recording {
		recorder.WriteAudio(data)
	}
}

func stopRecording() {
	recordMutex.Lock()
	defer recordMutex.Unlock()
	recorder.Close()
	recordingParams.recording = false
}
