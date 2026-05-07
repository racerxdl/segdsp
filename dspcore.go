package main

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/segdsp/demodcore"
)

var samplesFifo *fifo.Queue
var demodulator demodcore.DemodCore
var buffer []complex64
var dspCb func(interface{})
var dspReady chan struct{}
var dspDone chan struct{}

func addComplex(data []complex64) {
	samplesFifo.Add(data)
	select {
	case dspReady <- struct{}{}:
	default:
	}
}

func initDSP() {
	samplesFifo = fifo.NewQueue()
	dspReady = make(chan struct{}, 1)
	dspDone = make(chan struct{})
}

func startDSP() {
	go dspLoop()
}

func stopDSP() {
	close(dspReady)
	<-dspDone
}

func dspRun() {
	length := samplesFifo.Len()

	if length == 0 {
		return
	}

	if demodulator == nil || dspCb == nil {
		return
	}

	buffer = samplesFifo.Next().([]complex64)

	var out = demodulator.Work(buffer)

	if out != nil {
		dspCb(out)
	}
}

func dspLoop() {
	defer close(dspDone)
	for range dspReady {
		for samplesFifo.Len() > 0 {
			dspRun()
		}
	}
}
