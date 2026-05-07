package main

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/segdsp/demodcore"
	"runtime"
	"sync/atomic"
)

var samplesFifo *fifo.Queue
var demodulator demodcore.DemodCore
var running atomic.Bool
var buffer []complex64

var dspCb func(interface{})

func addComplex(data []complex64) {
	samplesFifo.Add(data)
}

func initDSP() {
	samplesFifo = fifo.NewQueue()
}

func startDSP() {
	if !running.Load() {
		running.Store(true)
		go dspLoop()
	}
}

func stopDSP() {
	running.Store(false)
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
	for running.Load() {
		dspRun()
		runtime.Gosched()
	}
}
