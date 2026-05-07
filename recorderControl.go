package main

import (
	"fmt"
	"github.com/racerxdl/segdsp/demodcore"
	"github.com/racerxdl/segdsp/recorders"
	"sync"
	"time"
)

type recordingMetadata struct {
	DemodParams  interface{}
	BaseFilename string
	Timestamp    time.Time
}

type RecordingManager struct {
	mu             sync.Mutex
	recorder       recorders.BaseRecorder
	baseFilename   string
	RecorderEnable bool
	recording      bool
	getDemodParams func() interface{}
}

func NewRecordingManager(recorder recorders.BaseRecorder, getDemodParams func() interface{}) *RecordingManager {
	return &RecordingManager{
		recorder:       recorder,
		baseFilename:   "%s-%s",
		getDemodParams: getDemodParams,
	}
}

func (r *RecordingManager) StartRecording(stationName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.RecorderEnable {
		return
	}
	if r.recorder != nil {
		r.recorder.Close()
	}
	var filename = fmt.Sprintf(r.baseFilename, stationName, time.Now().Local().Format("20060102_150405"))
	var newParams = []interface{}{
		filename,
		recordingMetadata{
			DemodParams:  r.getDemodParams(),
			BaseFilename: filename,
			Timestamp:    time.Now().Local(),
		},
	}

	r.recorder.Open(newParams)
	r.recording = true
}

func (r *RecordingManager) RecordAudio(data []float32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.RecorderEnable && r.recorder != nil && r.recording {
		r.recorder.WriteAudio(data)
	}
}

func (r *RecordingManager) StopRecording() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.recorder != nil {
		r.recorder.Close()
	}
	r.recording = false
}

func (r *RecordingManager) ProcessDemodData(data demodcore.DemodData) {
	r.RecordAudio(data.Data)
}
