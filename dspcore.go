package main

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/radioserver/client"
	"github.com/racerxdl/segdsp/demodcore"
	"runtime"
	"time"
)

var samplesFifo *fifo.Queue
var demodulator demodcore.DemodCore
var running = false
var buffer []complex64
var delta = 0.0
var count = 0

var dspCb func(interface{})

func addS16Fifo(data []client.ComplexInt16) {
	var chunk = make([]complex64, len(data))
	for i := 0; i < len(data); i++ {
		chunk[i] = complex(float32(data[i].Imag)/32768.0, float32(data[i].Real)/32768.0)
	}

	samplesFifo.Add(chunk)
}

func initDSP() {
	samplesFifo = fifo.NewQueue()
}

func startDSP() {
	if !running {
		running = true
		go dspLoop()
	}
}

func stopDSP() {
	if running {
		running = false
	}
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

	var t0 = time.Now()
	var out = demodulator.Work(buffer)
	var d = time.Since(t0)
	delta += d.Seconds()
	count++

	if out != nil {
		delta /= float64(count)
		//log.Println("Delta: ", delta, "seconds")
		delta = 0
		count = 0
		dspCb(out)
	}
}

func dspLoop() {
	for running {
		dspRun()
		runtime.Gosched()
	}
}
